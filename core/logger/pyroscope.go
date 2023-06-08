package logger

import (
	"runtime"

	"github.com/pyroscope-io/client/pyroscope"

	"github.com/smartcontractkit/chainlink/v2/core/static"
)

// PyroscopeConfig represents the expected configuration for Pyroscope to properly work
type PyroscopeProfilerConfig interface {
	Pyroscope() config.Pyroscope

type PprofConfig interface {
	BlockProfileRate() int
	MutexProfileFraction() int
}

// StartPyroscope starts continuous profiling of the Chainlink Node
func StartPyroscope(cfg PyroscopeProfilerConfig, pprofConfig PprofConfig) (*pyroscope.Profiler, error) {
	runtime.SetBlockProfileRate(pprofConfig.BlockProfileRate())
	runtime.SetMutexProfileFraction(pprofConfig.MutexProfileFraction())

	sha, ver := static.Short()

	return pyroscope.Start(pyroscope.Config{
		// Maybe configurable to identify the specific NOP - TBD
		ApplicationName: "chainlink-node",

		ServerAddress: cfg.Pyroscope().ServerAddress(),
		AuthToken:     cfg.Pyroscope().AuthToken(),

		// We disable logging the profiling info, it will be in the Pyroscope instance anyways...
		Logger: nil,

		Tags: map[string]string{
			"SHA":         sha,
			"Version":     ver,
			"Environment": cfg.Pyroscope().Environment(),
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
