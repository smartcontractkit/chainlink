package logger

import (
	"os"
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

	// Increase memory profiling sample rate for better granularity
	// Default is 512KB (524288 bytes) per sample
	// runtime.MemProfileRate = 512 * 1024 // 512KB per sample

	sha, ver := static.Short()

	return pyroscope.Start(pyroscope.Config{
		// Maybe configurable to identify the specific NOP - TBD
		ApplicationName: "chainlink-node",

		ServerAddress: pyroConfig.ServerAddress(),
		AuthToken:     pyroConfig.AuthToken(),

		// We disable logging the profiling info, it will be in the Pyroscope instance anyways...
		Logger: nil,

		Tags: func() map[string]string {
			hostname, _ := os.Hostname()
			return map[string]string{
				"SHA":         sha,
				"Version":     ver,
				"Environment": pyroConfig.Environment(),
				"hostname":    hostname, // set hostname, so we can distinguish between nodes in the same environment
			}
		}(),

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
