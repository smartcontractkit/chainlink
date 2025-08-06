package contracts

import (
	"testing"

	"github.com/stretchr/testify/require"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestOperationDefinition(t *testing.T) {
	t.Parallel()

	t.Run("operations are defined", func(t *testing.T) {
		require.NotNil(t, DeployWorkflowRegistryV2Op)
		require.NotNil(t, DeployWorkflowRegistryV2MultiChainOp)
	})
}

func TestOperationInputOutput(t *testing.T) {
	t.Parallel()

	t.Run("DeployWorkflowRegistryV2Input validation", func(t *testing.T) {
		input := DeployWorkflowRegistryV2Input{
			ChainSelector: 1,
		}
		require.Equal(t, uint64(1), input.ChainSelector)
	})

	t.Run("DeployWorkflowRegistryV2Request validation", func(t *testing.T) {
		req := keystone_changeset.DeployWorkflowRegistryV2Request{
			ChainSelectors: []uint64{1, 2, 3},
		}
		require.Len(t, req.ChainSelectors, 3)
		require.Equal(t, uint64(1), req.ChainSelectors[0])
	})

	t.Run("empty chain selectors", func(t *testing.T) {
		req := keystone_changeset.DeployWorkflowRegistryV2Request{
			ChainSelectors: []uint64{},
		}
		require.Len(t, req.ChainSelectors, 0)
	})
}
