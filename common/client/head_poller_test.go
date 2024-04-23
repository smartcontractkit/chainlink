package client

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func Test_Poller(t *testing.T) {
	// Mock polling function that returns a new value every time it's called
	var pollNumber int
	pollFunc := func() (Head, error) {
		pollNumber++
		h := head{
			BlockNumber:     int64(pollNumber),
			BlockDifficulty: big.NewInt(int64(pollNumber)),
		}
		return h.ToMockHead(t), nil
	}

	// data channel to receive updates from the poller
	channel := make(chan Head, 1)

	// Create poller and subscribe to receive data
	poller := NewHeadPoller[Head](time.Millisecond, pollFunc, channel, logger.Test(t))
	ctx := context.Background()

	require.NoError(t, poller.Start(ctx))
	defer poller.Unsubscribe()

	// Create goroutine to receive updates from the poller
	pollCount := 0
	pollMax := 50
	go func() {
		for ; pollCount < pollMax; pollCount++ {
			h := <-channel
			assert.Equal(t, int64(pollNumber), h.BlockNumber())
		}
	}()

	// Wait for a short duration to allow for some polling iterations
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, pollMax, pollCount)
}

// TODO: Test error in pollingFunc
// TODO: Test context cancellation
