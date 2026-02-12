package v1_6_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestDeployRegistryModuleChangeset(t *testing.T) {
	t.Parallel()

	t.Run("successfully deploys registry module on multiple chains", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector
		chain2 := chain_selectors.TEST_90000002.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1, chain2}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		chainSelectors := []uint64{chain1, chain2}

		// Deploy prerequisites (Link token and TokenAdminRegistry)
		prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
		for _, chain := range chainSelectors {
			prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
				ChainSelector: chain,
			})
		}

		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
				chainSelectors,
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
				changeset.DeployPrerequisiteConfig{
					Configs: prereqCfg,
				},
			),
		)
		require.NoError(t, err)

		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: chainSelectors,
		}

		// Use Apply to deploy and merge the results back into the environment
		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Load state - the environment now has the deployed contracts
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		// Verify we have registry modules deployed on both chains
		for _, chainSel := range chainSelectors {
			chainState, ok := state.Chains[chainSel]
			require.True(t, ok, "chain %d not found in state", chainSel)
			require.NotEmpty(t, chainState.RegistryModules1_6, "should have registry modules on chain %d", chainSel)
		}
	})

	t.Run("skips deployment if already deployed", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites
		prereqCfg := []changeset.DeployPrerequisiteConfigPerChain{
			{ChainSelector: chain1},
		}

		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
				[]uint64{chain1},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
				changeset.DeployPrerequisiteConfig{
					Configs: prereqCfg,
				},
			),
		)
		require.NoError(t, err)

		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{chain1},
		}

		// Deploy once using Apply
		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Verify deployment exists
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState, ok := state.Chains[chain1]
		require.True(t, ok, "chain %d not found in state", chain1)
		require.NotEmpty(t, chainState.RegistryModules1_6, "should have registry modules after first deployment")

		// Deploy again - should skip (Apply will handle this gracefully)
		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Verify still only one deployment
		state2, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)
		require.Equal(t, state.Chains[chain1].RegistryModules1_6, state2.Chains[chain1].RegistryModules1_6, "should not deploy twice")
	})

	t.Run("deploys on single chain", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites
		prereqCfg := []changeset.DeployPrerequisiteConfigPerChain{
			{ChainSelector: chain1},
		}

		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
				[]uint64{chain1},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
				changeset.DeployPrerequisiteConfig{
					Configs: prereqCfg,
				},
			),
		)
		require.NoError(t, err)

		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{chain1},
		}

		// Deploy using Apply
		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Verify only deployed on the specified chain
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState, ok := state.Chains[chain1]
		require.True(t, ok, "chain %d not found in state", chain1)
		require.NotEmpty(t, chainState.RegistryModules1_6, "should have registry modules on chain %d", chain1)
	})
}

// TestDeployRegistryModuleConfig_Validate tests the configuration validation logic
func TestDeployRegistryModuleConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("fails with empty chain selectors", func(t *testing.T) {
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{},
		}

		require.Empty(t, cfg.ChainSelectors)
	})

	t.Run("fails with nil chain selectors", func(t *testing.T) {
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: nil,
		}

		require.Nil(t, cfg.ChainSelectors)
	})

	t.Run("config structure is valid", func(t *testing.T) {
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{123456, 789012},
		}

		require.Len(t, cfg.ChainSelectors, 2)
		require.Equal(t, uint64(123456), cfg.ChainSelectors[0])
		require.Equal(t, uint64(789012), cfg.ChainSelectors[1])
	})
}

// TestNeedsRegistryModule16Deployment tests the NeedsRegistryModule16Deployment helper function
func TestNeedsRegistryModule16Deployment(t *testing.T) {
	t.Parallel()

	t.Run("returns true when registry modules is empty slice", func(t *testing.T) {
		chainState := evm.CCIPChainState{
			RegistryModules1_6: []*registry_module_owner_custom.RegistryModuleOwnerCustom{},
		}

		result := v1_6.NeedsRegistryModule16Deployment(chainState)
		require.True(t, result, "should need deployment when no registry modules exist")
	})

	t.Run("returns true when registry modules is nil", func(t *testing.T) {
		chainState := evm.CCIPChainState{
			RegistryModules1_6: nil,
		}

		result := v1_6.NeedsRegistryModule16Deployment(chainState)
		require.True(t, result, "should need deployment when registry modules is nil")
	})

	t.Run("returns false when at least one 1.6.0 registry module exists", func(t *testing.T) {
		chainState := evm.CCIPChainState{
			RegistryModules1_6: []*registry_module_owner_custom.RegistryModuleOwnerCustom{
				{}, // Empty struct is fine - we're just checking the length
			},
		}

		result := v1_6.NeedsRegistryModule16Deployment(chainState)
		require.False(t, result, "should not need deployment when 1.6.0 registry module exists")
	})

	t.Run("returns false when multiple 1.6.0 registry modules exist", func(t *testing.T) {
		chainState := evm.CCIPChainState{
			RegistryModules1_6: []*registry_module_owner_custom.RegistryModuleOwnerCustom{
				{},
				{},
			},
		}

		result := v1_6.NeedsRegistryModule16Deployment(chainState)
		require.False(t, result, "should not need deployment when multiple 1.6.0 registry modules exist")
	})
}
