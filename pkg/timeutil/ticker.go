package timeutil

import (
	"time"
)

// Ticker is like time.Ticker, but with a variable period.
type Ticker struct {
	C     <-chan time.Time
	stop  chan struct{}
	reset chan struct{}
}

// NewTicker returns a started Ticker which calls nextDur for each period.
// Ticker.Stop should be called to prevent goroutine leaks.
func NewTicker(nextDur func() time.Duration) *Ticker {
	c := make(chan time.Time) // unbuffered so we block and delay if not being handled
	t := Ticker{C: c, stop: make(chan struct{}), reset: make(chan struct{})}
	go t.run(c, nextDur)
	return &t
}

// Stop permanently stops the Ticker. It cannot be Reset.
func (t *Ticker) Stop() { close(t.stop) }

func (t *Ticker) run(c chan<- time.Time, nextDur func() time.Duration) {
	for {
		timer := time.NewTimer(nextDur())
		select {
		case <-t.stop:
			timer.Stop()
			return

		case <-t.reset:
			timer.Stop()

		case <-timer.C:
			timer.Stop()
			select {
			case <-t.stop:
				return
			case c <- time.Now():
			case <-t.reset:
			}
		}
	}
}

// Reset starts a new period.
func (t *Ticker) Reset() {
	select {
	case <-t.stop:
	case t.reset <- struct{}{}:
	default:
		// unnecessary
		return
	}
}
