package env

import (
	"fmt"
	"os"
	"strconv"

	"go.aporeto.io/trireme-lib/controller/constants"
	"go.aporeto.io/trireme-lib/controller/pkg/claimsheader"
)

// RemoteParameters holds all configuration objects that must be passed
// during the initialization of the monitor.
type RemoteParameters struct {
	LogToConsole      bool
	LogWithID         bool
	LogLevel          string
	LogFormat         string
	CompressedTags    claimsheader.CompressionType
	KubernetesEnabled bool
}

// GetParameters retrieves log parameters for Remote Enforcer.
func GetParameters() (logToConsole bool, logID string, logLevel string, logFormat string, compressedTagsVersion claimsheader.CompressionType, kubernetesEnabled bool) {

	logLevel = os.Getenv(constants.EnvLogLevel)
	if logLevel == "" {
		logLevel = "info"
	}
	//kubernetesEnabled := os.Getenv(constants.)
	logFormat = os.Getenv(constants.EnvLogFormat)
	if logLevel == "" {
		logFormat = "json"
	}

	if console := os.Getenv(constants.EnvLogToConsole); console == constants.EnvLogToConsoleEnable {
		logToConsole = true
	}

	logID = os.Getenv(constants.EnvLogID)

	compressedTagsVersion = claimsheader.CompressionTypeNone
	if console := os.Getenv(constants.EnvCompressedTags); console != string(claimsheader.CompressionTypeNone) {
		if console == string(claimsheader.CompressionTypeV1) {
			compressedTagsVersion = claimsheader.CompressionTypeV1
		} else if console == string(claimsheader.CompressionTypeV2) {
			compressedTagsVersion = claimsheader.CompressionTypeV2
		}
	}
	kubernetesEnabled, _ = strconv.ParseBool(os.Getenv(constants.EnvKubernetesEnabled))
	fmt.Println("\n\n Getting the parameters, k8senabled : ", kubernetesEnabled)
	return
}
