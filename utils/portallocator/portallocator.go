package portallocator

import (
	"net"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// allocator
type allocator struct {
	allocate chan string
	portNum  int
	size     int
	start    int
}

// New provides a new allocator
func New(start, size int) Allocator {

	a := &allocator{
		allocate: make(chan string, size),
		portNum:  start,
		start:    start,
		size:     size,
	}
	//count := 0
	zap.L().Debug("Started Binding for reserving ports", zap.Time("Start", time.Now()))
	for i := start; len(a.allocate) < size; i++ {
		if i > ((1 << 16) - 1) {
			zap.L().Error("Could not reserve 100 ports for enforcerproxy")
			return nil
		}
		addr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(i))

		if err != nil {
			zap.L().Debug("Resolve TCP failed", zap.Error(err))
			continue
		}

		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			zap.L().Debug("Socket failed", zap.Error(err))
			continue
		}
		if len(addr.IP) == 0 {
			addr.IP = net.IPv4zero
		}
		socketAddress := &syscall.SockaddrInet4{Port: addr.Port}
		copy(socketAddress.Addr[:], addr.IP.To4())
		//set REUSEPORT or REUSEADDR so application proxy can still bind to these later
		if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {

			return nil
		}
		if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, 15, 1); err != nil {
			return nil
		}
		if err = syscall.Bind(fd, socketAddress); err != nil {
			syscall.Close(fd) // nolint errcheck
			zap.L().Debug("Bind failed", zap.Error(err))
			continue
		}
		if err = syscall.Listen(fd, 100); err != nil {
			syscall.Close(fd) // nolint errcheck
			zap.L().Debug("Listen failed", zap.Error(err))
			continue
		}
		a.allocate <- strconv.Itoa(i)

	}
	zap.L().Debug("Done Binding for reserving ports", zap.Time("End", time.Now()), zap.Int("Reserved Ports", len(a.allocate)))
	return a
}

// Allocate allocates an item
func (p *allocator) Allocate() string {
	return <-p.allocate
}

// Release releases an item
func (p *allocator) Release(item string) {
	p.allocate <- item
}
