package v1_6_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	// 1.5.0 contracts
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"

	// 1.6.0 contracts
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	// Workspace specific

	// For V1_5DeploymentConfig
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	v1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestTranslateEVM2EVMOnRampsToFeeQuoterChangeset(t *testing.T) {
	ctx := testcontext.Get(t)

	// 1. 1.5 deployment
	v1_5_deployment_config := &changeset.V1_5DeploymentConfig{
		PriceRegStalenessThreshold: 60 * 60 * 24, // 1 day
		RMNConfig: &rmn_contract.RMNConfig{ // Dummy RMN config
			BlessWeightThreshold: 1,
			CurseWeightThreshold: 1,
			Voters: []rmn_contract.RMNVoter{
				{BlessWeight: 1, CurseWeight: 1, BlessVoteAddr: utils.RandomAddress(), CurseVoteAddr: utils.RandomAddress()},
			},
		},
	}
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(v1_5_deployment_config), //price registry
	)

	// e, _ := testhelpers.NewMemoryEnvironment(t) // testrouter
	tenv := e.Env
	allChainSelectors := tenv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	require.Len(t, allChainSelectors, 2, "Expected 2 EVM chains")
	sourceChainSelector := allChainSelectors[0]
	destChainSelector := allChainSelectors[1]
	// 2. Load initial onchain state
	state, err := stateview.LoadOnchainState(tenv)
	require.NoError(t, err, "Failed to load initial onchain state")
	sourceChainState := state.MustGetEVMChainState(sourceChainSelector)

	allChains := tenv.BlockChains.ListChainSelectors(
		cldf_chain.WithFamily(chainselectors.FamilyEVM),
		cldf_chain.WithChainSelectorsExclusion([]uint64{chainselectors.GETH_TESTNET.Selector}),
	)

	selectorA, selectorB := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: selectorA, DestChainSelector: selectorB},
	}

	// 3. Deploy 1.5 Lanes
	tenv = v1_5.AddLanes(t, tenv, state, pairs)
	require.NoError(t, err)

	// 4. reload state after adding lanes
	state, err = stateview.LoadOnchainState(tenv)
	require.NoError(t, err)

	// 5. Remove link token as it will be deployed by 1.6 contracts again
	ab := cldf.NewMemoryAddressBook()
	for _, sel := range allChains {
		require.NoError(t, ab.Save(sel, state.Chains[sel].LinkToken.Address().Hex(),
			cldf.NewTypeAndVersion("LinkToken", deployment.Version1_0_0)))
	}
	//nolint:staticcheck //SA1019 ignoring deprecated
	require.NoError(t, tenv.ExistingAddresses.Remove(ab))

	// 6. Set the test router as the source chain's router
	ab = cldf.NewMemoryAddressBook()
	for _, sel := range allChains {
		require.NoError(t, ab.Save(sel, utils.RandomAddress().Hex(),
			cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0)))
	}
	//nolint:staticcheck //SA1019 ignoring deprecated
	require.NoError(t, tenv.ExistingAddresses.Merge(ab))

	// 7. Deploy 1.6.0 Pre-reqs contracts
	DeployUtil(t, &tenv, sourceChainSelector)

	// 8. Validate all needed contracts are deployed
	state, err = stateview.LoadOnchainState(tenv)
	require.NoError(t, err, "Failed to load initial onchain state")
	sourceChainState = state.MustGetEVMChainState(sourceChainSelector)
	require.NotNil(t, sourceChainState.EVM2EVMOnRamp, "1.5.0 OnRamps should be deployed on source chain")
	onRamp1_5_info, _ := sourceChainState.EVM2EVMOnRamp[destChainSelector]
	require.NotNil(t, onRamp1_5_info, "1.5.0 OnRamp instance info should not be nil")

	onRamp1_5_contract, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRamp1_5_info.Address(), tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err)
	feeQuoterContract, err := fee_quoter.NewFeeQuoter(sourceChainState.FeeQuoter.Address(), tenv.BlockChains.EVMChains()[sourceChainSelector].Client)
	require.NoError(t, err)

	// 9. Apply Translation Changeset
	translateConfig := v1_6.TranslateEVM2EVMOnRampsToFeeQuoterConfig{
		SourceChainSelectors: []uint64{sourceChainSelector},
		MCMS:                 nil, // Not testing MCMS interactions in this specific test
	}

	_, err = v1_6.TranslateEVM2EVMOnRampsToFeeQuoterChangeset(tenv, translateConfig)
	require.NoError(t, err, "TranslateEVM2EVMOnRampsToFeeQuoterChangeset execution failed")

	// 10. get onramp & feequoter dynamic & default configs to compare
	onRampDynamicCfg, err := onRamp1_5_contract.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	require.NoError(t, err, "Failed to get DestChainConfig from 1.5 onramp")
	actualFeeQuoterDestCfg, err := feeQuoterContract.GetDestChainConfig(&bind.CallOpts{Context: ctx}, destChainSelector)
	require.NoError(t, err, "Failed to get DestChainConfig from FeeQuoter after translation")

	defaultCfgForFamily := v1_6.DefaultFeeQuoterDestChainConfig(true, destChainSelector)

	// 11.Compare the actual configuration with the expected one
	require.Equal(t, onRampDynamicCfg.MaxNumberOfTokensPerMsg, actualFeeQuoterDestCfg.MaxNumberOfTokensPerMsg, "MaxNumberOfTokensPerMsg mismatch")
	require.Equal(t, onRampDynamicCfg.MaxDataBytes, actualFeeQuoterDestCfg.MaxDataBytes, "MaxDataBytes mismatch")
	require.Equal(t, onRampDynamicCfg.MaxPerMsgGasLimit, actualFeeQuoterDestCfg.MaxPerMsgGasLimit, "MaxPerMsgGasLimit mismatch")
	require.Equal(t, onRampDynamicCfg.DestGasOverhead, actualFeeQuoterDestCfg.DestGasOverhead, "DestGasOverhead mismatch")
	require.Equal(t, onRampDynamicCfg.DefaultTokenFeeUSDCents, actualFeeQuoterDestCfg.DefaultTokenFeeUSDCents, "DefaultTokenFeeUSDCents mismatch")
	require.Equal(t, onRampDynamicCfg.DestGasPerPayloadByte, actualFeeQuoterDestCfg.DestGasPerPayloadByteBase, "DestGasPerPayloadByteBase mismatch")
	require.Equal(t, onRampDynamicCfg.DestDataAvailabilityOverheadGas, actualFeeQuoterDestCfg.DestDataAvailabilityOverheadGas, "DestDataAvailabilityOverheadGas mismatch")
	require.Equal(t, onRampDynamicCfg.DestGasPerDataAvailabilityByte, actualFeeQuoterDestCfg.DestGasPerDataAvailabilityByte, "DestGasPerDataAvailabilityByte mismatch")
	require.Equal(t, onRampDynamicCfg.DestDataAvailabilityMultiplierBps, actualFeeQuoterDestCfg.DestDataAvailabilityMultiplierBps, "DestDataAvailabilityMultiplierBps mismatch")
	require.Equal(t, onRampDynamicCfg.DefaultTokenDestGasOverhead, actualFeeQuoterDestCfg.DefaultTokenDestGasOverhead, "DefaultTokenDestGasOverhead mismatch")
	require.Equal(t, defaultCfgForFamily.ChainFamilySelector, actualFeeQuoterDestCfg.ChainFamilySelector, "ChainFamilySelector mismatch")

	t.Logf("Successfully verified translation of 1.5.0 OnRamp config for chain %d to 1.6.0 FeeQuoter DestChainConfig for destination %d", sourceChainSelector, destChainSelector)
	/*
		externalAdmin := utils.ZeroAddress
		SelectorA2B := createSymmetricRateLimits(100, 1000)
		SelectorB2A := createSymmetricRateLimits(100, 1000)
		addTokenE2EConfig := v1_5_1.AddTokensE2EConfig{
			MCMS: nil,
		}
		for _, chain := range tenv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
			if addTokenE2EConfig.Tokens == nil {
				addTokenE2EConfig.Tokens = make(map[shared.TokenSymbol]v1_5_1.AddTokenE2EConfig)
			}
			if _, ok := addTokenE2EConfig.Tokens[testhelpers.TestTokenSymbol]; !ok {
				addTokenE2EConfig.Tokens[testhelpers.TestTokenSymbol] = v1_5_1.AddTokenE2EConfig{
					PoolConfig: make(map[uint64]v1_5_1.E2ETokenAndPoolConfig),
				}
			}
			rateLimiterPerChain := make(map[uint64]v1_5_1.RateLimiterConfig)
			for range []uint64{selectorA, selectorB} {
				switch chain {
				case selectorA:
					rateLimiterPerChain[selectorB] = SelectorA2B
				case selectorB:
					rateLimiterPerChain[selectorA] = SelectorB2A
				}
			}
			poolConfig := addTokenE2EConfig.Tokens[testhelpers.TestTokenSymbol].PoolConfig
			var deployPoolConfig *v1_5_1.DeployTokenPoolInput
			var deployTokenConfig *v1_5_1.DeployTokenConfig
			var _ *cldf.ContractType
			recipientAddress := utils.RandomAddress()
			topupAmount := big.NewInt(1000)
			deployTokenConfig = &v1_5_1.DeployTokenConfig{
				TokenName:     string(testhelpers.TestTokenSymbol),
				TokenSymbol:   testhelpers.TestTokenSymbol,
				TokenDecimals: testhelpers.LocalTokenDecimals,
				MaxSupply:     big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
				Type:          shared.BurnMintToken,
				PoolType:      shared.BurnMintTokenPool,
				MintTokenForRecipients: map[common.Address]*big.Int{
					recipientAddress: topupAmount,
				},
			}
			poolConfig[chain] = v1_5_1.E2ETokenAndPoolConfig{
				TokenDeploymentConfig: deployTokenConfig,
				DeployPoolConfig:      deployPoolConfig,
				PoolVersion:           deployment.Version1_5_1,
				ExternalAdmin:         externalAdmin,
			}
		} */

	// 2. E2E AddTokens
	/* 	tenv, err = commonchangeset.Apply(t, tenv,
		commonchangeset.Configure(v1_5_1.AddTokensE2E, addTokenE2EConfig))
	require.NoError(t, err) */
}

func createSymmetricRateLimits(rate int64, capacity int64) v1_5_1.RateLimiterConfig {
	return v1_5_1.RateLimiterConfig{
		Inbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
		Outbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
	}
}

func DeployUtil(t *testing.T, e *cldf.Environment, homeChainSel uint64) {
	evmSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
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
				"NodeOperator": p2pIds,
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
	state, err := stateview.LoadOnchainState(*e)
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
