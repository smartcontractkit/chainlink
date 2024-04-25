package client

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func Test_Poller(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	t.Run("Test polling for heads", func(t *testing.T) {
		// Mock polling function that returns a new value every time it's called
		var pollNumber int
		pollLock := sync.Mutex{}
		pollFunc := func(ctx context.Context, args ...interface{}) (Head, error) {
			pollLock.Lock()
			defer pollLock.Unlock()
			pollNumber++
			h := head{
				BlockNumber:     int64(pollNumber),
				BlockDifficulty: big.NewInt(int64(pollNumber)),
			}
			return h.ToMockHead(t), nil
		}

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		// Create poller and start to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, nil, channel, &lggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Monitor error channel
		done := make(chan struct{})
		defer close(done)
		monitorPollingErrors(t, poller.Err(), done)

		// Receive updates from the poller
		func() {
			pollCount := 0
			pollMax := 50
			for ; pollCount < pollMax; pollCount++ {
				h := <-channel
				assert.Equal(t, int64(pollCount+1), h.BlockNumber())
			}
		}()
	})

	t.Run("Test polling errors", func(t *testing.T) {
		// Mock polling function that returns an error
		var pollNumber int
		pollLock := sync.Mutex{}
		pollFunc := func(ctx context.Context, args ...interface{}) (Head, error) {
			pollLock.Lock()
			defer pollLock.Unlock()
			pollNumber++
			return nil, fmt.Errorf("polling error %d", pollNumber)
		}

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		// Create poller and subscribe to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, nil, channel, &lggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Create goroutine to receive updates from the poller
		func() {
			pollCount := 0
			pollMax := 50
			for ; pollCount < pollMax; pollCount++ {
				select {
				case <-channel:
					require.Fail(t, "should not receive any data")
				case err := <-poller.Err():
					require.Error(t, err)
					require.Equal(t, fmt.Sprintf("polling error %d", pollCount+1), err.Error())
				}
			}
		}()
	})

	t.Run("Test polling timeout", func(t *testing.T) {
		pollFunc := func(ctx context.Context, args ...interface{}) (Head, error) {
			time.Sleep(10 * time.Millisecond)
			return nil, nil
		}

		// Set instant timeout
		pollingTimeout := time.Duration(0)

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		// Create poller and subscribe to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, &pollingTimeout, channel, &lggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Create goroutine to receive updates from the poller
		func() {
			err := <-poller.Err()
			require.Error(t, err)
			require.Equal(t, "polling timeout exceeded", err.Error())
		}()
	})

	t.Run("Test polling with args", func(t *testing.T) {
		pollFunc := func(ctx context.Context, args ...interface{}) (Head, error) {
			require.Equal(t, args[0], "arg1")
			require.Equal(t, args[1], "arg2")
			require.Equal(t, args[2], "arg3")
			return nil, nil
		}

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		// Create poller and subscribe to receive data
		args := []interface{}{"arg1", "arg2", "arg3"}
		poller := NewPoller[Head](time.Millisecond, pollFunc, nil, channel, &lggr, args...)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Ensure no errors are received
		done := make(chan struct{})
		defer close(done)
		monitorPollingErrors(t, poller.Err(), done)

		// Create goroutine to receive updates from the poller
		func() {
			h := <-channel
			require.Nil(t, h)
		}()
	})

	t.Run("Test panic in polling function", func(t *testing.T) {
		pollFunc := func(ctx context.Context, args ...interface{}) (Head, error) {
			panic("panic test")
		}

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		// Create poller and subscribe to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, nil, channel, &lggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Create goroutine to receive updates from the poller
		func() {
			err := <-poller.Err()
			require.Equal(t, "panic: panic test", err.Error())
		}()
	})
}

func Test_Poller_Unsubscribe(t *testing.T) {
	t.Run("Test multiple unsubscribe", func(t *testing.T) {
		poller := NewPoller[Head](time.Millisecond, nil, nil, nil, nil)
		err := poller.Start()
		require.NoError(t, err)
		poller.Unsubscribe()
		poller.Unsubscribe()
	})

	t.Run("Test unsubscribe with no subscribers", func(t *testing.T) {
		poller := NewPoller[Head](time.Millisecond, nil, nil, nil, nil)
		poller.Unsubscribe()
	})
}

// monitorPollingErrors fails the test if an error is received on the error channel
func monitorPollingErrors(t *testing.T, errCh <-chan error, done <-chan struct{}) {
	go func() {
		select {
		case err := <-errCh:
			require.NoError(t, err)
		case <-done:
			return
		}
	}()
}
