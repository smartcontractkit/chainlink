package client

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// TODO: Fix race conditions in tests!

func Test_Poller(t *testing.T) {
	pollingTimeout := 10 * time.Millisecond
	lggr, err := logger.New()
	require.NoError(t, err)

	t.Run("Test polling for heads", func(t *testing.T) {
		// Mock polling function that returns a new value every time it's called
		var pollNumber int
		pollFunc := func(ctx context.Context) (Head, error) {
			pollNumber++
			h := head{
				BlockNumber:     int64(pollNumber),
				BlockDifficulty: big.NewInt(int64(pollNumber)),
			}
			return h.ToMockHead(t), nil
		}

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)

		// Create poller and start to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, &pollingTimeout, channel, &lggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Create goroutine to receive updates from the poller
		done := make(chan struct{})
		go func() {
			pollCount := 0
			pollMax := 50
			for ; pollCount < pollMax; pollCount++ {
				h := <-channel
				assert.Equal(t, int64(pollNumber), h.BlockNumber())
			}
			close(done)
		}()
		<-done
	})

	t.Run("Test polling errors", func(t *testing.T) {
		// Mock polling function that returns an error
		pollFunc := func(ctx context.Context) (Head, error) {
			return nil, errors.New("polling error")
		}

		// data channel to receive updates from the poller
		channel := make(chan Head, 1)

		// Create poller and subscribe to receive data
		poller := NewPoller[Head](time.Millisecond, pollFunc, &pollingTimeout, channel, &lggr)
		require.NoError(t, poller.Start())
		defer poller.Unsubscribe()

		// Create goroutine to receive updates from the poller
		done := make(chan struct{})
		go func() {
			pollCount := 0
			pollMax := 50
			for ; pollCount < pollMax; pollCount++ {
				select {
				case <-channel:
					require.Fail(t, "should not receive any data")
				case err := <-poller.Err():
					require.Error(t, err)
					require.Equal(t, "polling error", err.Error())
				}
			}
			close(done)
		}()
		<-done
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
		// TODO: Add test case that ensures Unsubscribe exits even if no one is listening
		poller := NewPoller[Head](time.Millisecond, nil, nil, nil, nil)
		err := poller.Start()
		require.NoError(t, err)
		poller.Unsubscribe()
	})
}
