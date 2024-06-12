package timeutil

import (
	"testing"
	"time"
)

func TestJitterPct(t *testing.T) {
	for _, tt := range []struct {
		pct      JitterPct
		dur      time.Duration
		from, to time.Duration
	}{
		{0.1, 0, 0, 0},
		{0.1, time.Second, 900 * time.Millisecond, 1100 * time.Millisecond},
		{0.1, time.Minute, 54 * time.Second, 66 * time.Second},
		{0.1, 24 * time.Hour, 21*time.Hour + 36*time.Minute, 26*time.Hour + 24*time.Minute},
	} {
		t.Run(tt.dur.String(), func(t *testing.T) {
			for i := 0; i < 100; i++ {
				got := tt.pct.Apply(tt.dur)
				t.Logf("%d: %s", i, got)
				if got < tt.from || got > tt.to {
					t.Errorf("expected duration %s with jitter to be between (%s, %s) but got: %s", tt.dur, tt.from, tt.to, got)
				}
			}
		})
	}
}
