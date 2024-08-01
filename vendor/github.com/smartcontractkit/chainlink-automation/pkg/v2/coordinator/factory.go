package coordinator

import (
	"log"
	"time"

	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"
	"github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
)

// CoordinatorFactory provides a single method to create a new coordinator
type CoordinatorFactory struct {
	Logger     *log.Logger
	Encoder    Encoder
	Logs       LogProvider
	CacheClean time.Duration
}

// NewCoordinator returns a new coordinator with provided dependencies and
// config. The new coordinator is not automatically started.
func (f *CoordinatorFactory) NewCoordinator(c config.OffchainConfig) (ocr2keepers.Coordinator, error) {
	return NewReportCoordinator(
		time.Duration(c.PerformLockoutWindow)*time.Millisecond,
		f.CacheClean,
		f.Logs,
		c.MinConfirmations,
		f.Logger,
		f.Encoder,
	), nil
}
