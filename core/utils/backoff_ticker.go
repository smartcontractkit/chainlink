package utils

import (
	"sync"
	"time"

	"github.com/jpillora/backoff"
)

type timerFactory func(d time.Duration) *time.Timer

func newBackoffTicker(tf timerFactory, min, max time.Duration) BackoffTicker {
	c := make(chan time.Time, 1)
	return BackoffTicker{
		createTimer: tf,
		b: backoff.Backoff{
			Min: min,
			Max: max,
		},
		C:      c,
		chStop: make(chan struct{}),
	}
}

// BackoffTicker sends ticks with periods that increase over time, over a configured range.
type BackoffTicker struct {
	createTimer timerFactory
	b           backoff.Backoff
	timer       *time.Timer
	C           chan time.Time
	chStop      chan struct{}
	isRunning   bool
	sync.Mutex
}

// NewBackoffTicker returns a new BackoffTicker for the given range.
func NewBackoffTicker(min, max time.Duration) BackoffTicker {
	return newBackoffTicker(time.NewTimer, min, max)
}

// Start - Starts the ticker
// Returns true if the ticker was not running yet
func (t *BackoffTicker) Start() bool {
	t.Lock()
	defer t.Unlock()

	if t.isRunning {
		return false
	}

	// Reset the backoff
	t.b.Reset()
	go t.run()
	t.isRunning = true
	return true
}

// Stop stops the ticker. A ticker can be restarted by calling Start on a
// stopped ticker.
// Returns true if the ticker was actually stopped at this invocation (was previously running)
func (t *BackoffTicker) Stop() bool {
	t.Lock()
	defer t.Unlock()

	if !t.isRunning {
		return false
	}

	t.chStop <- struct{}{}
	t.timer = nil
	t.isRunning = false
	return true
}

func (t *BackoffTicker) run() {
	d := t.b.Duration()

	for {
		// Set up initial tick
		if t.timer == nil {
			t.timer = t.createTimer(d)
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

// Ticks returns the underlying channel.
func (t *BackoffTicker) Ticks() <-chan time.Time {
	return t.C
}

func (t *BackoffTicker) Bounds() (time.Duration, time.Duration) {
	return t.b.Min, t.b.Max
}
