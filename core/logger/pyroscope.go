package logger

import (
	"runtime"

	"github.com/pkg/errors"
	"github.com/pyroscope-io/client/pyroscope"
)

// ErrPyroscopeServerAddrMissing error to be triggered if the server address is missing
var ErrPyroscopeServerAddrMissing = errors.New("pyroscope server address is missing")

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
		return nil, ErrPyroscopeServerAddrMissing
	}

	runtime.SetBlockProfileRate(cfg.AutoPprofBlockProfileRate())
	runtime.SetMutexProfileFraction(cfg.AutoPprofMutexProfileFraction())

	return pyroscope.Start(pyroscope.Config{
		// Maybe configurable to identify the specific NOP - TBD
		ApplicationName: "chainlink-node",

		ServerAddress: cfg.PyroscopeServerAddress(),
		AuthToken:     cfg.PyroscopeAuthToken(),

		// We disable logging the profiling info, it will be in the Pyroscope instance anyways...
		Logger: nil,

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
