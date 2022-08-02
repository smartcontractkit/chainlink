package logger

import (
	"fmt"
	"runtime"

	"github.com/pyroscope-io/client/pyroscope"
)

type PyroscopeConfig interface {
	PyroscopeServerAddress() string
	PyroscopeAuthToken() string

	AutoPprofBlockProfileRate() int
	AutoPprofMutexProfileFraction() int
}

// StartPyroscope starts continuous profiling of the Chainlink Node
func StartPyroscope(cfg PyroscopeConfig) (*pyroscope.Profiler, error) {
	runtime.SetBlockProfileRate(cfg.AutoPprofBlockProfileRate())
	runtime.SetMutexProfileFraction(cfg.AutoPprofMutexProfileFraction())

	return pyroscope.Start(pyroscope.Config{
		// Maybe configurable to identify the specific NOP
		ApplicationName: fmt.Sprintf("chainlink-node"),

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
