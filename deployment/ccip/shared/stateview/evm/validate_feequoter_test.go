package evm_test

import (
	"math/big"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	fqv2ops "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/fee_quoter"
	fqv2seq "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/sequences"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	v1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

func TestValidateFeeQuoter_HappyPath(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	slices.Sort(evmChains)
	for _, sel := range evmChains {
		chainState := state.MustGetEVMChainState(sel)
		v16Active := buildV16ActiveChains(t, tenv, state)
		connectedChains, err := chainState.ValidateRouter(tenv.Env, false, v16Active)
		require.NoError(t, err, "router validation failed for chain %d", sel)

		err = chainState.ValidateFeeQuoter(tenv.Env, sel, connectedChains, nil, nil)
		require.NoError(t, err, "FeeQuoter validation failed for chain %d", sel)
	}
}

func TestValidateFeeQuoter_NilFeeQuoter(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chainState := state.MustGetEVMChainState(evmChains[0])
	chainState.FeeQuoter = nil
	err = chainState.ValidateFeeQuoter(tenv.Env, evmChains[0], evmChains[1:], nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no FeeQuoter")
}

// TestValidateFeeQuoter_CrossVersionValidation deploys v1.5, v1.6, and v2.0 FeeQuoter/OnRamp.
// "wrong_values" sets mismatched configs and asserts cross-version errors are reported.
// "fixed_values" aligns all configs and asserts validation passes.
func TestValidateFeeQuoter_CrossVersionValidation(t *testing.T) {
	t.Parallel()
	// 1. Deploy v1.5 prerequisites.
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
	source := allChainSelectors[0]
	dest := allChainSelectors[1]

	state, err := stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Capture old LinkToken addresses — v1.5 PriceRegistry still references these.
	oldLinkTokens := make(map[uint64]common.Address)
	for _, sel := range allChainSelectors {
		oldLinkTokens[sel] = state.Chains[sel].LinkToken.Address()
	}

	// 2. Remove LinkToken (re-deployed by v1.6).
	ab := cldf.NewMemoryAddressBook()
	for _, sel := range allChainSelectors {
		require.NoError(t, ab.Save(sel, state.Chains[sel].LinkToken.Address().Hex(),
			cldf.NewTypeAndVersion("LinkToken", deployment.Version1_0_0)))
	}
	require.NoError(t, tenv.ExistingAddresses.Remove(ab))

	// 3. Add TestRouter placeholder.
	ab = cldf.NewMemoryAddressBook()
	for _, sel := range allChainSelectors {
		require.NoError(t, ab.Save(sel, utils.RandomAddress().Hex(),
			cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0)))
	}
	require.NoError(t, tenv.ExistingAddresses.Merge(ab))

	// 4. Deploy v1.6 contracts.
	deployV16Contracts(t, &tenv, e.HomeChainSel)

	// 5. Add v1.5 lanes.
	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: source, DestChainSelector: dest},
		{SourceChainSelector: dest, DestChainSelector: source},
	}
	tenv = v1_5.AddLanes(t, tenv, state, pairs)

	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	sourceState := state.MustGetEVMChainState(source)
	require.NotNil(t, sourceState.EVM2EVMOnRamp[dest], "v1.5 OnRamp must exist for source→dest")

	evmChain := tenv.BlockChains.EVMChains()[source]
	linkAddr, err := sourceState.LinkTokenAddress()
	require.NoError(t, err)
	wethAddr := sourceState.Weth9.Address()

	// 6. Deploy v2.0 FeeQuoter and register fee tokens.
	fqV2Addr, fqV2 := deployV20FeeQuoter(t, evmChain, linkAddr)
	updateV20FeeQuoterFeeTokens(t, evmChain, fqV2, []common.Address{linkAddr, wethAddr})

	connectedChains := []uint64{dest}

	t.Run("wrong_values", func(t *testing.T) {
		// Values deliberately mismatched with v1.5 and v1.6 business rules.
		// Must still satisfy on-chain invariants (defaultTxGasLimit <= maxPerMsgGasLimit, etc.).
		badV16Cfg := fee_quoter.FeeQuoterDestChainConfig{
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           99,
			DestGasOverhead:                   999_999,
			DestDataAvailabilityOverheadGas:   1,
			DestGasPerDataAvailabilityByte:    1,
			DestDataAvailabilityMultiplierBps: 1,
			MaxDataBytes:                      50_000,
			MaxPerMsgGasLimit:                 5_000_000,
			DefaultTokenDestGasOverhead:       1_000,
			DefaultTokenFeeUSDCents:           99,
			EnforceOutOfOrder:                 true,
			DestGasPerPayloadByteBase:         99,
			DestGasPerPayloadByteHigh:         99,
			DestGasPerPayloadByteThreshold:    99,
			DefaultTxGasLimit:                 500_000,
			NetworkFeeUSDCents:                77,
			GasPriceStalenessThreshold:        0,
			GasMultiplierWeiPerEth:            99,
			ChainFamilySelector:               [4]byte{0x28, 0x12, 0xd5, 0x2c},
		}
		tenv, err = commonchangeset.Apply(t, tenv,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
				v1_6.UpdateFeeQuoterDestsConfig{
					UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
						source: {dest: badV16Cfg},
					},
				},
			),
		)
		require.NoError(t, err)

		// v2.0 config: different wrong values to trigger v1.6↔v2.0 cross-check failures.
		badV20Cfg := fqv2ops.DestChainConfig{
			IsEnabled:                   true,
			MaxDataBytes:                20_000,
			MaxPerMsgGasLimit:           400_000, // must be >= DefaultTxGasLimit
			DestGasOverhead:             111_111,
			DestGasPerPayloadByteBase:   50,
			ChainFamilySelector:         [4]byte{0x28, 0x12, 0xd5, 0x2c},
			DefaultTokenFeeUSDCents:     11,
			DefaultTokenDestGasOverhead: 500,
			DefaultTxGasLimit:           300_000,
			NetworkFeeUSDCents:          55,
			LinkFeeMultiplierPercent:    99,
		}
		updateV20FeeQuoterDestConfig(t, evmChain, fqV2Addr, dest, badV20Cfg)

		state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)
		chainState := state.MustGetEVMChainState(source)

		err = chainState.ValidateFeeQuoter(tenv, source, connectedChains, fqV2, evmChain.Client)
		require.Error(t, err, "validation must fail with wrong values")
		errMsg := err.Error()

		for _, field := range []string{
			"DestGasOverhead",
			"MaxDataBytes",
			"MaxPerMsgGasLimit",
			"MaxNumberOfTokensPerMsg",
			"DefaultTokenDestGasOverhead",
			"DefaultTokenFeeUSDCents",
			"EnforceOutOfOrder",
			"DestGasPerPayloadByteBase",
		} {
			assert.Contains(t, errMsg, field,
				"expected v1.5↔v1.6 cross-check to catch %s", field)
		}

		assert.Contains(t, errMsg, "NetworkFeeUSDCents", "v1.6 business rule")
		assert.Contains(t, errMsg, "DefaultTxGasLimit", "v1.6 business rule")
		assert.Contains(t, errMsg, "GasPriceStalenessThreshold", "v1.6 business rule")
		assert.Contains(t, errMsg, "v1.6<->v2.0", "v1.6↔v2.0 cross-check")
		assert.Contains(t, errMsg, "LinkFeeMultiplierPercent", "v2.0 business rule")
	})

	t.Run("fixed_values", func(t *testing.T) {
		callOpts := &bind.CallOpts{Context: t.Context()}
		state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)
		sourceState := state.MustGetEVMChainState(source)
		onRamp := sourceState.EVM2EVMOnRamp[dest]
		require.NotNil(t, onRamp)
		v15Cfg, err := onRamp.GetDynamicConfig(callOpts)
		require.NoError(t, err)

		// Align v1.6 FeeQuoter with v1.5 OnRamp dynamic config.
		goodV16Cfg := fee_quoter.FeeQuoterDestChainConfig{
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           v15Cfg.MaxNumberOfTokensPerMsg,
			DestGasOverhead:                   v15Cfg.DestGasOverhead,
			DestDataAvailabilityOverheadGas:   v15Cfg.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    v15Cfg.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: v15Cfg.DestDataAvailabilityMultiplierBps,
			MaxDataBytes:                      v15Cfg.MaxDataBytes,
			MaxPerMsgGasLimit:                 v15Cfg.MaxPerMsgGasLimit,
			DefaultTokenDestGasOverhead:       v15Cfg.DefaultTokenDestGasOverhead,
			DefaultTokenFeeUSDCents:           25, // topology: non-ETH EVM→EVM
			EnforceOutOfOrder:                 v15Cfg.EnforceOutOfOrder,
			DestGasPerPayloadByteBase:         uint8(v15Cfg.DestGasPerPayloadByte), //nolint:gosec // match v1.5 truncation
			DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
			DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
			DefaultTxGasLimit:                 200_000,
			NetworkFeeUSDCents:                10,
			GasPriceStalenessThreshold:        86400,
			ChainFamilySelector:               [4]byte{0x28, 0x12, 0xd5, 0x2c},
			GasMultiplierWeiPerEth:            11e17, // deploy-constant
		}
		tenv, err = commonchangeset.Apply(t, tenv,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
				v1_6.UpdateFeeQuoterDestsConfig{
					UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
						source: {dest: goodV16Cfg},
					},
				},
			),
		)
		require.NoError(t, err)

		// Align v2.0 FeeQuoter with v1.6.
		goodV20Cfg := fqv2ops.DestChainConfig{
			IsEnabled:                   true,
			MaxDataBytes:                goodV16Cfg.MaxDataBytes,
			MaxPerMsgGasLimit:           goodV16Cfg.MaxPerMsgGasLimit,
			DestGasOverhead:             goodV16Cfg.DestGasOverhead,
			DestGasPerPayloadByteBase:   goodV16Cfg.DestGasPerPayloadByteBase,
			ChainFamilySelector:         goodV16Cfg.ChainFamilySelector,
			DefaultTokenFeeUSDCents:     25, // topology: non-ETH EVM→EVM
			DefaultTokenDestGasOverhead: goodV16Cfg.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:           200_000,
			NetworkFeeUSDCents:          10,
			LinkFeeMultiplierPercent:    fqv2seq.LinkFeeMultiplierPercent,
		}
		updateV20FeeQuoterDestConfig(t, evmChain, fqV2Addr, dest, goodV20Cfg)

		// v1.5 default has NetworkFeeUSDCents=100; fix to 10.
		updateV15OnRampFeeTokenConfig(t, evmChain, onRamp, linkAddr, wethAddr)

		// Add old v1.5 LinkToken + WETH to v1.6 FeeQuoter for superset check.
		v16FQ := sourceState.FeeQuoter
		require.NotNil(t, v16FQ, "v1.6 FeeQuoter must exist")
		extraFeeTokens := []common.Address{wethAddr, oldLinkTokens[source]}
		fqTx, err := v16FQ.ApplyFeeTokensUpdates(evmChain.DeployerKey, nil, extraFeeTokens)
		require.NoError(t, err, "ApplyFeeTokensUpdates on v1.6 FeeQuoter")
		_, err = evmChain.Confirm(fqTx)
		require.NoError(t, err)

		updateV20FeeQuoterFeeTokens(t, evmChain, fqV2, []common.Address{oldLinkTokens[source]})

		state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)
		chainState := state.MustGetEVMChainState(source)

		err = chainState.ValidateFeeQuoter(tenv, source, connectedChains, fqV2, evmChain.Client)
		// v2.0 ownership can't be transferred in test (Timelock can't call acceptOwnership).
		if err != nil {
			assert.Contains(t, err.Error(), "not owned by Timelock",
				"only expected remaining error is v2.0 ownership")
			assert.NotContains(t, err.Error(), "mismatch",
				"no config mismatch errors should remain")
			assert.NotContains(t, err.Error(), "missing fee token",
				"no missing fee token errors should remain")
		}
	})
}

func deployV16Contracts(t *testing.T, tenv *cldf.Environment, homeChainSel uint64) {
	t.Helper()
	evmSelectors := tenv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	nodes, err := deployment.NodeInfo(tenv.NodeIDs, tenv.Offchain)
	require.NoError(t, err)
	p2pIDs := nodes.NonBootstraps().PeerIDs()

	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]ccipseq.ChainContractParams)
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, sel := range evmSelectors {
		cfg[sel] = proposalutils.SingleGroupTimelockConfigV2(t)
		contractParams[sel] = ccipseq.ChainContractParams{
			FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
			OffRampParams:   ccipops.DefaultOffRampParams(),
		}
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{ChainSelector: sel})
	}

	eVal, err := commonchangeset.Apply(t, *tenv, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
		v1_6.DeployHomeChainConfig{
			HomeChainSel:     homeChainSel,
			RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
			RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
			NodeOperators:    testhelpers.NewTestNodeOperator(tenv.BlockChains.EVMChains()[homeChainSel].DeployerKey.From),
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
		changeset.DeployPrerequisiteConfig{Configs: prereqCfg},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
		ccipseq.DeployChainContractsConfig{
			HomeChainSelector:      homeChainSel,
			ContractParamsPerChain: contractParams,
		},
	))
	require.NoError(t, err)
	*tenv = eVal
}

func deployV20FeeQuoter(t *testing.T, evmChain cldf_evm.Chain, linkToken common.Address) (common.Address, *fqv2ops.FeeQuoterContract) {
	t.Helper()
	parsedABI, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	require.NoError(t, err)

	fqV2Addr, tx, _, err := bind.DeployContract(
		evmChain.DeployerKey,
		parsedABI,
		common.FromHex(fqv2ops.FeeQuoterBin),
		evmChain.Client,
		fqv2ops.StaticConfig{
			MaxFeeJuelsPerMsg: big.NewInt(1e18),
			LinkToken:         linkToken,
		},
		[]common.Address{evmChain.DeployerKey.From},
		[]fqv2ops.TokenTransferFeeConfigArgs{},
		[]fqv2ops.DestChainConfigArgs{},
	)
	require.NoError(t, err)
	_, err = evmChain.Confirm(tx)
	require.NoError(t, err)

	fqV2, err := fqv2ops.NewFeeQuoterContract(fqV2Addr, evmChain.Client)
	require.NoError(t, err)
	return fqV2Addr, fqV2
}

// updateV20FeeQuoterFeeTokens registers fee tokens on v2.0 via updatePrices (auto-adds to s_feeTokens).
func updateV20FeeQuoterFeeTokens(t *testing.T, evmChain cldf_evm.Chain, fqV2 *fqv2ops.FeeQuoterContract, tokens []common.Address) {
	t.Helper()
	updates := make([]fqv2ops.TokenPriceUpdate, len(tokens))
	for i, tok := range tokens {
		updates[i] = fqv2ops.TokenPriceUpdate{
			SourceToken: tok,
			UsdPerToken: big.NewInt(1e18), // 1 USD — placeholder price
		}
	}
	tx, err := fqV2.UpdatePrices(evmChain.DeployerKey, fqv2ops.PriceUpdates{
		TokenPriceUpdates: updates,
	})
	require.NoError(t, err, "updatePrices for fee tokens")
	_, err = evmChain.Confirm(tx)
	require.NoError(t, err)
}

func updateV20FeeQuoterDestConfig(t *testing.T, evmChain cldf_evm.Chain, fqAddr common.Address, destSel uint64, cfg fqv2ops.DestChainConfig) {
	t.Helper()
	parsedABI, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	require.NoError(t, err)
	bc := bind.NewBoundContract(fqAddr, parsedABI, evmChain.Client, evmChain.Client, evmChain.Client)
	tx, err := bc.Transact(evmChain.DeployerKey, "applyDestChainConfigUpdates", []fqv2ops.DestChainConfigArgs{
		{DestChainSelector: destSel, DestChainConfig: cfg},
	})
	require.NoError(t, err, "applyDestChainConfigUpdates")
	_, err = evmChain.Confirm(tx)
	require.NoError(t, err)
}

func updateV15OnRampFeeTokenConfig(t *testing.T, evmChain cldf_evm.Chain, onRamp *evm_2_evm_onramp.EVM2EVMOnRamp, linkAddr, wethAddr common.Address) {
	t.Helper()
	tx, err := onRamp.SetFeeTokenConfig(evmChain.DeployerKey, []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
		{Token: linkAddr, NetworkFeeUSDCents: 10, GasMultiplierWeiPerEth: 1e18, PremiumMultiplierWeiPerEth: 9e17, Enabled: true},
		{Token: wethAddr, NetworkFeeUSDCents: 10, GasMultiplierWeiPerEth: 1e18, PremiumMultiplierWeiPerEth: 1e18, Enabled: true},
	})
	require.NoError(t, err, "SetFeeTokenConfig")
	_, err = evmChain.Confirm(tx)
	require.NoError(t, err)
}
