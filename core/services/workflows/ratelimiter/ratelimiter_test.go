package ratelimiter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

func TestRateLimiter(t *testing.T) {
	t.Parallel()

	config := Config{
		GlobalRPS:      3.0,
		GlobalBurst:    3,
		PerSenderRPS:   1.0,
		PerSenderBurst: 2,
	}
	rl, err := NewRateLimiter(config, limits.Factory{Logger: logger.Test(t)})
	require.NoError(t, err)
	ctx1 := contexts.WithCRE(t.Context(), contexts.CRE{Owner: "user1", Workflow: "wf-1"})
	require.True(t, rl.Allow(ctx1))
	require.True(t, rl.Allow(contexts.WithCRE(t.Context(), contexts.CRE{Owner: "user2", Workflow: "wf-2"})))
	require.True(t, rl.Allow(ctx1))
	require.False(t, rl.Allow(ctx1))
	require.False(t, rl.Allow(contexts.WithCRE(t.Context(), contexts.CRE{Owner: "user3", Workflow: "wf-3"})))

}
