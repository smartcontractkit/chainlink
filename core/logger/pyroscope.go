package logger

import (
	"runtime"

	"github.com/grafana/pyroscope-go"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

type PprofConfig interface {
	BlockProfileRate() int
	MutexProfileFraction() int
}

// StartPyroscope starts continuous profiling of the Chainlink Node
func StartPyroscope(pyroConfig config.Pyroscope, pprofConfig PprofConfig) (*pyroscope.Profiler, error) {
	runtime.SetBlockProfileRate(pprofConfig.BlockProfileRate())
	runtime.SetMutexProfileFraction(pprofConfig.MutexProfileFraction())

	sha, ver := static.Short()

	return pyroscope.Start(pyroscope.Config{
		// Maybe configurable to identify the specific NOP - TBD
		ApplicationName: "chainlink-node",

		ServerAddress: pyroConfig.ServerAddress(),
		AuthToken:     pyroConfig.AuthToken(),

		// We disable logging the profiling info, it will be in the Pyroscope instance anyways...
		Logger: nil,

		Tags: map[string]string{
			"SHA":         sha,
			"Version":     ver,
			"Environment": pyroConfig.Environment(),
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
