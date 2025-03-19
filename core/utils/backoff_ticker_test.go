package utils

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackoffTicker_Bounds(t *testing.T) {
	t.Parallel()

	bt := NewBackoffTicker(1*time.Millisecond, 2*time.Second)
	min, max := bt.Bounds()
	assert.Equal(t, 1*time.Millisecond, min)
	assert.Equal(t, 2*time.Second, max)
}

func TestBackoffTicker_StartTwice(t *testing.T) {
	t.Parallel()

	bt := NewBackoffTicker(1*time.Second, 10*time.Second)
	defer bt.Stop()

	ok := bt.Start()
	assert.True(t, ok)

	ok = bt.Start()
	assert.False(t, ok)
}

func TestBackoffTicker_StopTwice(t *testing.T) {
	t.Parallel()

	bt := NewBackoffTicker(1*time.Second, 10*time.Second)
	ok := bt.Start()
	assert.True(t, ok)

	ok = bt.Stop()
	assert.True(t, ok)

	ok = bt.Stop()
	assert.False(t, ok)
}

func TestBackoffTicker_NoTicksAfterStop(t *testing.T) {
	t.Parallel()

	min := 100 * time.Millisecond
	max := 5 * time.Second

	chTime := make(chan time.Time, 1)
	defer close(chTime)

	newFakeTimer := func(d time.Duration) *time.Timer {
		assert.Equal(t, min, d)
		realTimer := time.NewTimer(max)
		realTimer.C = chTime
		return realTimer
	}

	bt := newBackoffTicker(newFakeTimer, min, max)

	ok := bt.Start()
	assert.True(t, ok)

	ok = bt.Stop()
	assert.True(t, ok)

	chTime <- time.Now()

	select {
	case <-time.After(2 * min):
	case <-bt.Ticks():
		assert.FailNow(t, "received a tick after Stop()")
	}
}

func TestBackoffTicker_Ticks(t *testing.T) {
	t.Parallel()

	min := 100 * time.Millisecond
	max := 5 * time.Second

	chTime := make(chan time.Time)
	defer close(chTime)

	newFakeTimer := func(d time.Duration) *time.Timer {
		assert.Equal(t, min, d)
		realTimer := time.NewTimer(max)
		realTimer.C = chTime
		return realTimer
	}

	bt := newBackoffTicker(newFakeTimer, min, max)

	ok := bt.Start()
	assert.True(t, ok)
	defer bt.Stop()

	t1 := time.Now()
	t2 := t1.Add(1 * time.Second)
	t3 := t2.Add(1 * time.Second)
	times := []time.Time{t1, t2, t3}

	go func() {
		for _, tm := range times {
			chTime <- tm
		}
	}()

	for _, tm := range times {
		tick := <-bt.Ticks()
		assert.Equal(t, tm, tick)
	}

	select {
	case <-time.After(2 * min):
	case <-bt.Ticks():
		assert.FailNow(t, "received an unexpected tick")
	}
}

func TestBackoffTicker_Restart(t *testing.T) {
	t.Parallel()

	min := 1 * time.Second
	max := 10 * time.Second

	var newTimerCount atomic.Int32

	newFakeTimer := func(d time.Duration) *time.Timer {
		newTimerCount.Add(1)
		assert.Equal(t, min, d)
		return time.NewTimer(max)
	}

	bt := newBackoffTicker(newFakeTimer, min, max)

	ok := bt.Start()
	assert.True(t, ok)

	ok = bt.Stop()
	assert.True(t, ok)

	ok = bt.Start()
	assert.True(t, ok)
	defer bt.Stop()

	assert.Eventually(t, func() bool {
		return newTimerCount.Load() == 2
	}, min, min/100, "expected timer factory to be triggered twice")
}
