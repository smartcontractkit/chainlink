package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSmokeState(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv := NewMemoryEnvironmentWithJobsAndContracts(t, lggr, memory.MemoryEnvironmentConfig{
		Chains:             3,
		Nodes:              4,
		Bootstraps:         1,
		NumOfUsersPerChain: 1,
	}, nil)
	state, err := LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	_, err = state.View(tenv.Env.AllChainSelectors())
	require.NoError(t, err)
}
