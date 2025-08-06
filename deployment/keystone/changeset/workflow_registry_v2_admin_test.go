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
	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func setupWorkflowRegistryV2Test(t *testing.T) (cldf.Environment, uint64, cldf_evm.Chain, common.Address) {
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

func TestSetMetadataConfig(t *testing.T) {
	t.Parallel()

	env, chainSelector, chain, registryAddr := setupWorkflowRegistryV2Test(t)

	t.Run("successful metadata config", func(t *testing.T) {
		output, err := SetMetadataConfig(env, SetMetadataConfigRequest{
			ChainSelector: chainSelector,
			NameLen:       50,
			TagLen:        20,
			URLLen:        200,
			AttrLen:       1000,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)

		// Verify the configuration was set (if the contract has a getter)
		registry, err := workflow_registry_v2.NewWorkflowRegistry(registryAddr, chain.Client)
		require.NoError(t, err)
		require.NotNil(t, registry)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		_, err := SetMetadataConfig(env, SetMetadataConfigRequest{
			ChainSelector: 999999,
			NameLen:       50,
			TagLen:        20,
			URLLen:        200,
			AttrLen:       1000,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestSetWorkflowOwnerConfig(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupWorkflowRegistryV2Test(t)

	t.Run("successful workflow owner config", func(t *testing.T) {
		owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
		config := []byte(`{"maxWorkflows": 100}`)

		output, err := SetWorkflowOwnerConfig(env, SetWorkflowOwnerConfigRequest{
			ChainSelector: chainSelector,
			Owner:         owner,
			Config:        config,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
		config := []byte(`{"maxWorkflows": 100}`)

		_, err := SetWorkflowOwnerConfig(env, SetWorkflowOwnerConfigRequest{
			ChainSelector: 999999,
			Owner:         owner,
			Config:        config,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestSetDONLimit(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupWorkflowRegistryV2Test(t)

	t.Run("successful DON limit", func(t *testing.T) {
		output, err := SetDONLimit(env, SetDONLimitRequest{
			ChainSelector: chainSelector,
			DONFamily:     "test_don_family",
			Limit:         50,
			Enabled:       true,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		_, err := SetDONLimit(env, SetDONLimitRequest{
			ChainSelector: 999999,
			DONFamily:     "test_don_family",
			Limit:         50,
			Enabled:       true,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestSetUserDONOverride(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupWorkflowRegistryV2Test(t)

	t.Run("successful user DON override", func(t *testing.T) {
		user := common.HexToAddress("0x1234567890123456789012345678901234567890")

		output, err := SetUserDONOverride(env, SetUserDONOverrideRequest{
			ChainSelector: chainSelector,
			User:          user,
			DONFamily:     "test_don_family",
			Limit:         100,
			Enabled:       true,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		user := common.HexToAddress("0x1234567890123456789012345678901234567890")

		_, err := SetUserDONOverride(env, SetUserDONOverrideRequest{
			ChainSelector: 999999,
			User:          user,
			DONFamily:     "test_don_family",
			Limit:         100,
			Enabled:       true,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestSetDONRegistry(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupWorkflowRegistryV2Test(t)

	t.Run("successful DON registry", func(t *testing.T) {
		donRegistry := common.HexToAddress("0x1234567890123456789012345678901234567890")

		output, err := SetDONRegistry(env, SetDONRegistryRequest{
			ChainSelector:    chainSelector,
			Registry:         donRegistry,
			ChainSelectorDON: chainSelector,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		donRegistry := common.HexToAddress("0x1234567890123456789012345678901234567890")

		_, err := SetDONRegistry(env, SetDONRegistryRequest{
			ChainSelector:    999999,
			Registry:         donRegistry,
			ChainSelectorDON: chainSelector,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestUpdateAllowedSigners(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, _ := setupWorkflowRegistryV2Test(t)

	t.Run("successful update allowed signers", func(t *testing.T) {
		signers := []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			common.HexToAddress("0x2345678901234567890123456789012345678901"),
		}

		output, err := UpdateAllowedSigners(env, UpdateAllowedSignersRequest{
			ChainSelector: chainSelector,
			Signers:       signers,
			Allowed:       true,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("successful remove allowed signers", func(t *testing.T) {
		signers := []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
		}

		output, err := UpdateAllowedSigners(env, UpdateAllowedSignersRequest{
			ChainSelector: chainSelector,
			Signers:       signers,
			Allowed:       false,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		signers := []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
		}

		_, err := UpdateAllowedSigners(env, UpdateAllowedSignersRequest{
			ChainSelector: 999999,
			Signers:       signers,
			Allowed:       true,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})

	t.Run("empty signers list", func(t *testing.T) {
		output, err := UpdateAllowedSigners(env, UpdateAllowedSignersRequest{
			ChainSelector: chainSelector,
			Signers:       []common.Address{},
			Allowed:       true,
		})
		require.NoError(t, err)
		require.Equal(t, cldf.ChangesetOutput{}, output)
	})
}

func TestGetWorkflowRegistryV2(t *testing.T) {
	t.Parallel()

	env, chainSelector, _, registryAddr := setupWorkflowRegistryV2Test(t)

	t.Run("successful registry retrieval", func(t *testing.T) {
		registry, err := getWorkflowRegistryV2(env, chainSelector)
		require.NoError(t, err)
		require.NotNil(t, registry)

		// Verify it's the correct contract
		address := registry.Address()
		require.Equal(t, registryAddr, address)
	})

	t.Run("chain not found", func(t *testing.T) {
		_, err := getWorkflowRegistryV2(env, 999999)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get addresses for chain selector 999999")
	})

	t.Run("registry not found in address book", func(t *testing.T) {
		// Create new environment without workflow registry
		lggr, err := logger.New()
		require.NoError(t, err)
		emptyEnv := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
			Chains: 1,
			Nodes:  4,
		})

		var emptyChainSelector uint64
		for sel := range emptyEnv.BlockChains.EVMChains() {
			emptyChainSelector = sel
			break
		}

		_, err = getWorkflowRegistryV2(emptyEnv, emptyChainSelector)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no addresses found for chain selector")
	})
}
