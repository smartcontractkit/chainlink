package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

func TestRateLimiter_Simple(t *testing.T) {
	t.Parallel()

	rl := handlers.NewRateLimiter(3.0, 3, 1.0, 2)
	require.True(t, rl.Allow("user1"))
	require.True(t, rl.Allow("user2"))
	require.True(t, rl.Allow("user1"))
	require.False(t, rl.Allow("user1"))
	require.False(t, rl.Allow("user3"))
}
