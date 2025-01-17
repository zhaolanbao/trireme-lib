// +build linux

package markedconn

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"
	"unsafe"

	"go.aporeto.io/trireme-lib/utils/netinterfaces"
	"go.uber.org/zap"
)

const (
	sockOptOriginalDst = 80
)

// DialMarkedWithContext will dial a TCP connection to the provide address and mark the socket
// with the provided mark.
func DialMarkedWithContext(ctx context.Context, network string, addr string, mark int) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, _ string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {

				if err := syscall.SetNonblock(int(fd), false); err != nil {
					zap.L().Error("unable to set socket options", zap.Error(err))
				}
				if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, mark); err != nil {
					zap.L().Error("Failed to assing mark to socket", zap.Error(err))
				}
				if err := syscall.SetsockoptInt(int(fd), syscall.SOL_TCP, 30, 1); err != nil {
					zap.L().Debug("Failed to set fast open socket option", zap.Error(err))
				}
			})
		},
	}

	conn, err := d.DialContext(ctx, network, addr)
	if err != nil {
		zap.L().Error("Failed to dial to downstream node",
			zap.Error(err),
			zap.String("Address", addr),
			zap.String("Network type", network),
		)
	}
	return conn, err
}

// NewSocketListener will create a listener and mark the socket with the provided mark.
func NewSocketListener(ctx context.Context, port string, mark int) (net.Listener, error) {
	listenerCfg := net.ListenConfig{
		Control: func(_, _ string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, mark); err != nil {
					zap.L().Error("Failed to mark connection", zap.Error(err))
				}
				if err := syscall.SetsockoptInt(int(fd), syscall.SOL_TCP, 23, 16*1024); err != nil {
					zap.L().Error("Cannot set tcp fast open options", zap.Error(err))
				}
			})
		},
	}

	listener, err := listenerCfg.Listen(ctx, "tcp", port)

	if err != nil {
		return nil, fmt.Errorf("Failed to create listener: %s", err)
	}

	return ProxiedListener{netListener: listener, mark: mark}, nil
}

// ProxiedConnection is a proxied connection where we can recover the
// original destination.
type ProxiedConnection struct {
	originalIP            net.IP
	originalPort          int
	originalTCPConnection *net.TCPConn
}

// GetOriginalDestination sets the original destination of the connection.
func (p *ProxiedConnection) GetOriginalDestination() (net.IP, int) {
	return p.originalIP, p.originalPort
}

// GetTCPConnection returns the TCP connection object.
func (p *ProxiedConnection) GetTCPConnection() *net.TCPConn {
	return p.originalTCPConnection
}

// LocalAddr implements the corresponding method of net.Conn, but returns the original
// address.
func (p *ProxiedConnection) LocalAddr() net.Addr {

	return &net.TCPAddr{
		IP:   p.originalIP,
		Port: p.originalPort,
	}
}

// RemoteAddr returns the remote address
func (p *ProxiedConnection) RemoteAddr() net.Addr {
	return p.originalTCPConnection.RemoteAddr()
}

// Read reads data from the connection.
func (p *ProxiedConnection) Read(b []byte) (n int, err error) {
	return p.originalTCPConnection.Read(b)
}

// Write writes data to the connection.
func (p *ProxiedConnection) Write(b []byte) (n int, err error) {
	return p.originalTCPConnection.Write(b)
}

// Close closes the connection.
func (p *ProxiedConnection) Close() error {
	return p.originalTCPConnection.Close()
}

// SetDeadline passes the read deadline to the original TCP connection.
func (p *ProxiedConnection) SetDeadline(t time.Time) error {
	return p.originalTCPConnection.SetDeadline(t)
}

// SetReadDeadline implements the call by passing it to the original connection.
func (p *ProxiedConnection) SetReadDeadline(t time.Time) error {
	return p.originalTCPConnection.SetReadDeadline(t)
}

// SetWriteDeadline implements the call by passing it to the original connection.
func (p *ProxiedConnection) SetWriteDeadline(t time.Time) error {
	return p.originalTCPConnection.SetWriteDeadline(t)
}

// ProxiedListener is a proxied listener that uses proxied connections.
type ProxiedListener struct {
	netListener net.Listener
	mark        int
}

// Accept implements the accept method of the interface.
func (l ProxiedListener) Accept() (c net.Conn, err error) {
	nc, err := l.netListener.Accept()
	if err != nil {
		return nil, err
	}

	tcpConn, ok := nc.(*net.TCPConn)
	if !ok {
		zap.L().Error("Received a non-TCP connection - this should never happen", zap.Error(err))
		return nil, fmt.Errorf("Not a tcp connection - ignoring")
	}

	ip, port, err := GetOriginalDestination(tcpConn)
	if err != nil {
		zap.L().Error("Failed to discover original destination - aborting", zap.Error(err))
		return nil, err
	}

	return &ProxiedConnection{
		originalIP:            ip,
		originalPort:          port,
		originalTCPConnection: tcpConn,
	}, nil
}

// Addr implements the Addr method of net.Listener.
func (l ProxiedListener) Addr() net.Addr {
	return l.netListener.Addr()
}

// Close implements the Close method of the net.Listener.
func (l ProxiedListener) Close() error {
	return l.netListener.Close()
}

type sockaddr4 struct {
	family uint16
	data   [14]byte
}

type sockaddr6 struct {
	family   uint16
	port     [2]byte
	flowInfo [4]byte //nolint
	ip       [16]byte
	scopeID  [4]byte //nolint
}

type origDest func(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)
type passFD interface {
	Control(func(uintptr)) error
}

func getOriginalDestInternal(rawConn passFD, v4Proto bool, getOrigDest origDest) (net.IP, int, error) { // nolint interfacer{
	var getsockopt func(fd uintptr)
	var netIP net.IP
	var port int
	var err error

	getsockopt4 := func(fd uintptr) {
		var addr sockaddr4
		size := uint32(unsafe.Sizeof(addr))
		_, _, e1 := getOrigDest(syscall.SYS_GETSOCKOPT, uintptr(fd), uintptr(syscall.SOL_IP), uintptr(sockOptOriginalDst), uintptr(unsafe.Pointer(&addr)), uintptr(unsafe.Pointer(&size)), 0) //nolint

		if e1 != 0 {
			err = fmt.Errorf("Failed to get original destination: %s", e1)
			return
		}

		if addr.family != syscall.AF_INET {
			err = fmt.Errorf("invalid address family. Expected AF_INET")
			return
		}

		netIP = addr.data[2:6]
		port = int(addr.data[0])<<8 + int(addr.data[1])
	}

	getsockopt6 := func(fd uintptr) {
		var addr sockaddr6
		size := uint32(unsafe.Sizeof(addr))

		_, _, e1 := getOrigDest(syscall.SYS_GETSOCKOPT, uintptr(fd), uintptr(syscall.SOL_IPV6), uintptr(sockOptOriginalDst), uintptr(unsafe.Pointer(&addr)), uintptr(unsafe.Pointer(&size)), 0) //nolint
		if e1 != 0 {
			err = fmt.Errorf("Failed to get original destination: %s", e1)
			return
		}

		if addr.family != syscall.AF_INET6 {
			err = fmt.Errorf("invalid address family. Expected AF_INET6")
			return
		}

		netIP = addr.ip[:]
		port = int(addr.port[0])<<8 + int(addr.port[1])
	}

	if v4Proto {
		getsockopt = getsockopt4
	} else {
		getsockopt = getsockopt6
	}

	if err1 := rawConn.Control(getsockopt); err1 != nil {
		return nil, 0, fmt.Errorf("Failed to get original destination: %s", err)
	}

	if err != nil {
		return nil, 0, err
	}

	return netIP, port, nil
}

// GetOriginalDestination -- Func to get original destination a connection
func GetOriginalDestination(conn *net.TCPConn) (net.IP, int, error) { // nolint interfacer

	rawconn, err := conn.SyscallConn()
	if err != nil {
		return nil, 0, err
	}

	localIPString, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return nil, 0, err
	}

	localIP := net.ParseIP(localIPString)
	return getOriginalDestInternal(rawconn, localIP.To4() != nil, syscall.Syscall6)
}

// GetInterfaces retrieves all the local interfaces.
func GetInterfaces() map[string]struct{} {
	ipmap := map[string]struct{}{}

	ifaces, err := netinterfaces.GetInterfacesInfo()
	if err != nil {
		zap.L().Debug("Unable to get interfaces info", zap.Error(err))
	}

	for _, iface := range ifaces {
		for _, ip := range iface.IPs {
			ipmap[ip.String()] = struct{}{}
		}
	}

	return ipmap
}
