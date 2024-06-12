package services

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
)

// DefaultJitter is +/-10%
const DefaultJitter timeutil.JitterPct = 0.1

// NewTicker returns a new timeutil.Ticker configured to:
// - fire the first tick immediately
// - apply DefaultJitter to each period
func NewTicker(period time.Duration) *timeutil.Ticker {
	return TickerConfig{JitterPct: DefaultJitter}.NewTicker(period)
}

type TickerConfig struct {
	// Initial delay before the first tick.
	Initial time.Duration
	// JitterPct to apply to each period.
	JitterPct timeutil.JitterPct
}

func (c TickerConfig) NewTicker(period time.Duration) *timeutil.Ticker {
	first := true
	return timeutil.NewTicker(func() time.Duration {
		if first {
			first = false
			return c.Initial
		}
		p := period
		if c.JitterPct != 0.0 {
			p = c.JitterPct.Apply(p)
		}
		return p
	})
}
