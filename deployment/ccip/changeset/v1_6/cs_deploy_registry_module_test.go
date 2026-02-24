package v1_6_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestDeployRegistryModuleChangeset(t *testing.T) {
	t.Parallel()

	t.Run("successfully deploys registry module to single chain", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites first
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

		// Deploy registry module
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{chain1},
		}

		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Verify deployment
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState := state.Chains[chain1]
		require.NotEmpty(t, chainState.RegistryModules1_6, "should have deployed registry module")

		// Verify the registry module has correct owner (TokenAdminRegistry)
		// for _, module := range chainState.RegistryModules1_6 {
		// 	owner, err := module.Owner(nil)
		// 	require.NoError(t, err)
		// 	require.Equal(t, chainState.TokenAdminRegistry.Address(), owner, "registry module owner should be TokenAdminRegistry")
		// }
	})

	t.Run("successfully deploys registry module to multiple chains", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector
		chain2 := chain_selectors.TEST_90000002.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1, chain2}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		chainSelectors := []uint64{chain1, chain2}

		// Deploy prerequisites
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

		// Deploy registry modules
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: chainSelectors,
		}

		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Verify deployment on all chains
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		for _, chainSel := range chainSelectors {
			chainState := state.Chains[chainSel]
			require.NotEmpty(t, chainState.RegistryModules1_6, "should have deployed registry module on chain %d", chainSel)

			// Verify owner
			// for _, module := range chainState.RegistryModules1_6 {
			// 	owner, err := module.Owner(nil)
			// 	require.NoError(t, err)
			// 	require.Equal(t, chainState.TokenAdminRegistry.Address(), owner, "registry module owner should be TokenAdminRegistry on chain %d", chainSel)
			// }
		}
	})

	t.Run("skips deployment if registry module already exists", func(t *testing.T) {
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

		// Deploy first time
		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Get the deployed module address
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)
		var firstAddr string
		for _, module := range state.Chains[chain1].RegistryModules1_6 {
			firstAddr = module.Address().Hex()
			break
		}

		// Try deploying again - should skip
		_, err = v1_6.DeployRegistryModuleChangeset(*env, cfg)
		require.NoError(t, err)
		// When skipping deployment, just verify no error occurred
		// The changeset may return a DataStore even when skipping

		// Verify the same module still exists
		state, err = stateview.LoadOnchainState(*env)
		require.NoError(t, err)
		var secondAddr string
		for _, module := range state.Chains[chain1].RegistryModules1_6 {
			secondAddr = module.Address().Hex()
			break
		}
		require.Equal(t, firstAddr, secondAddr, "should be the same registry module")
	})

	t.Run("fails without prerequisites", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Try to deploy registry module without prerequisites
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{chain1},
		}

		_, err = v1_6.DeployRegistryModuleChangeset(*env, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "TokenAdminRegistry not found")
	})

	t.Run("deploys with MCMS", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites and MCMS
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
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
				map[uint64]commontypes.MCMSWithTimelockConfigV2{
					chain1: proposalutils.SingleGroupTimelockConfigV2(t),
				},
			),
		)
		require.NoError(t, err)

		// Deploy registry module
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{chain1},
		}

		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Verify deployment
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState := state.Chains[chain1]
		require.NotEmpty(t, chainState.RegistryModules1_6, "should have deployed registry module")
	})

	t.Run("deploys to multiple chains with MCMS", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector
		chain2 := chain_selectors.TEST_90000002.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1, chain2}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		chainSelectors := []uint64{chain1, chain2}

		// Deploy prerequisites and MCMS
		prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
		for _, chain := range chainSelectors {
			prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
				ChainSelector: chain,
			})
		}

		mcmsConfigs := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
		for _, chain := range chainSelectors {
			mcmsConfigs[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
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
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
				mcmsConfigs,
			),
		)
		require.NoError(t, err)

		// Deploy registry modules
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: chainSelectors,
		}

		*env, err = commonchangeset.Apply(t, *env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				cfg,
			),
		)
		require.NoError(t, err)

		// Verify deployment on all chains
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		for _, chainSel := range chainSelectors {
			chainState := state.Chains[chainSel]
			require.NotEmpty(t, chainState.RegistryModules1_6, "should have deployed registry module on chain %d", chainSel)
		}
	})
}

// TestDeployRegistryModuleConfig_Validate tests the configuration validation
func TestDeployRegistryModuleConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("fails with empty chain selectors", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no chain selectors provided")
	})

	t.Run("fails with nil chain selectors", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: nil,
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no chain selectors provided")
	})

	t.Run("fails with invalid chain selector", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{999999},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid chain selector")
	})

	t.Run("fails with chain not in environment", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector
		chain2 := chain_selectors.TEST_90000002.Selector

		// Only add chain1 to environment
		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Try to configure chain2 which doesn't exist in environment
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{chain2},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain state not found")
	})

	t.Run("fails when TokenAdminRegistry not deployed", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Don't deploy prerequisites
		cfg := v1_6.DeployRegistryModuleConfig{
			ChainSelectors: []uint64{chain1},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "TokenAdminRegistry not found")
	})

	t.Run("succeeds with valid config", func(t *testing.T) {
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

		err = cfg.Validate(*env)
		require.NoError(t, err)
	})

	t.Run("succeeds with multiple chains", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector
		chain2 := chain_selectors.TEST_90000002.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1, chain2}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		chainSelectors := []uint64{chain1, chain2}

		// Deploy prerequisites on both chains
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

		err = cfg.Validate(*env)
		require.NoError(t, err)
	})
}
