package client

import (
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
	poller := NewPoller[Head](time.Millisecond, pollFunc, channel, logger.Test(t))
	require.NoError(t, poller.Subscribe())

	// Create goroutine to receive updates from the subscriber
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
