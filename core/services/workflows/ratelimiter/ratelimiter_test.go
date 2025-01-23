package ratelimiter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRateLimiter(t *testing.T) {
	t.Parallel()

	config := Config{
		GlobalRPS:      3.0,
		GlobalBurst:    3,
		PerSenderRPS:   1.0,
		PerSenderBurst: 2,
	}
	rl, err := NewRateLimiter(config)
	require.NoError(t, err)
	allowUserSender, allowUserGlobal := rl.Allow("user1")
	require.True(t, allowUserSender && allowUserGlobal)
	allowUserSender, allowUserGlobal = rl.Allow("user2")
	require.True(t, allowUserSender && allowUserGlobal)
	allowUserSender, allowUserGlobal = rl.Allow("user1")
	require.True(t, allowUserSender && allowUserGlobal)
	allowUserSender, allowUserGlobal = rl.Allow("user1")
	require.False(t, allowUserSender && allowUserGlobal)
	allowUserSender, allowUserGlobal = rl.Allow("user3")
	require.False(t, allowUserSender && allowUserGlobal)
}
