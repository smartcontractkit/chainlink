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
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func Test_Poller(t *testing.T) {
	lggr := logger.Test(t)

	t.Run("Test multiple start", func(t *testing.T) {
		pollFunc := func(ctx context.Context) (Head, error) {
			return nil, nil
		}

		channel := make(chan Head, 1)
		defer close(channel)

		poller := NewPoller[Head](time.Millisecond, pollFunc, time.Second, channel, lggr)
		err := poller.Start()
		require.NoError(t, err)

		err = poller.Start()
		require.Error(t, err)
		poller.Unsubscribe()
	})

	t.Run("Test polling for heads", func(t *testing.T) {
		// Mock polling function that returns a new value every time it's called
		var pollNumber int
		pollLock := sync.Mutex{}
		pollFunc := func(ctx context.Context) (Head, error) {
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
		poller := NewPoller[Head](time.Millisecond, pollFunc, time.Second, channel, lggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Receive updates from the poller
		pollCount := 0
		pollMax := 50
		for ; pollCount < pollMax; pollCount++ {
			h := <-channel
			assert.Equal(t, int64(pollCount+1), h.BlockNumber())
		}
	})

	t.Run("Test polling errors", func(t *testing.T) {
		// Mock polling function that returns an error
		var pollNumber int
		pollLock := sync.Mutex{}
		pollFunc := func(ctx context.Context) (Head, error) {
			pollLock.Lock()
			defer pollLock.Unlock()
			pollNumber++
			return nil, fmt.Errorf("polling error %d", pollNumber)
		}

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		olggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)

		// Create poller and subscribe to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, time.Second, channel, olggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Ensure that all errors were logged as expected
		logsSeen := func() bool {
			for pollCount := 0; pollCount < 50; pollCount++ {
				numLogs := observedLogs.FilterMessage(fmt.Sprintf("polling error: polling error %d", pollCount+1)).Len()
				if numLogs != 1 {
					return false
				}
			}
			return true
		}
		require.Eventually(t, logsSeen, time.Second, time.Millisecond)
	})

	t.Run("Test polling timeout", func(t *testing.T) {
		pollFunc := func(ctx context.Context) (Head, error) {
			if <-ctx.Done(); true {
				return nil, ctx.Err()
			}
			return nil, nil
		}

		// Set instant timeout
		pollingTimeout := time.Duration(0)

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		olggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)

		// Create poller and subscribe to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, pollingTimeout, channel, olggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Ensure that timeout errors were logged as expected
		logsSeen := func() bool {
			return observedLogs.FilterMessage("polling error: context deadline exceeded").Len() >= 1
		}
		require.Eventually(t, logsSeen, time.Second, time.Millisecond)
	})

	t.Run("Test unsubscribe during polling", func(t *testing.T) {
		wait := make(chan struct{})
		pollFunc := func(ctx context.Context) (Head, error) {
			close(wait)
			// Block in polling function until context is cancelled
			if <-ctx.Done(); true {
				return nil, ctx.Err()
			}
			return nil, nil
		}

		// Set long timeout
		pollingTimeout := time.Minute

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)
		defer close(channel)

		olggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)

		// Create poller and subscribe to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, pollingTimeout, channel, olggr)
		require.NoError(t, poller.Start())

		// Unsubscribe while blocked in polling function
		<-wait
		poller.Unsubscribe()

		// Ensure error was logged
		logsSeen := func() bool {
			return observedLogs.FilterMessage("polling error: context canceled").Len() >= 1
		}
		require.Eventually(t, logsSeen, time.Second, time.Millisecond)
	})
}

func Test_Poller_Unsubscribe(t *testing.T) {
	lggr := logger.Test(t)
	pollFunc := func(ctx context.Context) (Head, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			h := head{
				BlockNumber:     0,
				BlockDifficulty: big.NewInt(0),
			}
			return h.ToMockHead(t), nil
		}
	}

	t.Run("Test multiple unsubscribe", func(t *testing.T) {
		channel := make(chan Head, 1)
		poller := NewPoller[Head](time.Millisecond, pollFunc, time.Second, channel, lggr)
		err := poller.Start()
		require.NoError(t, err)

		<-channel
		poller.Unsubscribe()
		poller.Unsubscribe()
	})

	t.Run("Test unsubscribe with closed channel", func(t *testing.T) {
		channel := make(chan Head, 1)
		poller := NewPoller[Head](time.Millisecond, pollFunc, time.Second, channel, lggr)
		err := poller.Start()
		require.NoError(t, err)

		<-channel
		close(channel)
		poller.Unsubscribe()
	})
}
