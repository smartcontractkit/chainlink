package v1_6_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

func TestAddRegistryModuleChangeset(t *testing.T) {
	t.Parallel()

	t.Run("successfully adds registry module to single chain with MCMS", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites and registry module
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
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				v1_6.DeployRegistryModuleConfig{
					ChainSelectors: []uint64{chain1},
				},
			),
		)
		require.NoError(t, err)

		// Load state to get the deployed registry module address
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState := state.Chains[chain1]
		require.NotEmpty(t, chainState.RegistryModules1_6, "should have deployed registry module")

		// Get the registry module address
		var registryModuleAddr common.Address
		for _, module := range chainState.RegistryModules1_6 {
			registryModuleAddr = module.Address()
			break
		}

		// Create MCMS config for testing
		mcmsConfig := &proposalutils.TimelockConfig{
			MinDelay: 0,
		}

		// Add the registry module to TokenAdminRegistry with MCMS
		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				chain1: registryModuleAddr,
			},
			MCMSConfig: mcmsConfig,
		}

		*env, _, err = commonchangeset.ApplyChangesets(t, *env, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.AddRegistryModuleChangeset),
				cfg,
			),
		})
		require.NoError(t, err)

		// Verify the registry module was added
		state, err = stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState = state.Chains[chain1]
		isModule, err := chainState.TokenAdminRegistry.IsRegistryModule(nil, registryModuleAddr)
		require.NoError(t, err)
		require.True(t, isModule, "registry module should be added to TokenAdminRegistry")
	})

	t.Run("successfully adds registry module to multiple chains with MCMS", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector
		chain2 := chain_selectors.TEST_90000002.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1, chain2}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		chainSelectors := []uint64{chain1, chain2}

		// Deploy prerequisites and registry modules
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
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				v1_6.DeployRegistryModuleConfig{
					ChainSelectors: chainSelectors,
				},
			),
		)
		require.NoError(t, err)

		// Load state and collect registry module addresses
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		registryModuleAddrs := make(map[uint64]common.Address)
		for _, chainSel := range chainSelectors {
			chainState := state.Chains[chainSel]
			require.NotEmpty(t, chainState.RegistryModules1_6, "should have deployed registry module on chain %d", chainSel)

			for _, module := range chainState.RegistryModules1_6 {
				registryModuleAddrs[chainSel] = module.Address()
				break
			}
		}

		// Create MCMS config for testing
		mcmsConfig := &proposalutils.TimelockConfig{
			MinDelay: 0,
		}

		// Add registry modules to TokenAdminRegistry on all chains with MCMS
		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: registryModuleAddrs,
			MCMSConfig:          mcmsConfig,
		}

		*env, _, err = commonchangeset.ApplyChangesets(t, *env, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.AddRegistryModuleChangeset),
				cfg,
			),
		})
		require.NoError(t, err)

		// Verify registry modules were added on all chains
		state, err = stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		for _, chainSel := range chainSelectors {
			chainState := state.Chains[chainSel]
			addr := registryModuleAddrs[chainSel]
			isModule, err := chainState.TokenAdminRegistry.IsRegistryModule(nil, addr)
			require.NoError(t, err)
			require.True(t, isModule, "registry module should be added to TokenAdminRegistry on chain %d", chainSel)
		}
	})

	t.Run("skips if registry module already added", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites and registry module
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
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				v1_6.DeployRegistryModuleConfig{
					ChainSelectors: []uint64{chain1},
				},
			),
		)
		require.NoError(t, err)

		// Load state to get the registry module address
		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState := state.Chains[chain1]
		var registryModuleAddr common.Address
		for _, module := range chainState.RegistryModules1_6 {
			registryModuleAddr = module.Address()
			break
		}

		mcmsConfig := &proposalutils.TimelockConfig{
			MinDelay: 0,
		}

		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				chain1: registryModuleAddr,
			},
			MCMSConfig: mcmsConfig,
		}

		// Add registry module first time
		*env, _, err = commonchangeset.ApplyChangesets(t, *env, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.AddRegistryModuleChangeset),
				cfg,
			),
		})
		require.NoError(t, err)

		// Try adding again - should skip (no proposal generated)
		output2, err := v1_6.AddRegistryModuleChangeset(*env, cfg)
		require.NoError(t, err)
		require.Empty(t, output2.MCMSTimelockProposals, "should not generate proposal when already added")

		// Verify still registered
		state, err = stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		chainState = state.Chains[chain1]
		isModule, err := chainState.TokenAdminRegistry.IsRegistryModule(nil, registryModuleAddr)
		require.NoError(t, err)
		require.True(t, isModule, "registry module should still be registered")
	})

	t.Run("fails without MCMS config", func(t *testing.T) {
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
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
				map[uint64]commontypes.MCMSWithTimelockConfigV2{
					chain1: proposalutils.SingleGroupTimelockConfigV2(t),
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.DeployRegistryModuleChangeset),
				v1_6.DeployRegistryModuleConfig{
					ChainSelectors: []uint64{chain1},
				},
			),
		)
		require.NoError(t, err)

		state, err := stateview.LoadOnchainState(*env)
		require.NoError(t, err)

		var registryModuleAddr common.Address
		for _, module := range state.Chains[chain1].RegistryModules1_6 {
			registryModuleAddr = module.Address()
			break
		}

		// Try without MCMS config
		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				chain1: registryModuleAddr,
			},
			MCMSConfig: nil, // Should fail
		}

		_, err = v1_6.AddRegistryModuleChangeset(*env, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "mcmsConfig is required")
	})
}

// TestAddRegistryModuleConfig_Validate tests the configuration validation
func TestAddRegistryModuleConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("fails with empty registry module addresses", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{},
			MCMSConfig:          &proposalutils.TimelockConfig{MinDelay: 0},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no registry module addresses provided")
	})

	t.Run("fails with nil registry module addresses", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: nil,
			MCMSConfig:          &proposalutils.TimelockConfig{MinDelay: 0},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no registry module addresses provided")
	})

	t.Run("fails with zero address", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites to have TokenAdminRegistry
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

		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				chain1: common.Address{}, // Zero address
			},
			MCMSConfig: &proposalutils.TimelockConfig{MinDelay: 0},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "zero address")
	})

	t.Run("fails with invalid chain selector", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				999999: common.HexToAddress("0x1234567890123456789012345678901234567890"),
			},
			MCMSConfig: &proposalutils.TimelockConfig{MinDelay: 0},
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
		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				chain2: common.HexToAddress("0x1234567890123456789012345678901234567890"),
			},
			MCMSConfig: &proposalutils.TimelockConfig{MinDelay: 0},
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found in environment")
	})

	t.Run("fails without MCMS config", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites to have TokenAdminRegistry
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

		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				chain1: common.HexToAddress("0x1234567890123456789012345678901234567890"),
			},
			MCMSConfig: nil, // Should fail
		}

		err = cfg.Validate(*env)
		require.Error(t, err)
		require.Contains(t, err.Error(), "mcmsConfig is required")
	})

	t.Run("succeeds with valid config", func(t *testing.T) {
		chain1 := chain_selectors.TEST_90000001.Selector

		env, err := environment.New(t.Context(),
			environment.WithEVMSimulated(t, []uint64{chain1}),
			environment.WithLogger(logger.Test(t)),
		)
		require.NoError(t, err)

		// Deploy prerequisites to have TokenAdminRegistry
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

		cfg := v1_6.AddRegistryModuleConfig{
			RegistryModuleAddrs: map[uint64]common.Address{
				chain1: common.HexToAddress("0x1234567890123456789012345678901234567890"),
			},
			MCMSConfig: &proposalutils.TimelockConfig{MinDelay: 0},
		}

		err = cfg.Validate(*env)
		require.NoError(t, err)
	})
}
