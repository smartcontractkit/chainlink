package utils

import (
	"sync"
	"time"

	"github.com/jpillora/backoff"
)

type BackoffTicker struct {
	b         backoff.Backoff
	timer     *time.Timer
	C         chan time.Time
	chStop    chan struct{}
	isRunning bool
	sync.Mutex
}

func NewBackoffTicker(min, max time.Duration) BackoffTicker {
	c := make(chan time.Time, 1)
	return BackoffTicker{
		b: backoff.Backoff{
			Min: min,
			Max: max,
		},
		C:      c,
		chStop: make(chan struct{}),
	}
}

// Starts the ticker
func (t *BackoffTicker) Start() {
	t.Lock()
	defer t.Unlock()

	if t.isRunning {
		return
	}

	// Reset the backoff
	t.b.Reset()
	go t.run()
	t.isRunning = true
}

// Stop stops the ticker. A ticker can be restarted by calling Start on a
// stopped ticker.
func (t *BackoffTicker) Stop() {
	t.Lock()
	defer t.Unlock()

	if !t.isRunning {
		return
	}

	t.chStop <- struct{}{}
	t.timer = nil
	t.isRunning = false
}

func (t *BackoffTicker) run() {
	d := t.b.Duration()

	for {
		// Set up initial tick
		if t.timer == nil {
			t.timer = time.NewTimer(d)
		}

		select {
		case tickTime := <-t.timer.C:
			t.C <- tickTime
			t.timer.Reset(t.b.Duration())

			continue
		case <-t.chStop:
			return
		}
	}
}

func (t *BackoffTicker) Ticks() <-chan time.Time {
	return t.C
}
