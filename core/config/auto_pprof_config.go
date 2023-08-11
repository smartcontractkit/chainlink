package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type AutoPprof interface {
	BlockProfileRate() int
	CPUProfileRate() int
	Enabled() bool
	GatherDuration() models.Duration
	GatherTraceDuration() models.Duration
	GoroutineThreshold() int
	MaxProfileSize() utils.FileSize
	MemProfileRate() int
	MemThreshold() utils.FileSize
	MutexProfileFraction() int
	PollInterval() models.Duration
	ProfileRoot() string
}
