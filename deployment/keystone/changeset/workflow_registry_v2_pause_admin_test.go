package changeset

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func setupPauseAdminTest(t *testing.T) (cldf.Environment, uint64, cldf_evm.Chain, common.Address) {
	lggr, err := logger.New()
	require.NoError(t, err)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  4,
	})

	var chainSelector uint64
	var chain cldf_evm.Chain
	for sel, ch := range env.BlockChains.EVMChains() {
		chainSelector = sel
		chain = ch
		break
	}

	// Deploy workflow registry and update environment with DataStore
	ds := datastore.NewMemoryDataStore()
	resp, err := DeployWorkflowRegistryV2ToChain(context.Background(), chain, ds)
	require.NoError(t, err)

	// Update environment with DataStore for admin functions to find the registry
	env.DataStore = ds.Seal()

	return env, chainSelector, chain, resp.Address
}

func TestAdminPauseWorkflow(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupPauseAdminTest(t)

	t.Run("successful workflow pause", func(t *testing.T) {
		workflowID := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

		output, err := AdminPauseWorkflow(env, AdminPauseWorkflowRequest{
			ChainSelector: chainSelector,
			WorkflowID:    workflowID,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		workflowID := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

		_, err := AdminPauseWorkflow(env, AdminPauseWorkflowRequest{
			ChainSelector: 999999,
			WorkflowID:    workflowID,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})

	t.Run("zero workflow ID", func(t *testing.T) {
		workflowID := [32]byte{} // All zeros

		output, err := AdminPauseWorkflow(env, AdminPauseWorkflowRequest{
			ChainSelector: chainSelector,
			WorkflowID:    workflowID,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})
}

func TestAdminBatchPauseWorkflows(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupPauseAdminTest(t)

	t.Run("successful batch pause", func(t *testing.T) {
		workflowIDs := [][32]byte{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		}

		output, err := AdminBatchPauseWorkflows(env, AdminBatchPauseWorkflowsRequest{
			ChainSelector: chainSelector,
			WorkflowIDs:   workflowIDs,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("empty workflow IDs list", func(t *testing.T) {
		output, err := AdminBatchPauseWorkflows(env, AdminBatchPauseWorkflowsRequest{
			ChainSelector: chainSelector,
			WorkflowIDs:   [][32]byte{},
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("single workflow ID in batch", func(t *testing.T) {
		workflowIDs := [][32]byte{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		}

		output, err := AdminBatchPauseWorkflows(env, AdminBatchPauseWorkflowsRequest{
			ChainSelector: chainSelector,
			WorkflowIDs:   workflowIDs,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		workflowIDs := [][32]byte{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		}

		_, err := AdminBatchPauseWorkflows(env, AdminBatchPauseWorkflowsRequest{
			ChainSelector: 999999,
			WorkflowIDs:   workflowIDs,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestAdminPauseAllByOwner(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupPauseAdminTest(t)

	t.Run("successful pause all by owner", func(t *testing.T) {
		owner := common.HexToAddress("0x1234567890123456789012345678901234567890")

		output, err := AdminPauseAllByOwner(env, AdminPauseAllByOwnerRequest{
			ChainSelector: chainSelector,
			Owner:         owner,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("zero address owner", func(t *testing.T) {
		owner := common.Address{} // Zero address

		output, err := AdminPauseAllByOwner(env, AdminPauseAllByOwnerRequest{
			ChainSelector: chainSelector,
			Owner:         owner,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		owner := common.HexToAddress("0x1234567890123456789012345678901234567890")

		_, err := AdminPauseAllByOwner(env, AdminPauseAllByOwnerRequest{
			ChainSelector: 999999,
			Owner:         owner,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestAdminPauseAllByDON(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupPauseAdminTest(t)

	t.Run("successful pause all by DON", func(t *testing.T) {
		donFamily := "test_don_family"

		output, err := AdminPauseAllByDON(env, AdminPauseAllByDONRequest{
			ChainSelector: chainSelector,
			DONFamily:     donFamily,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("empty DON family", func(t *testing.T) {
		output, err := AdminPauseAllByDON(env, AdminPauseAllByDONRequest{
			ChainSelector: chainSelector,
			DONFamily:     "",
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("long DON family name", func(t *testing.T) {
		donFamily := "very_long_don_family_name_that_might_test_limits_of_string_handling"

		output, err := AdminPauseAllByDON(env, AdminPauseAllByDONRequest{
			ChainSelector: chainSelector,
			DONFamily:     donFamily,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		donFamily := "test_don_family"

		_, err := AdminPauseAllByDON(env, AdminPauseAllByDONRequest{
			ChainSelector: 999999,
			DONFamily:     donFamily,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

// TestPauseAdminSequence tests a sequence of pause operations
func TestPauseAdminSequence(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupPauseAdminTest(t)

	// Test a realistic sequence of pause operations
	t.Run("pause sequence", func(t *testing.T) {
		// First, pause a specific workflow
		workflowID := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
		_, err := AdminPauseWorkflow(env, AdminPauseWorkflowRequest{
			ChainSelector: chainSelector,
			WorkflowID:    workflowID,
		})
		require.NoError(t, err)

		// Then pause multiple workflows in batch
		workflowIDs := [][32]byte{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		}
		_, err = AdminBatchPauseWorkflows(env, AdminBatchPauseWorkflowsRequest{
			ChainSelector: chainSelector,
			WorkflowIDs:   workflowIDs,
		})
		require.NoError(t, err)

		// Then pause all workflows by owner
		owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
		_, err = AdminPauseAllByOwner(env, AdminPauseAllByOwnerRequest{
			ChainSelector: chainSelector,
			Owner:         owner,
		})
		require.NoError(t, err)

		// Finally, pause all workflows by DON
		_, err = AdminPauseAllByDON(env, AdminPauseAllByDONRequest{
			ChainSelector: chainSelector,
			DONFamily:     "emergency_don",
		})
		require.NoError(t, err)
	})
}
