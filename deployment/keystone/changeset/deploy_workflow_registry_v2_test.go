package changeset

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployWorkflowRegistryV2(t *testing.T) {
	t.Parallel()

	lggr, err := logger.New()
	require.NoError(t, err)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
		Nodes:  4,
	})

	// Test deploying to all chains
	t.Run("deploy to all chains", func(t *testing.T) {
		output, err := DeployWorkflowRegistryV2(env, DeployWorkflowRegistryV2Request{
			ChainSelectors: nil, // Deploy to all chains
		})
		require.NoError(t, err)
		require.NotNil(t, output.AddressBook)
		require.NotNil(t, output.DataStore)

		// Verify deployment on each chain
		for chainSelector := range env.BlockChains.EVMChains() {
			addresses, err := output.DataStore.Addresses().Fetch()
			require.NoError(t, err)
			found := false
			for _, addr := range addresses {
				if addr.ChainSelector == chainSelector && addr.Type == "WorkflowRegistry" {
					require.NotEmpty(t, addr.Address)
					require.Equal(t, chainSelector, addr.ChainSelector)
					found = true
					break
				}
			}
			require.True(t, found, "WorkflowRegistry not found for chain %d", chainSelector)
		}
	})

	// Test deploying to specific chains
	t.Run("deploy to specific chains", func(t *testing.T) {
		chainSelectors := make([]uint64, 0, 1)
		for selector := range env.BlockChains.EVMChains() {
			chainSelectors = append(chainSelectors, selector)
			break // Only take the first chain
		}

		output, err := DeployWorkflowRegistryV2(env, DeployWorkflowRegistryV2Request{
			ChainSelectors: chainSelectors,
		})
		require.NoError(t, err)

		// Verify deployment only on specified chain
		addresses, err := output.DataStore.Addresses().Fetch()
		require.NoError(t, err)
		found := false
		for _, addr := range addresses {
			if addr.ChainSelector == chainSelectors[0] && addr.Type == "WorkflowRegistry" {
				require.NotEmpty(t, addr.Address)
				found = true
				break
			}
		}
		require.True(t, found, "WorkflowRegistry not found for chain %d", chainSelectors[0])
	})

	// Test invalid chain selector
	t.Run("invalid chain selector", func(t *testing.T) {
		_, err := DeployWorkflowRegistryV2(env, DeployWorkflowRegistryV2Request{
			ChainSelectors: []uint64{999999}, // Non-existent chain
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

func TestDeployWorkflowRegistryV2ToChain(t *testing.T) {
	t.Parallel()

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

	t.Run("successful deployment", func(t *testing.T) {
		ds := datastore.NewMemoryDataStore()
		resp, err := DeployWorkflowRegistryV2ToChain(context.Background(), chain, ds)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Verify response
		require.NotEqual(t, common.Address{}, resp.Address)
		require.NotEqual(t, common.Hash{}, resp.Tx)
		require.NotEmpty(t, resp.Tv.String())

		// Verify datastore entry
		addresses, err := ds.Addresses().Fetch()
		require.NoError(t, err)
		found := false
		for _, addr := range addresses {
			if addr.ChainSelector == chainSelector && addr.Address == resp.Address.String() {
				found = true
				break
			}
		}
		require.True(t, found, "Address not found in datastore")
	})
}

func TestDeployWorkflowRegistryV2Single(t *testing.T) {
	t.Parallel()

	lggr, err := logger.New()
	require.NoError(t, err)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
		Nodes:  4,
	})

	var chainSelector uint64
	for sel := range env.BlockChains.EVMChains() {
		chainSelector = sel
		break
	}

	t.Run("successful single chain deployment", func(t *testing.T) {
		output, err := DeployWorkflowRegistryV2Single(env, chainSelector)
		require.NoError(t, err)
		require.NotNil(t, output.AddressBook)
		require.NotNil(t, output.DataStore)

		// Verify deployment only on the specified chain
		addresses, err := output.DataStore.Addresses().Fetch()
		require.NoError(t, err)
		found := false
		for _, addr := range addresses {
			if addr.ChainSelector == chainSelector && addr.Type == "WorkflowRegistry" {
				require.NotEmpty(t, addr.Address)
				found = true
				break
			}
		}
		require.True(t, found, "WorkflowRegistry not found for chain %d", chainSelector)
	})

	t.Run("invalid chain selector", func(t *testing.T) {
		_, err := DeployWorkflowRegistryV2Single(env, 999999)
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain with selector 999999 not found")
	})
}

// TestWorkflowRegistryV2Integration tests the deployed contract functionality
func TestWorkflowRegistryV2Integration(t *testing.T) {
	t.Parallel()

	lggr, err := logger.New()
	require.NoError(t, err)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  4,
	})

	var chain cldf_evm.Chain
	for _, ch := range env.BlockChains.EVMChains() {
		chain = ch
		break
	}

	// Deploy the workflow registry
	ds := datastore.NewMemoryDataStore()
	resp, err := DeployWorkflowRegistryV2ToChain(context.Background(), chain, ds)
	require.NoError(t, err)

	// Test calling a simple function on the deployed contract
	t.Run("verify contract deployment and basic functionality", func(t *testing.T) {
		// Create contract instance
		registry, err := workflow_registry_v2.NewWorkflowRegistry(resp.Address, chain.Client)
		require.NoError(t, err)

		// Test calling TypeAndVersion
		tv, err := registry.TypeAndVersion(&bind.CallOpts{})
		require.NoError(t, err)
		require.Contains(t, tv, "WorkflowRegistry")

		// Test calling Owner
		owner, err := registry.Owner(&bind.CallOpts{})
		require.NoError(t, err)
		require.NotEqual(t, common.Address{}, owner)
	})
}
