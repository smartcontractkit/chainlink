package offchainreporting

import (
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"golang.org/x/time/rate"
)

type limiterFactory struct{}

var _ types.LimiterFactory = (*limiterFactory)(nil)

func (l *limiterFactory) MakeLimiter(numEventsPerSec uint) types.Limiter {
	return rate.NewLimiter(rate.Limit(numEventsPerSec), 1)
}
