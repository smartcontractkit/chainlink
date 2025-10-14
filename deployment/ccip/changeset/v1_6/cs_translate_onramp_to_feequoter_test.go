package v1_6_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
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

	allChains := tenv.BlockChains.ListChainSelectors(
		cldf_chain.WithFamily(chain_selectors.FamilyEVM),
		cldf_chain.WithChainSelectorsExclusion([]uint64{chain_selectors.GETH_TESTNET.Selector}),
	)

	selectorA, selectorB := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: selectorA, DestChainSelector: selectorB},
		{SourceChainSelector: selectorB, DestChainSelector: selectorA},
	}

	// 3. Remove link token as it will be deployed by 1.6 contracts again
	ab := cldf.NewMemoryAddressBook()
	for _, sel := range allChains {
		require.NoError(t, ab.Save(sel, state.Chains[sel].LinkToken.Address().Hex(),
			cldf.NewTypeAndVersion("LinkToken", deployment.Version1_0_0)))
	}
	require.NoError(t, tenv.ExistingAddresses.Remove(ab))

	// 4. Set the test router as the source chain's router
	ab = cldf.NewMemoryAddressBook()
	for _, sel := range allChains {
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

	// 8. get onramp & feequoter dynamic, tokenCfg & default configs to compare

	feeTokenCfg, err := onRamp1_5Contract.GetFeeTokenConfig(&bind.CallOpts{Context: ctx}, allFeeTokens[0])
	require.NoError(t, err, "Failed to get FeeTokenConfig from OnRamp")
	actualFeeQuoterDestCfg, err := feeQuoterContract.GetDestChainConfig(&bind.CallOpts{Context: ctx}, destChainSelector)
	require.NoError(t, err, "Failed to get DestChainConfig from FeeQuoter after translation")

	defaultCfgForFamily := v1_6.DefaultFeeQuoterDestChainConfig(true, destChainSelector)

	// 9.Compare the actual configuration with the expected one
	// TODO: ensure all the fields are compared
	// Criteria 1: Ports dynamic config from all 1.5.0 OnRamps into FeeQuoter DestChainConfig
	require.Equal(t, onRampDynamicCfg.MaxNumberOfTokensPerMsg, actualFeeQuoterDestCfg.MaxNumberOfTokensPerMsg, "MaxNumberOfTokensPerMsg mismatch")
	require.Equal(t, onRampDynamicCfg.MaxDataBytes, actualFeeQuoterDestCfg.MaxDataBytes, "MaxDataBytes mismatch")
	require.Equal(t, onRampDynamicCfg.MaxPerMsgGasLimit, actualFeeQuoterDestCfg.MaxPerMsgGasLimit, "MaxPerMsgGasLimit mismatch")
	require.Equal(t, onRampDynamicCfg.DestGasOverhead, actualFeeQuoterDestCfg.DestGasOverhead, "DestGasOverhead mismatch")
	require.Equal(t, onRampDynamicCfg.DefaultTokenFeeUSDCents, actualFeeQuoterDestCfg.DefaultTokenFeeUSDCents, "DefaultTokenFeeUSDCents mismatch")
	require.Equal(t, onRampDynamicCfg.DestGasPerPayloadByte, uint16(actualFeeQuoterDestCfg.DestGasPerPayloadByteBase), "DestGasPerPayloadByteBase mismatch")
	require.Equal(t, onRampDynamicCfg.DestDataAvailabilityOverheadGas, actualFeeQuoterDestCfg.DestDataAvailabilityOverheadGas, "DestDataAvailabilityOverheadGas mismatch")
	require.Equal(t, onRampDynamicCfg.DestGasPerDataAvailabilityByte, actualFeeQuoterDestCfg.DestGasPerDataAvailabilityByte, "DestGasPerDataAvailabilityByte mismatch")
	require.Equal(t, onRampDynamicCfg.DestDataAvailabilityMultiplierBps, actualFeeQuoterDestCfg.DestDataAvailabilityMultiplierBps, "DestDataAvailabilityMultiplierBps mismatch")
	require.Equal(t, onRampDynamicCfg.DefaultTokenDestGasOverhead, actualFeeQuoterDestCfg.DefaultTokenDestGasOverhead, "DefaultTokenDestGasOverhead mismatch")
	require.Equal(t, defaultCfgForFamily.ChainFamilySelector, actualFeeQuoterDestCfg.ChainFamilySelector, "ChainFamilySelector mismatch")
	// These two should come from the GetFeeTokenConfig
	// Criteria 4 (b): Ports fee token config args from all 1.5.0 OnRamps into PremiumMultiplierWeiPerEthArgs
	require.Equal(t, newFeeQuoterParams.GasMultiplierWeiPerEth, actualFeeQuoterDestCfg.GasMultiplierWeiPerEth, "GasMultiplierWeiPerEth mismatch")
	require.Equal(t, newFeeQuoterParams.NetworkFeeUSDCents, actualFeeQuoterDestCfg.NetworkFeeUSDCents, "NetworkFeeUSDCents mismatch")

	// Criteria 3: Port supported fee tokens to the FeeQuoter if they do not yet exist on the FeeQuoter
	feetokensFromFeeQ, err = feeQuoterContract.GetFeeTokens(&bind.CallOpts{Context: ctx})
	require.NoError(t, err, "Failed to get GetFeeTokens from FeeQuoter")
	require.Len(t, feetokensFromFeeQ, 3, "Expected 3 fee tokens in FeeQuoter after translation") // same common token already exists between 1.5 & 1.6

	// Criteria 4: Ports fee token config args from all 1.5.0 OnRamps into PremiumMultiplierWeiPerEthArgs
	fqPremiumMultiplierCfg, err := feeQuoterContract.GetPremiumMultiplierWeiPerEth(&bind.CallOpts{Context: ctx}, allFeeTokens[0])
	require.NoError(t, err, "Failed to get PremiumMultiplierWeiPerEth from FeeQuoter")
	require.Equal(t, feeTokenCfg.PremiumMultiplierWeiPerEth, fqPremiumMultiplierCfg, "PremiumMultiplierWeiPerEth should match after translation")
	t.Logf("Successfully verified translation of 1.5.0 OnRamp config for chain %d to 1.6.0 FeeQuoter DestChainConfig for destination %d", sourceChainSelector, destChainSelector)

	// 10. E2E AddTokens & TokenPools
	tenv = DeployTokensAndTokenPools(t, tenv, &tenv.ExistingAddresses)

	tokenArContract, err := token_admin_registry.NewTokenAdminRegistry(sourceChainState.TokenAdminRegistry.Address(), tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err, "Failed to create TokenAdminRegistry contract binding")
	allTokens, _ := tokenArContract.GetAllConfiguredTokens(&bind.CallOpts{Context: ctx}, 0, 1000)
	require.Len(t, allTokens, 1, "Expected 1 token in TokenAdminRegistry after AddTokensE2E")
	tokenTransferFeeCfgFromOnRamp, err := onRamp1_5Contract.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, allTokens[0])
	require.NoError(t, err, "Failed to get TokenTransferFeeConfig from Onramp")
	tokenTransferFeeCfgFromFeeQ, err := feeQuoterContract.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, destChainSelector, allTokens[0])
	require.NoError(t, err, "Failed to get TokenTransferFeeConfig from FeeQuoter")

	// Should not match before translation
	require.NotEqual(t, tokenTransferFeeCfgFromOnRamp.DestBytesOverhead, tokenTransferFeeCfgFromFeeQ.DestBytesOverhead, "TokenTransferFeeConfig should not match before translation")

	// 11. Translate TokenTransferFeeConfig from OnRamp to FeeQuoter TokenTransferFeeConfig
	_, err = v1_6.TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset(tenv, translateConfig)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFeeQuoterChangeset execution failed")

	// Should match after translation
	// Criteria 2: Ports token transfer fee config args from all 1.5.0 OnRamps into FeeQuoter
	tokenTransferFeeCfgFromOnRamp, err = onRamp1_5Contract.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, allTokens[0])
	require.NoError(t, err, "Failed to get tokenTransferFeeCfgFromOnRamp from FeeQuoter")
	tokenTransferFeeCfgFromFeeQ, err = feeQuoterContract.GetTokenTransferFeeConfig(&bind.CallOpts{Context: ctx}, destChainSelector, allTokens[0])
	require.NoError(t, err, "Failed to get TokenTransferFeeConfig from FeeQuoter")
	require.Equal(t, tokenTransferFeeCfgFromOnRamp.DestBytesOverhead, tokenTransferFeeCfgFromFeeQ.DestBytesOverhead, "TokenTransferFeeConfig should match after translation(DestBytesOverhead)")
	require.Equal(t, tokenTransferFeeCfgFromOnRamp.MaxFeeUSDCents, tokenTransferFeeCfgFromFeeQ.MaxFeeUSDCents, "TokenTransferFeeConfig should match after translation (MaxFeeUSDCents)")
	require.Equal(t, tokenTransferFeeCfgFromOnRamp.DeciBps, tokenTransferFeeCfgFromFeeQ.DeciBps, "TokenTransferFeeConfig should match after translation (DeciBps)")
	require.Equal(t, tokenTransferFeeCfgFromOnRamp.DestGasOverhead, tokenTransferFeeCfgFromFeeQ.DestGasOverhead, "TokenTransferFeeConfig should match after translation (DestGasOverhead)")
	require.Equal(t, tokenTransferFeeCfgFromOnRamp.MinFeeUSDCents, tokenTransferFeeCfgFromFeeQ.MinFeeUSDCents, "TokenTransferFeeConfig should match after translation (MinFeeUSDCents)")
	require.Equal(t, tokenTransferFeeCfgFromOnRamp.IsEnabled, tokenTransferFeeCfgFromFeeQ.IsEnabled, "TokenTransferFeeConfig should match after translation (IsEnabled)")
	t.Logf("Successfully verified translation of 1.5.0 token transfer fee config args OnRamp config for chain %d to 1.6.0 FeeQuoter %d", sourceChainSelector, destChainSelector)
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
