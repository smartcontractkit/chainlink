package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type AutoPprof interface {
	AutoPprofBlockProfileRate() int
	AutoPprofCPUProfileRate() int
	AutoPprofGatherDuration() models.Duration
	AutoPprofGatherTraceDuration() models.Duration
	AutoPprofGoroutineThreshold() int
	AutoPprofMaxProfileSize() utils.FileSize
	AutoPprofMemProfileRate() int
	AutoPprofMemThreshold() utils.FileSize
	AutoPprofMutexProfileFraction() int
	AutoPprofPollInterval() models.Duration
	AutoPprofProfileRoot() string
}
