package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

func TestRateLimiter_PerSender(t *testing.T) {
	t.Parallel()

	config := common.RateLimiterConfig{
		GlobalRPS:      3.0,
		GlobalBurst:    3,
		PerSenderRPS:   1.0,
		PerSenderBurst: 2,
	}
	rl, err := common.NewRateLimiter(config)
	require.NoError(t, err)
	require.True(t, rl.Allow("user1"))
	require.True(t, rl.Allow("user2"))
	require.True(t, rl.Allow("user1"))
	require.False(t, rl.Allow("user1"))
	require.False(t, rl.Allow("user3"))
}
