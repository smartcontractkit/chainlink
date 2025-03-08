package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRateLimiter_PerWorkflow(t *testing.T) {
	t.Parallel()

	config := RateLimiterConfig{
		GlobalRPS:        3.0,
		GlobalBurst:      3,
		PerWorkflowRPS:   1.0,
		PerWorkflowBurst: 2,
	}
	rl, err := NewRateLimiter(config)
	require.NoError(t, err)
	require.True(t, rl.Allow("user1"), "workflowID1")
	require.True(t, rl.Allow("user2"), "workflowID2")
	require.True(t, rl.Allow("user3"), "workflowID1")
	require.False(t, rl.Allow("user4"), "workflowID1")
	require.False(t, rl.Allow("user5"), "workflowID3")
}
