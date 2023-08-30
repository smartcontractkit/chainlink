package vrfcommon

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForceFulfillRateLimiter_FulfillmentPerformed(t *testing.T) {
	rl := NewForceFulfillRateLimiter()
	subId := big.NewInt(1)
	rl.FulfillmentPerformed(subId)
	require.Equal(t, 1, rl.NumFulfilled(subId))
}

func TestForceFulfillRateLimiter_NumFulfilled(t *testing.T) {
	rl := NewForceFulfillRateLimiter()
	expectedNumFulfillments := 10
	for i := 0; i < expectedNumFulfillments; i++ {
		rl.FulfillmentPerformed(big.NewInt(1))
	}
	require.Equal(t, expectedNumFulfillments, rl.NumFulfilled(big.NewInt(1)))
}

func TestForceFulfillRateLimiter_SetLatestHead(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		rl := NewForceFulfillRateLimiter()
		rl.SetLatestHead(100)
		require.Equal(t, uint64(100), rl.latestHead)
	})

	t.Run("with pruning", func(t *testing.T) {
		rl := NewForceFulfillRateLimiter()
		rl.SetLatestHead(100)
		rl.FulfillmentPerformed(big.NewInt(1))
		rl.SetLatestHead(100 + PruneInterval)
		require.Equal(t, 0, rl.NumFulfilled(big.NewInt(1)))
	})
}

func TestForceFulfillRateLimiter_ShouldFulfill(t *testing.T) {
	rl := NewForceFulfillRateLimiter()
	subId := big.NewInt(1)
	require.True(t, rl.ShouldFulfill(subId))
	// after MaxForceFulfillments have been done, ShouldFulfill should return false
	for i := 0; i < MaxForceFulfillments; i++ {
		rl.FulfillmentPerformed(subId)
	}
	require.False(t, rl.ShouldFulfill(subId))
}

func TestForceFulfillRateLimiter_prune(t *testing.T) {
	rl := NewForceFulfillRateLimiter()
	subId := big.NewInt(1)
	rl.FulfillmentPerformed(subId)
	rl.prune()
	require.Equal(t, 0, rl.NumFulfilled(subId))
	require.Empty(t, rl.forceFulfillsCount)
}
