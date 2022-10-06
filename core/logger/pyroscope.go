package logger

import (
	"runtime"

	"github.com/pyroscope-io/client/pyroscope"

	"github.com/smartcontractkit/chainlink/core/static"
)

// PyroscopeConfig represents the expected configuration for Pyroscope to properly work
type PyroscopeConfig interface {
	PyroscopeServerAddress() string
	PyroscopeAuthToken() string
	PyroscopeEnvironment() string

	AutoPprofBlockProfileRate() int
	AutoPprofMutexProfileFraction() int
}

// StartPyroscope starts continuous profiling of the Chainlink Node
func StartPyroscope(cfg PyroscopeConfig) (*pyroscope.Profiler, error) {
	runtime.SetBlockProfileRate(cfg.AutoPprofBlockProfileRate())
	runtime.SetMutexProfileFraction(cfg.AutoPprofMutexProfileFraction())

	sha, ver := static.Short()

	return pyroscope.Start(pyroscope.Config{
		// Maybe configurable to identify the specific NOP - TBD
		ApplicationName: "chainlink-node",

		ServerAddress: cfg.PyroscopeServerAddress(),
		AuthToken:     cfg.PyroscopeAuthToken(),

		// We disable logging the profiling info, it will be in the Pyroscope instance anyways...
		Logger: nil,

		Tags: map[string]string{
			"SHA":         sha,
			"Version":     ver,
			"Environment": cfg.PyroscopeEnvironment(),
		},

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
