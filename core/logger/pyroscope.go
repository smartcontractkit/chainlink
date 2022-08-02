package logger

import (
	"runtime"

	"github.com/pkg/errors"
	"github.com/pyroscope-io/client/pyroscope"
)

// PyroscopeServerAddrMissingErr error to be triggered if the server address is missing
var PyroscopeServerAddrMissingErr = errors.New("pyroscope server address is missing")

// PyroscopeConfig represents the expected configuration for Pyroscope to properly work
type PyroscopeConfig interface {
	PyroscopeServerAddress() string
	PyroscopeAuthToken() string

	AutoPprofBlockProfileRate() int
	AutoPprofMutexProfileFraction() int
}

// StartPyroscope starts continuous profiling of the Chainlink Node
func StartPyroscope(cfg PyroscopeConfig) (*pyroscope.Profiler, error) {
	if cfg.PyroscopeServerAddress() == "" {
		return nil, PyroscopeServerAddrMissingErr
	}

	runtime.SetBlockProfileRate(cfg.AutoPprofBlockProfileRate())
	runtime.SetMutexProfileFraction(cfg.AutoPprofMutexProfileFraction())

	return pyroscope.Start(pyroscope.Config{
		// Maybe configurable to identify the specific NOP
		ApplicationName: "chainlink-node",

		// TBD
		ServerAddress: cfg.PyroscopeServerAddress(),
		AuthToken:     cfg.PyroscopeAuthToken(),

		// may end up disabling logging
		Logger: pyroscope.StandardLogger,

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
}
