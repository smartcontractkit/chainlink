package v1_6_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	v1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	migrate_seq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/migration"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

func TestTranslateEVM2EVMOnRampsToFeeQuoterChangeset(t *testing.T) {
	ctx := testcontext.Get(t)

	// 1. Deploy 1.5 pre-requisites
	v1_5DeploymentConfig := &changeset.V1_5DeploymentConfig{
		PriceRegStalenessThreshold: 60 * 60 * 24, // 1 day
		RMNConfig: &rmn_contract.RMNConfig{
			BlessWeightThreshold: 1,
			CurseWeightThreshold: 1,
			Voters: []rmn_contract.RMNVoter{
				{BlessWeight: 1, CurseWeight: 1, BlessVoteAddr: utils.RandomAddress(), CurseVoteAddr: utils.RandomAddress()},
			},
		},
	}

	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(v1_5DeploymentConfig), // price registry
	)

	tenv := e.Env

	allChainSelectors := tenv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChainSelectors, 2, "Expected 2 EVM chains")
	sourceChainSelector := allChainSelectors[0]
	destChainSelector := allChainSelectors[1]
	// 2. Load initial onchain state
	state, err := stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err, "Failed to load initial onchain state")

	selectorA, selectorB := allChainSelectors[0], allChainSelectors[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: selectorA, DestChainSelector: selectorB},
		{SourceChainSelector: selectorB, DestChainSelector: selectorA},
	}

	// 3. Remove link token as it will be deployed by 1.6 contracts again
	ab := cldf.NewMemoryAddressBook()
	for _, sel := range allChainSelectors {
		require.NoError(t, ab.Save(sel, state.Chains[sel].LinkToken.Address().Hex(),
			cldf.NewTypeAndVersion("LinkToken", deployment.Version1_0_0)))
	}
	require.NoError(t, tenv.ExistingAddresses.Remove(ab))

	// 4. Set the test router as the source chain's router
	ab = cldf.NewMemoryAddressBook()
	for _, sel := range allChainSelectors {
		require.NoError(t, ab.Save(sel, utils.RandomAddress().Hex(),
			cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0)))
	}
	require.NoError(t, tenv.ExistingAddresses.Merge(ab))

	// 4. Deploy 1.6.0 Pre-reqs contracts
	DeployUtil(t, &tenv, sourceChainSelector)
	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))

	// 5. Deploy 1.5 Lanes
	tenv = v1_5.AddLanes(t, tenv, state, pairs)
	require.NoError(t, err)

	// 6. Validate all needed contracts are deployed
	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err, "Failed to load initial onchain state")
	sourceChainState := state.MustGetEVMChainState(sourceChainSelector)
	require.NotNil(t, sourceChainState, "Src Chain state should not be nil")
	destChainState := state.MustGetEVMChainState(destChainSelector)
	require.NotNil(t, destChainState.EVM2EVMOnRamp, "1.5.0 OnRamps should be deployed on dest chain")
	onRamp1_5Info := sourceChainState.EVM2EVMOnRamp[destChainSelector]
	require.NotNil(t, onRamp1_5Info, "1.5.0 OnRamp instance info should not be nil")

	onRamp1_5Contract, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRamp1_5Info.Address(), tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err)
	feeQuoterContract, err := fee_quoter.NewFeeQuoter(sourceChainState.FeeQuoter.Address(), tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err)
	feetokensFromFeeQ, err := feeQuoterContract.GetFeeTokens(&bind.CallOpts{Context: ctx})
	require.NoError(t, err, "Failed to get GetFeeTokens from FeeQuoter")
	require.Len(t, feetokensFromFeeQ, 2, "Expected 2 fee token in FeeQuoter before translation changeset")
	onRampDynamicCfg, err := onRamp1_5Contract.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	require.NoError(t, err, "Failed to get DestChainConfig from 1.5 onramp")
	priceReg, err := price_registry.NewPriceRegistry(onRampDynamicCfg.PriceRegistry, tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err, "Failed to create PriceRegistry contract binding")
	allFeeTokens, err := priceReg.GetFeeTokens(nil)
	require.NoError(t, err, "Failed to get all fee tokens from PriceRegistry")
	require.Len(t, allFeeTokens, 2, "Expected 2 fee tokens in PriceRegistry before translation")

	// 7. Apply Translation Changeset
	newFeeQuoterParams := migrate_seq.NewFeeQuoterDestChainConfigParams{
		DestGasPerPayloadByteBase:      ccipevm.CalldataGasPerByteBase,
		DestGasPerPayloadByteHigh:      ccipevm.CalldataGasPerByteHigh,
		DestGasPerPayloadByteThreshold: ccipevm.CalldataGasPerByteThreshold,
		DefaultTxGasLimit:              200_000,
		ChainFamilySelector:            [4]byte{0x28, 0x12, 0xd5, 0x2c},
		GasPriceStalenessThreshold:     0,
		GasMultiplierWeiPerEth:         11e17,
		NetworkFeeUSDCents:             10,
	}
	newFeeQuoterParamsPerSource := make(map[uint64]migrate_seq.NewFeeQuoterDestChainConfigParams)
	for _, chain := range tenv.BlockChains.EVMChains() {
		if chain.Selector == destChainSelector {
			continue
		}
		newFeeQuoterParamsPerSource[chain.Selector] = newFeeQuoterParams
	}
	translateConfig := v1_6.TranslateEVM2EVMOnRampsToFeeQuoterConfig{
		NewFeeQuoterParamsPerSource: newFeeQuoterParamsPerSource,
		DestChainSelector:           destChainSelector,
		MCMS:                        nil, // Not testing MCMS interactions in this specific test
	}

	_, err = v1_6.TranslateEVM2EVMOnRampsToFeeQuoterChangeset(tenv, translateConfig)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFeeQuoterChangeset execution failed")

	// 8+9. Verify all FeeQuoter dest chain config fields (criteria 1), fee tokens (criteria 3),
	// and premium multipliers (criteria 4) on-chain.
	verifyFeeQuoterOnChainTranslation(ctx, t, onRamp1_5Contract, feeQuoterContract, newFeeQuoterParams, allFeeTokens, destChainSelector)
	t.Logf("Successfully verified translation of 1.5.0 OnRamp config for chain %d to 1.6.0 FeeQuoter DestChainConfig for destination %d", sourceChainSelector, destChainSelector)

	// 10. E2E AddTokens & TokenPools
	tenv = DeployTokensAndTokenPools(t, tenv, &tenv.ExistingAddresses)

	tokenArContract, err := token_admin_registry.NewTokenAdminRegistry(sourceChainState.TokenAdminRegistry.Address(), tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err, "Failed to create TokenAdminRegistry contract binding")
	allTokens, _ := tokenArContract.GetAllConfiguredTokens(&bind.CallOpts{Context: ctx}, 0, 1000)
	require.Len(t, allTokens, 1, "Expected 1 token in TokenAdminRegistry after AddTokensE2E")

	// Pre-condition: token transfer fee config does not match before translation.
	tokenTransferFeeCfgFromOnRamp, err := onRamp1_5Contract.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, allTokens[0])
	require.NoError(t, err)
	tokenTransferFeeCfgFromFeeQ, err := feeQuoterContract.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, destChainSelector, allTokens[0])
	require.NoError(t, err)
	require.NotEqual(t, tokenTransferFeeCfgFromOnRamp.DestBytesOverhead, tokenTransferFeeCfgFromFeeQ.DestBytesOverhead, "TokenTransferFeeConfig should not match before translation")

	// 11. Translate TokenTransferFeeConfig from OnRamp to FeeQuoter (criteria 2).
	_, err = v1_6.TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset(tenv, translateConfig)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset execution failed")

	verifyTokenTransferFeeConfigTranslation(ctx, t, onRamp1_5Contract, feeQuoterContract, allTokens, destChainSelector)
	t.Logf("Successfully verified translation of 1.5.0 token transfer fee config args OnRamp config for chain %d to 1.6.0 FeeQuoter %d", sourceChainSelector, destChainSelector)
}

// TestTranslateEVM2EVMOnRampsToFeeQuoterChangeset_WithMCMS verifies the full translation pipeline
// through MCMS: the FeeQuoter is transferred to the timelock, proposals are signed and executed
// via commonchangeset.Apply, and the same on-chain assertions as the direct-execution test are run.
func TestTranslateEVM2EVMOnRampsToFeeQuoterChangeset_WithMCMS(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)

	v1_5DeploymentConfig := &changeset.V1_5DeploymentConfig{
		PriceRegStalenessThreshold: 60 * 60 * 24,
		RMNConfig: &rmn_contract.RMNConfig{
			BlessWeightThreshold: 1,
			CurseWeightThreshold: 1,
			Voters: []rmn_contract.RMNVoter{
				{BlessWeight: 1, CurseWeight: 1, BlessVoteAddr: utils.RandomAddress(), CurseVoteAddr: utils.RandomAddress()},
			},
		},
	}
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(v1_5DeploymentConfig),
	)
	tenv := e.Env

	allChainSelectors := tenv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChainSelectors, 2)
	sourceChainSelector := allChainSelectors[0]
	destChainSelector := allChainSelectors[1]

	state, err := stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	ab := cldf.NewMemoryAddressBook()
	for _, sel := range allChainSelectors {
		require.NoError(t, ab.Save(sel, state.Chains[sel].LinkToken.Address().Hex(),
			cldf.NewTypeAndVersion("LinkToken", deployment.Version1_0_0)))
	}
	require.NoError(t, tenv.ExistingAddresses.Remove(ab))
	ab = cldf.NewMemoryAddressBook()
	for _, sel := range allChainSelectors {
		require.NoError(t, ab.Save(sel, utils.RandomAddress().Hex(),
			cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0)))
	}
	require.NoError(t, tenv.ExistingAddresses.Merge(ab))

	DeployUtil(t, &tenv, sourceChainSelector)
	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: sourceChainSelector, DestChainSelector: destChainSelector},
		{SourceChainSelector: destChainSelector, DestChainSelector: sourceChainSelector},
	}
	tenv = v1_5.AddLanes(t, tenv, state, pairs)
	tenv = DeployTokensAndTokenPools(t, tenv, &tenv.ExistingAddresses)

	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	onRamp1_5Info := state.MustGetEVMChainState(sourceChainSelector).EVM2EVMOnRamp[destChainSelector]
	require.NotNil(t, onRamp1_5Info)
	onRamp1_5Contract, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRamp1_5Info.Address(), tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err)

	onRampDynamicCfg, err := onRamp1_5Contract.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	priceReg, err := price_registry.NewPriceRegistry(onRampDynamicCfg.PriceRegistry, tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err)
	allFeeTokens, err := priceReg.GetFeeTokens(nil)
	require.NoError(t, err)

	tokenArContract, err := token_admin_registry.NewTokenAdminRegistry(
		state.MustGetEVMChainState(sourceChainSelector).TokenAdminRegistry.Address(),
		tenv.BlockChains.EVMChains()[sourceChainSelector].Client,
	)
	require.NoError(t, err)
	allTokens, _ := tokenArContract.GetAllConfiguredTokens(&bind.CallOpts{Context: ctx}, 0, 1000)
	require.Len(t, allTokens, 1)

	newFeeQuoterParams := migrate_seq.NewFeeQuoterDestChainConfigParams{
		DestGasPerPayloadByteBase:      ccipevm.CalldataGasPerByteBase,
		DestGasPerPayloadByteHigh:      ccipevm.CalldataGasPerByteHigh,
		DestGasPerPayloadByteThreshold: ccipevm.CalldataGasPerByteThreshold,
		DefaultTxGasLimit:              200_000,
		ChainFamilySelector:            [4]byte{0x28, 0x12, 0xd5, 0x2c},
		GasPriceStalenessThreshold:     0,
		GasMultiplierWeiPerEth:         11e17,
		NetworkFeeUSDCents:             10,
	}
	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay:   0 * time.Second,
		MCMSAction: types.TimelockActionSchedule,
	}
	translateConfig := v1_6.TranslateEVM2EVMOnRampsToFeeQuoterConfig{
		NewFeeQuoterParamsPerSource: map[uint64]migrate_seq.NewFeeQuoterDestChainConfigParams{
			sourceChainSelector: newFeeQuoterParams,
		},
		DestChainSelector: destChainSelector,
		MCMS:              mcmsConfig,
	}

	// Transfer the source chain's FeeQuoter ownership to the MCMS timelock so the
	// changeset proposal can be routed through it.
	tenv, err = commonchangeset.Apply(t, tenv,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					sourceChainSelector: {state.MustGetEVMChainState(sourceChainSelector).FeeQuoter.Address()},
				},
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: types.TimelockActionSchedule,
				},
			},
		),
	)
	require.NoError(t, err, "TransferToMCMSWithTimelockV2 failed")

	// Inspect proposals before executing
	csOut1, err := v1_6.TranslateEVM2EVMOnRampsToFeeQuoterChangeset(tenv, translateConfig)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFeeQuoterChangeset proposal generation failed")
	assertNoBatchForDestChain(t, csOut1, destChainSelector)

	csOut2, err := v1_6.TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset(tenv, translateConfig)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset proposal generation failed")
	assertNoBatchForDestChain(t, csOut2, destChainSelector)

	// Execute both changesets through MCMS — Apply signs and submits proposals via the timelock.
	tenv, err = commonchangeset.Apply(t, tenv,
		commonchangeset.Configure(v1_6.TranslateEVM2EVMOnRampsToFQDestConfig, translateConfig),
	)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFQDestConfig via MCMS failed")

	tenv, err = commonchangeset.Apply(t, tenv,
		commonchangeset.Configure(v1_6.TranslateEVM2EVMOnRampsToFQTokenTransferConfig, translateConfig),
	)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFQTokenTransferConfig via MCMS failed")

	// Reload state and bind FeeQuoter for on-chain assertions.
	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	feeQuoterContract, err := fee_quoter.NewFeeQuoter(
		state.MustGetEVMChainState(sourceChainSelector).FeeQuoter.Address(),
		tenv.BlockChains.EVMChains()[sourceChainSelector].Client,
	)
	require.NoError(t, err)

	// Same on-chain assertions as the direct-execution test — the only difference is that
	// the writes went through the MCMS timelock.
	verifyFeeQuoterOnChainTranslation(ctx, t, onRamp1_5Contract, feeQuoterContract, newFeeQuoterParams, allFeeTokens, destChainSelector)
	verifyTokenTransferFeeConfigTranslation(ctx, t, onRamp1_5Contract, feeQuoterContract, allTokens, destChainSelector)
	t.Logf("WithMCMS: verified on-chain translation for source %d → dest %d", sourceChainSelector, destChainSelector)
}

// verifyFeeQuoterOnChainTranslation asserts all FeeQuoter DestChainConfig fields, fee tokens, and
// premium multipliers
func verifyFeeQuoterOnChainTranslation(
	ctx context.Context,
	t *testing.T,
	onRamp1_5 *evm_2_evm_onramp.EVM2EVMOnRamp,
	feeQuoter *fee_quoter.FeeQuoter,
	newParams migrate_seq.NewFeeQuoterDestChainConfigParams,
	allFeeTokens []common.Address,
	destChainSelector uint64,
) {
	t.Helper()
	onRampDynCfg, err := onRamp1_5.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	actualFQDestCfg, err := feeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, destChainSelector)
	require.NoError(t, err)
	defaultCfg := v1_6.DefaultFeeQuoterDestChainConfig(true, destChainSelector)

	// Criteria 1: fields sourced from 1.5.0 OnRamp DynamicConfig
	require.True(t, actualFQDestCfg.IsEnabled, "IsEnabled must be hardcoded true")
	require.Equal(t, onRampDynCfg.MaxNumberOfTokensPerMsg, actualFQDestCfg.MaxNumberOfTokensPerMsg, "MaxNumberOfTokensPerMsg")
	require.Equal(t, onRampDynCfg.MaxDataBytes, actualFQDestCfg.MaxDataBytes, "MaxDataBytes")
	require.Equal(t, onRampDynCfg.MaxPerMsgGasLimit, actualFQDestCfg.MaxPerMsgGasLimit, "MaxPerMsgGasLimit")
	require.Equal(t, onRampDynCfg.DestGasOverhead, actualFQDestCfg.DestGasOverhead, "DestGasOverhead")
	require.Equal(t, onRampDynCfg.DefaultTokenFeeUSDCents, actualFQDestCfg.DefaultTokenFeeUSDCents, "DefaultTokenFeeUSDCents")
	require.Equal(t, onRampDynCfg.DestGasPerDataAvailabilityByte, actualFQDestCfg.DestGasPerDataAvailabilityByte, "DestGasPerDataAvailabilityByte")
	require.Equal(t, onRampDynCfg.DestDataAvailabilityOverheadGas, actualFQDestCfg.DestDataAvailabilityOverheadGas, "DestDataAvailabilityOverheadGas")
	require.Equal(t, onRampDynCfg.DestDataAvailabilityMultiplierBps, actualFQDestCfg.DestDataAvailabilityMultiplierBps, "DestDataAvailabilityMultiplierBps")
	require.Equal(t, onRampDynCfg.DefaultTokenDestGasOverhead, actualFQDestCfg.DefaultTokenDestGasOverhead, "DefaultTokenDestGasOverhead")
	require.Equal(t, onRampDynCfg.EnforceOutOfOrder, actualFQDestCfg.EnforceOutOfOrder, "EnforceOutOfOrder")
	// NewFeeQuoterDestChainConfigParams fields (no 1.5 equivalent, supplied explicitly)
	require.Equal(t, newParams.DestGasPerPayloadByteBase, actualFQDestCfg.DestGasPerPayloadByteBase, "DestGasPerPayloadByteBase")
	require.Equal(t, newParams.DestGasPerPayloadByteHigh, actualFQDestCfg.DestGasPerPayloadByteHigh, "DestGasPerPayloadByteHigh")
	require.Equal(t, newParams.DestGasPerPayloadByteThreshold, actualFQDestCfg.DestGasPerPayloadByteThreshold, "DestGasPerPayloadByteThreshold")
	require.Equal(t, newParams.DefaultTxGasLimit, actualFQDestCfg.DefaultTxGasLimit, "DefaultTxGasLimit")
	require.Equal(t, newParams.GasPriceStalenessThreshold, actualFQDestCfg.GasPriceStalenessThreshold, "GasPriceStalenessThreshold")
	require.Equal(t, defaultCfg.ChainFamilySelector, actualFQDestCfg.ChainFamilySelector, "ChainFamilySelector")
	require.Equal(t, newParams.GasMultiplierWeiPerEth, actualFQDestCfg.GasMultiplierWeiPerEth, "GasMultiplierWeiPerEth")
	require.Equal(t, newParams.NetworkFeeUSDCents, actualFQDestCfg.NetworkFeeUSDCents, "NetworkFeeUSDCents")

	// Criteria 3: fee tokens ported from PriceRegistry
	feeTokensFromFQ, err := feeQuoter.GetFeeTokens(&bind.CallOpts{Context: ctx})
	require.NoError(t, err)
	require.Len(t, feeTokensFromFQ, 3, "expected 3 fee tokens after translation (2 already in FeeQuoter + 1 migrated from PriceRegistry)")

	// Criteria 4: all fee tokens have their PremiumMultiplierWeiPerEth correctly translated.
	for _, ft := range allFeeTokens {
		onRampFTCfg, ftErr := onRamp1_5.GetFeeTokenConfig(&bind.CallOpts{Context: ctx}, ft)
		require.NoError(t, ftErr, "GetFeeTokenConfig for %s", ft)
		fqPremium, fqErr := feeQuoter.GetPremiumMultiplierWeiPerEth(&bind.CallOpts{Context: ctx}, ft)
		require.NoError(t, fqErr, "GetPremiumMultiplierWeiPerEth for %s", ft)
		require.Equal(t, onRampFTCfg.PremiumMultiplierWeiPerEth, fqPremium, "PremiumMultiplierWeiPerEth mismatch for token %s", ft)
	}
}

// verifyTokenTransferFeeConfigTranslation asserts all TokenTransferFeeConfig fields on-chain
// after TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset
func verifyTokenTransferFeeConfigTranslation(
	ctx context.Context,
	t *testing.T,
	onRamp1_5 *evm_2_evm_onramp.EVM2EVMOnRamp,
	feeQuoter *fee_quoter.FeeQuoter,
	allTokens []common.Address,
	destChainSelector uint64,
) {
	t.Helper()
	for _, token := range allTokens {
		fromOnRamp, err := onRamp1_5.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, token)
		require.NoError(t, err, "GetTokenTransferFeeConfig (OnRamp) for %s", token)
		fromFQ, err := feeQuoter.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, destChainSelector, token)
		require.NoError(t, err, "GetTokenTransferFeeConfig (FeeQuoter) for %s", token)
		require.Equal(t, fromOnRamp.MinFeeUSDCents, fromFQ.MinFeeUSDCents, "MinFeeUSDCents for %s", token)
		require.Equal(t, fromOnRamp.MaxFeeUSDCents, fromFQ.MaxFeeUSDCents, "MaxFeeUSDCents for %s", token)
		require.Equal(t, fromOnRamp.DeciBps, fromFQ.DeciBps, "DeciBps for %s", token)
		require.Equal(t, fromOnRamp.DestGasOverhead, fromFQ.DestGasOverhead, "DestGasOverhead for %s", token)
		require.Equal(t, fromOnRamp.DestBytesOverhead, fromFQ.DestBytesOverhead, "DestBytesOverhead for %s", token)
		require.Equal(t, fromOnRamp.IsEnabled, fromFQ.IsEnabled, "IsEnabled for %s", token)
	}
}

// assertNoBatchForDestChain is the regression check for the empty-batch bug.
// Before the fix, toSequenceInput unconditionally added an entry for every chain in the env
// including the dest chain, emitting an empty applyTokenTransferFeeConfigUpdates call for it.
func assertNoBatchForDestChain(t *testing.T, csOut cldf.ChangesetOutput, destChainSelector uint64) {
	t.Helper()
	for _, prop := range csOut.MCMSTimelockProposals {
		for _, op := range prop.Operations {
			require.NotEqual(t, uint64(op.ChainSelector), destChainSelector,
				"dest chain %d must not appear in proposal operations — empty batch regression", destChainSelector)
		}
	}
}

func DeployUtil(t *testing.T, e *cldf.Environment, homeChainSel uint64) {
	evmSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIDs := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]ccipseq.ChainContractParams)
	for _, chain := range e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		cfg[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
		contractParams[chain] = ccipseq.ChainContractParams{
			FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
			OffRampParams:   ccipops.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	eVal, err := commonchangeset.Apply(t, *e, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
		v1_6.DeployHomeChainConfig{
			HomeChainSel:     homeChainSel,
			RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
			RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
			NodeOperators:    testhelpers.NewTestNodeOperator(e.BlockChains.EVMChains()[homeChainSel].DeployerKey.From),
			NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
				"NodeOperator": p2pIDs,
			},
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
		evmSelectors,
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		cfg,
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
		changeset.DeployPrerequisiteConfig{
			Configs: prereqCfg,
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
		ccipseq.DeployChainContractsConfig{
			HomeChainSelector:      homeChainSel,
			ContractParamsPerChain: contractParams,
		},
	))
	require.NoError(t, err)
	*e = eVal // Update the environment pointed to by e

	// load onchain state
	state, err := stateview.LoadOnchainState(*e, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// verify all contracts populated
	require.NotNil(t, state.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, state.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, state.Chains[homeChainSel].RMNHome)
	for _, sel := range evmSelectors {
		require.NotNil(t, state.Chains[sel].LinkToken)
		require.NotNil(t, state.Chains[sel].Weth9)
		require.NotNil(t, state.Chains[sel].TokenAdminRegistry)
		require.NotNil(t, state.Chains[sel].RegistryModules1_6)
		require.NotNil(t, state.Chains[sel].Router)
		require.NotNil(t, state.Chains[sel].RMNRemote)
		require.NotNil(t, state.Chains[sel].TestRouter)
		require.NotNil(t, state.Chains[sel].NonceManager)
		require.NotNil(t, state.Chains[sel].FeeQuoter)
		require.NotNil(t, state.Chains[sel].OffRamp)
		require.NotNil(t, state.Chains[sel].OnRamp)
	}
}

const (
	LocalTokenDecimals                    = 18
	TestTokenSymbol    shared.TokenSymbol = "LINK"
)

func DeployTokensAndTokenPools(t *testing.T, e cldf.Environment, addressBook *cldf.AddressBook) cldf.Environment {
	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	selectorA, selectorB := selectors[0], selectors[1]
	state, err := stateview.LoadOnchainState(e, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	newPools := map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintTokenPool,
			TokenAddress:       state.Chains[selectorA].LinkToken.Address(),
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintTokenPool,
			TokenAddress:       state.Chains[selectorB].LinkToken.Address(),
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}

	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.DeployTokenPoolContractsChangeset),
			v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: TestTokenSymbol,
				NewPools:    newPools,
			},
		),
	)
	require.NoError(t, err)
	SelectorA2B := testhelpers.CreateSymmetricRateLimits(100, 1000)
	SelectorB2A := testhelpers.CreateSymmetricRateLimits(100, 1000)
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
			v1_5_1.ConfigureTokenPoolContractsConfig{
				TokenSymbol: TestTokenSymbol,
				MCMS:        nil,
				PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
					selectorA: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
						ChainUpdates: v1_5_1.RateLimiterPerChain{
							selectorB: SelectorA2B,
						},
					},
					selectorB: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
						ChainUpdates: v1_5_1.RateLimiterPerChain{
							selectorA: SelectorB2A,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	e, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			MCMS: nil,
			Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
				selectorA: {
					TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
				selectorB: {
					TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
			},
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			MCMS: nil,
			Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
				selectorA: {
					TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
				selectorB: {
					TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
			},
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.SetPoolChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			MCMS: nil,
			Pools: map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
				selectorA: {
					TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
				selectorB: {
					TestTokenSymbol: {
						Type:    shared.BurnMintTokenPool,
						Version: deployment.Version1_5_1,
					},
				},
			},
		},
	))
	require.NoError(t, err)
	return e
}
