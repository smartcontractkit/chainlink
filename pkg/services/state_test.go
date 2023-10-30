package services

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateMachine_StartOnce_StopOnce(t *testing.T) {
	t.Parallel()

	var sm StateMachine

	ch := make(chan int, 3)

	ready := make(chan bool)

	go func() {
		assert.NoError(t, sm.StartOnce("slow service", func() (err error) {
			ch <- 1
			ready <- true
			<-time.After(time.Millisecond * 500) // wait for StopOnce to happen
			ch <- 2

			return nil
		}))
	}()

	go func() {
		<-ready // try stopping halfway through startup
		assert.NoError(t, sm.StopOnce("slow service", func() (err error) {
			ch <- 3

			return nil
		}))
	}()

	require.Equal(t, 1, <-ch)
	require.Equal(t, 2, <-ch)
	require.Equal(t, 3, <-ch)
}

func TestStateMachine_MultipleStartNoBlock(t *testing.T) {
	t.Parallel()

	var sm StateMachine

	ch := make(chan int, 3)

	ready := make(chan bool)
	next := make(chan bool)

	go func() {
		ch <- 1
		assert.NoError(t, sm.StartOnce("slow service", func() (err error) {
			ready <- true
			<-next // continue after the other StartOnce call fails

			return nil
		}))
		<-next
		ch <- 2

	}()

	go func() {
		<-ready // try starting halfway through startup
		assert.Error(t, sm.StartOnce("slow service", func() (err error) {
			return nil
		}))
		next <- true
		ch <- 3
		next <- true

	}()

	require.Equal(t, 1, <-ch)
	require.Equal(t, 3, <-ch) // 3 arrives before 2 because it returns immediately
	require.Equal(t, 2, <-ch)
}

func TestStateMachine(t *testing.T) {
	t.Parallel()

	var callsCount atomic.Int32
	incCount := func() {
		callsCount.Add(1)
	}

	var s StateMachine
	ok := s.IfStarted(incCount)
	assert.False(t, ok)
	ok = s.IfNotStopped(incCount)
	assert.True(t, ok)
	assert.Equal(t, int32(1), callsCount.Load())

	err := s.StartOnce("foo", func() error { return nil })
	assert.NoError(t, err)

	assert.True(t, s.IfStarted(incCount))
	assert.Equal(t, int32(2), callsCount.Load())

	err = s.StopOnce("foo", func() error { return nil })
	assert.NoError(t, err)
	ok = s.IfNotStopped(incCount)
	assert.False(t, ok)
	assert.Equal(t, int32(2), callsCount.Load())
}

func TestStateMachine_StartErrors(t *testing.T) {
	var s StateMachine

	err := s.StartOnce("foo", func() error { return errors.New("foo") })
	assert.Error(t, err)

	var callsCount atomic.Int32
	incCount := func() {
		callsCount.Add(1)
	}

	assert.False(t, s.IfStarted(incCount))
	assert.Equal(t, int32(0), callsCount.Load())

	err = s.StartOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo has already been started once")
	err = s.StopOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo cannot be stopped from this state; state=StartFailed")

	assert.Equal(t, stateStartFailed, s.loadState())
}

func TestStateMachine_StopErrors(t *testing.T) {
	var s StateMachine

	err := s.StartOnce("foo", func() error { return nil })
	require.NoError(t, err)

	var callsCount atomic.Int32
	incCount := func() {
		callsCount.Add(1)
	}

	err = s.StopOnce("foo", func() error { return errors.New("explodey mcsplode") })
	assert.Error(t, err)

	assert.False(t, s.IfStarted(incCount))
	assert.Equal(t, int32(0), callsCount.Load())
	assert.True(t, s.IfNotStopped(incCount))
	assert.Equal(t, int32(1), callsCount.Load())

	err = s.StartOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo has already been started once")
	err = s.StopOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo cannot be stopped from this state; state=StopFailed")

	assert.Equal(t, stateStopFailed, s.loadState())
}

func TestStateMachine_Ready_Healthy(t *testing.T) {
	t.Parallel()

	var s StateMachine
	assert.Error(t, s.Ready())
	assert.Error(t, s.Healthy())

	err := s.StartOnce("foo", func() error { return nil })
	assert.NoError(t, err)
	assert.NoError(t, s.Ready())
	assert.NoError(t, s.Healthy())
}
