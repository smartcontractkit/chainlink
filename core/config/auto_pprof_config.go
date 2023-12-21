package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type AutoPprof interface {
	BlockProfileRate() int
	CPUProfileRate() int
	Enabled() bool
	GatherDuration() sqlutil.Duration
	GatherTraceDuration() sqlutil.Duration
	GoroutineThreshold() int
	MaxProfileSize() utils.FileSize
	MemProfileRate() int
	MemThreshold() utils.FileSize
	MutexProfileFraction() int
	PollInterval() sqlutil.Duration
	ProfileRoot() string
}
