package constants

import "time"

const (
	// DefaultProcMountPoint The default proc mountpoint
	DefaultProcMountPoint = "/proc"
	// DefaultAporetoProcMountPoint The aporeto proc mountpoint just in case we are launched with some specific docker config
	DefaultAporetoProcMountPoint = "/aporetoproc"
	// DefaultSecretsPath is the default path for the secrets proxy.
	DefaultSecretsPath = "@secrets"
)

const (
	// DefaultRemoteArg is the default arguments for a remote enforcer
	DefaultRemoteArg = "enforce"
	// DefaultConnMark is the default conn mark for all data packets
	DefaultConnMark = uint32(0xEEEE)
	// DeleteConnmark is the mark used to trigger udp handshake.
	DeleteConnmark = uint32(0xABCD)
)

const (

	// EnvMountPoint is an environment variable which will contain the mount point
	EnvMountPoint = "TRIREME_ENV_PROC_MOUNTPOINT"

	// EnvContextSocket stores the path to the context specific socket
	EnvContextSocket = "TRIREME_ENV_SOCKET_PATH"

	// EnvStatsChannel stores the path to the stats channel
	EnvStatsChannel = "TRIREME_ENV_STATS_CHANNEL_PATH"

	// EnvDebugChannel stores the path to the debug channel
	EnvDebugChannel = "TRIREME_ENV_DEBUG_CHANNEL_PATH"
	// EnvRPCClientSecret is the secret used between RPC client/server
	EnvRPCClientSecret = "TRIREME_ENV_SECRET"

	// EnvStatsSecret is the secret to be used for the stats channel
	EnvStatsSecret = "TRIREME_ENV_STATS_SECRET"

	// EnvContainerPID is the PID of the container
	EnvContainerPID = "TRIREME_ENV_CONTAINER_PID"

	// EnvNSPath is the path of the network namespace
	EnvNSPath = "TRIREME_ENV_NS_PATH"

	// EnvNsenterErrorState stores the error state as reported by remote enforcer
	EnvNsenterErrorState = "TRIREME_ENV_NSENTER_ERROR_STATE"

	// EnvNsenterLogs stores the logs as reported by remote enforcer
	EnvNsenterLogs = "TRIREME_ENV_NSENTER_LOGS"

	// EnvLogLevel store the log level to be used.
	EnvLogLevel = "TRIREME_ENV_LOG_LEVEL"

	// EnvLogFormat store the log format to be used.
	EnvLogFormat = "TRIREME_ENV_LOG_FORMAT"

	// EnvLogToConsole specifies if logs should be sent out to console.
	EnvLogToConsole = "TRIREME_ENV_LOG_TO_CONSOLE"

	// EnvLogToConsoleEnable specifies value to enable logging to console.
	EnvLogToConsoleEnable = "1"

	// EnvLogID store the context Id for the log file to be used.
	EnvLogID = "TRIREME_ENV_LOG_ID"

	// EnvCompressedTags stores whether we should be using compressed tags.
	EnvCompressedTags = "TRIREME_ENV_COMPRESSED_TAGS"

	// EnvDisableLogWrite tells us if we are running in kubernetes, if true don't write the logs to a file.
	EnvDisableLogWrite = "TRIREME_ENV_DISABLE_LOG_WRITE"
)

// ModeType defines the mode of the enforcement and supervisor.
type ModeType int

const (
	// RemoteContainer indicates that the Supervisor is implemented in the
	// container namespace
	RemoteContainer ModeType = iota
	// LocalServer indicates that the Supervisor applies to Linux processes
	LocalServer
	// Sidecar indicates the controller to be in sidecar mode
	Sidecar
)

// LogLevel corresponds to log level of any logger. eg: zap.
type LogLevel string

// Various log levels.
const (
	Info  LogLevel = "Info"
	Debug LogLevel = "Debug"
	Trace LogLevel = "Trace"
	Error LogLevel = "Error"
	Warn  LogLevel = "Warn"
)

// API service related constants
const (
	CallbackURIExtension = "/aporeto/oidc/callback"
)

// Protocol constants
const (
	TCPProtoNum    = "6"
	UDPProtoNum    = "17"
	TCPProtoString = "TCP"
	UDPProtoString = "UDP"
)

// sockets
const (
	StatsChannel = "/var/run/statschannel.sock"
	DebugChannel = "/var/run/debugchannel.sock"
)

// PortNumberLabelString is the label to use for port numbers
const (
	PortNumberLabelString = "@sys:port"
)

// Token and cache default validities. These have performance implications.
// The faster the datapath issues new tokens it affects performance. However,
// making it too slow can potentially allow reuse of the tokens. The
// token issuance rate must be always faster than the expiration rate.
const (
	// SynTokenCacheValiditity determines how often the data path creates new tokens.
	SynTokenCacheValiditity = 10 * time.Second

	// DatapathTokenValidity determines how long the tokens are valid.
	DatapathTokenValidity = 1 * time.Minute
)
