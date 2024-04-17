package client

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type TestHead struct {
	blockNumber int64
}

var _ Head = &TestHead{}

func (th *TestHead) BlockNumber() int64 {
	return th.blockNumber
}

func (th *TestHead) BlockDifficulty() *big.Int {
	return nil
}

func (th *TestHead) IsValid() bool {
	return true
}

func Test_Polling_Transformer(t *testing.T) {
	t.Parallel()

	// Mock polling function that returns a new value every time it's called
	var lastBlockNumber int64
	pollFunc := func() (Head, error) {
		lastBlockNumber++
		return &TestHead{lastBlockNumber}, nil
	}

	pt := NewPollingTransformer(time.Millisecond, pollFunc, logger.Test(t))
	pt.StartPolling()
	defer pt.StopPolling()

	// Create a subscriber channel
	subscriber := make(chan Head)
	pt.Subscribe(subscriber)
	defer pt.Unsubscribe(subscriber)

	// Create a goroutine to receive updates from the subscriber
	pollCount := 0
	pollMax := 50
	go func() {
		for i := 0; i < pollMax; i++ {
			value := <-subscriber
			pollCount++
			require.Equal(t, int64(pollCount), value.BlockNumber())
		}
	}()

	// Wait for a short duration to allow for some polling iterations
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, pollMax, pollCount)
}
