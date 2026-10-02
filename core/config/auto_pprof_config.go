package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type AutoPprof interface {
	BlockProfileRate() int
	CPUProfileRate() int
	Enabled() bool
	GatherDuration() config.Duration
	GatherTraceDuration() config.Duration
	GoroutineThreshold() int
	MaxProfileSize() utils.FileSize
	MemProfileRate() int
	MemThreshold() utils.FileSize
	MutexProfileFraction() int
	PollInterval() config.Duration
	ProfileRoot() string
}
