package aptos_test

import (
	"math/big"
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	module_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	aptosconfig "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/operation/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestDynamicCS_Apply(t *testing.T) {
	t.Parallel()

	// Setup environment with Aptos chain and deployed contracts
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
	)
	env := deployedEnvironment.Env

	// Load onchain state to get deployed contract addresses
	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err, "must load onchain state")

	aptosChainSel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))[0]
	aptosState := state.AptosChains[aptosChainSel]
	aptosChain := env.BlockChains.AptosChains()[aptosChainSel]

	// Get a token address for price update (use any deployed token or a mock address)
	// For this test, we'll use the CCIP address as a mock token address
	mockTokenAddr := aptosState.CCIPAddress.StringLong()

	registry := operations.NewOperationRegistry(operation.GetAptosOperations()...)
	env.OperationsBundle.OperationRegistry = registry

	// Define the operations to execute
	defs := []operations.Definition{
		operation.ApplyAllowedOfframpUpdatesOp.Def(),
		operation.UpdateFeeQuoterDestsOp.Def(),
		operation.UpdateFeeQuoterPricesOp.Def(),
		aptos.CurseMultipleOp.Def(),
	}

	arbSubject := globals.FamilyAwareSelectorToSubject(
		chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
		chain_selectors.FamilyEVM,
	)
	bscSubject := globals.FamilyAwareSelectorToSubject(
		chain_selectors.BINANCE_SMART_CHAIN_MAINNET.Selector,
		chain_selectors.FamilyEVM,
	)

	// Define the inputs for each operation
	inputs := []any{
		// Input for ApplyAllowedOfframpUpdatesOp
		operations.EmptyInput{},
		// Input for UpdateFeeQuoterDestsOp
		operation.UpdateFeeQuoterDestsInput{
			Updates: map[uint64]module_fee_quoter.DestChainConfig{
				chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID: aptosTestDestFeeQuoterConfig(t),
			},
		},
		// Input for UpdateFeeQuoterPricesOp
		operation.UpdateFeeQuoterPricesInput{
			TokenPrices: map[string]*big.Int{
				mockTokenAddr: big.NewInt(1000001),
			},
			GasPrices: map[uint64]*big.Int{
				chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.EvmChainID: big.NewInt(500000), // Mock gas price
			},
		},
		// Input for CurseSubjectsOp
		aptos.CurseMultipleInput{
			CCIPAddress: aptosState.CCIPAddress,
			Subjects: [][]byte{
				arbSubject[:],
				bscSubject[:],
			},
		},
	}

	// Configure the dynamic changeset
	cfg := aptosconfig.DynamicConfig{
		Defs:          defs,
		Inputs:        inputs,
		ChainSelector: aptosChainSel,
		Description:   "Test dynamic changeset with multiple operations",
		MCMSConfig: &proposalutils.TimelockConfig{
			MinDelay:     time.Duration(1) * time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
	}

	// Apply the dynamic changeset
	env, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.DynamicCS{}, cfg),
	})
	require.NoError(t, err, "dynamic changeset should apply successfully")
	// Re-register operations after ApplyChangesets (bundle may be rebuilt)
	env.OperationsBundle.OperationRegistry = operations.NewOperationRegistry(operation.GetAptosOperations()...)

	// Verify the operations were executed successfully by checking the state
	// 1. Verify FeeQuoter prices were updated
	ccipBind := ccip.Bind(aptosState.CCIPAddress, aptosChain.Client)

	// Check token price
	tokenPrice, err := ccipBind.FeeQuoter().GetTokenPrice(nil, aptosState.CCIPAddress)
	require.NoError(t, err)
	require.NotNil(t, tokenPrice)
	require.Equal(t, big.NewInt(1000001), tokenPrice.Value, "token price should be updated")

	// 2. Verify allowed offramp updates were applied
	// The ApplyAllowedOfframpUpdatesOp adds the CCIP owner to the allowlist
	// Bind MCMS to get the owner address
	mcmsBind := mcmsbind.Bind(aptosState.MCMSAddress, aptosChain.Client)

	// Get CCIP owner address
	ccipOwnerAddress, err := mcmsBind.MCMSRegistry().GetRegisteredOwnerAddress(nil, aptosState.CCIPAddress)
	require.NoError(t, err)

	// Get the list of allowed offramps
	allowedOfframps, err := ccipBind.Auth().GetAllowedOfframps(nil)
	require.NoError(t, err)

	// Verify CCIP owner is in the allowlist
	found := false
	for _, addr := range allowedOfframps {
		if addr == ccipOwnerAddress {
			found = true
			break
		}
	}
	require.True(t, found, "CCIP owner should be in the allowlist after ApplyAllowedOfframpUpdatesOp")

	// 3. Verify subjects were cursed
	arbU128Selector := new(big.Int).SetUint64(chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector)
	bscU128Selector := new(big.Int).SetUint64(chain_selectors.BINANCE_SMART_CHAIN_MAINNET.Selector)
	isCursedU128, err := ccipBind.RMNRemote().IsCursedU128(nil, arbU128Selector)
	require.NoError(t, err)
	require.True(t, isCursedU128, "should be cursed")

	isCursed, err := ccipBind.RMNRemote().IsCursed(nil, bscSubject[:])
	require.NoError(t, err)
	require.True(t, isCursed, "should be cursed")

	// define the operations to execute
	defs = []operations.Definition{
		aptos.UncurseMultipleOp.Def(),
	}

	inputs = []any{
		aptos.UncurseMultipleInput{
			CCIPAddress: aptosState.CCIPAddress,
			Subjects: [][]byte{
				arbSubject[:],
			},
		},
	}

	cfg = aptosconfig.DynamicConfig{
		Defs:          defs,
		Inputs:        inputs,
		ChainSelector: aptosChainSel,
		Description:   "Test dynamic changeset with uncurse subjects operation",
		MCMSConfig: &proposalutils.TimelockConfig{
			MinDelay:     time.Duration(1) * time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
	}

	env, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.DynamicCS{}, cfg),
	})
	require.NoError(t, err, "dynamic changeset should apply successfully")
	// Re-register operations after ApplyChangesets (bundle may be rebuilt)
	env.OperationsBundle.OperationRegistry = operations.NewOperationRegistry(operation.GetAptosOperations()...)

	// Verify the operations were executed successfully by checking the state
	isCursedU128, err = ccipBind.RMNRemote().IsCursedU128(nil, arbU128Selector)
	require.NoError(t, err)
	require.False(t, isCursedU128, "should not be cursed")

	isCursed, err = ccipBind.RMNRemote().IsCursed(nil, arbSubject[:])
	require.NoError(t, err)
	require.False(t, isCursed, "should not be cursed")

	isCursedU128, err = ccipBind.RMNRemote().IsCursedU128(nil, bscU128Selector)
	require.NoError(t, err)
	require.True(t, isCursedU128, "should be cursed")

	isCursed, err = ccipBind.RMNRemote().IsCursed(nil, bscSubject[:])
	require.NoError(t, err)
	require.True(t, isCursed, "should be cursed")

	// define the operations to execute
	defs = []operations.Definition{
		aptos.CurseMultipleOp.Def(),
	}

	globalSubject := globals.GlobalCurseSubject()

	inputs = []any{
		aptos.CurseMultipleInput{
			CCIPAddress: aptosState.CCIPAddress,
			Subjects:    [][]byte{globalSubject[:]},
		},
	}

	cfg = aptosconfig.DynamicConfig{
		Defs:          defs,
		Inputs:        inputs,
		ChainSelector: aptosChainSel,
		Description:   "Test dynamic changeset with global curse operation",
		MCMSConfig: &proposalutils.TimelockConfig{
			MinDelay:     time.Duration(1) * time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
	}

	env, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.DynamicCS{}, cfg),
	})
	require.NoError(t, err, "dynamic changeset should apply successfully")

	// Verify the operations were executed successfully by checking the state
	isCursedGlobal, err := ccipBind.RMNRemote().IsCursedGlobal(nil)
	require.NoError(t, err)
	require.True(t, isCursedGlobal, "should be cursed globally")

	optimismSubject := globals.FamilyAwareSelectorToSubject(
		chain_selectors.ETHEREUM_MAINNET_OPTIMISM_1.Selector,
		chain_selectors.FamilyEVM,
	)

	isCursed, err = ccipBind.RMNRemote().IsCursed(nil, optimismSubject[:])
	require.NoError(t, err)
	require.True(t, isCursed, "should be cursed")
}
