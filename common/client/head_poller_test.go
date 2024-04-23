package client

import (
	"github.com/pkg/errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	poller := NewHeadPoller[Head](time.Millisecond, pollFunc, channel)

	require.NoError(t, poller.Start())
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

func Test_Poller_Error(t *testing.T) {
	// Mock polling function that returns an error every time it's called
	pollFunc := func() (Head, error) {
		return nil, errors.New("polling error")
	}

	// data channel to receive updates from the poller
	channel := make(chan Head, 1)

	// Create poller and subscribe to receive data
	poller := NewHeadPoller[Head](time.Millisecond, pollFunc, channel)

	require.NoError(t, poller.Start())
	defer poller.Unsubscribe()

	// Create goroutine to receive updates from the poller
	pollCount := 0
	pollMax := 50
	go func() {
		for ; pollCount < pollMax; pollCount++ {
			err := <-poller.Err()
			require.Error(t, err)
			require.Equal(t, "polling error", err.Error())
		}
	}()

	// Wait for a short duration to allow for some polling iterations
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, pollMax, pollCount)
}
