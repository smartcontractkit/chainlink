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
	workflowAllow1, globalAllow1 := rl.Allow("workflowID1")
	require.True(t, workflowAllow1)
	require.True(t, globalAllow1)
	workflowAllow2, globalAllow2 := rl.Allow("workflowID2")
	require.True(t, workflowAllow2)
	require.True(t, globalAllow2)
	workflowAllow3, globalAllow3 := rl.Allow("workflowID1")
	require.True(t, workflowAllow3)
	require.True(t, globalAllow3)
	workflowAllow4, globalAllow4 := rl.Allow("workflowID1")
	require.False(t, workflowAllow4)
	require.False(t, globalAllow4)
	workflowAllow5, globalAllow5 := rl.Allow("workflowID3")
	require.True(t, workflowAllow5)
	require.False(t, globalAllow5)
}
