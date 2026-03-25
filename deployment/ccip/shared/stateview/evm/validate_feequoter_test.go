package evm_test

import (
	"math/big"
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

// TestValidateFeeQuoter_HappyPath verifies that ValidateFeeQuoter passes on a
// correctly-deployed memory environment with all lanes configured (v1.6 only).
func TestValidateFeeQuoter_HappyPath(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	for _, sel := range evmChains {
		chainState := state.MustGetEVMChainState(sel)
		v16Active := buildV16ActiveChains(t, tenv, state)
		connectedChains, err := chainState.ValidateRouter(tenv.Env, false, v16Active)
		require.NoError(t, err, "router validation failed for chain %d", sel)

		err = chainState.ValidateFeeQuoter(tenv.Env, sel, connectedChains, nil, nil)
		require.NoError(t, err, "FeeQuoter validation failed for chain %d", sel)
	}
}

// TestValidateFeeQuoter_NilFeeQuoter verifies that ValidateFeeQuoter returns an
// error when no FeeQuoter contract is present.
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

// TestValidateFeeQuoter_CrossVersionValidation deploys v1.5 OnRamp + PriceRegistry,
// v1.6.3 FeeQuoter, and v2.0 FeeQuoter, then:
//   - Subtest "wrong_values": sets deliberately wrong dest chain configs on both FeeQuoters,
//     validates, and asserts that cross-version mismatches and business-rule violations are reported.
//   - Subtest "fixed_values": corrects all configs to align v1.5↔v1.6↔v2.0 with correct
//     business rules, validates, and asserts no errors.
func TestValidateFeeQuoter_CrossVersionValidation(t *testing.T) {
	t.Parallel()

	// ===== SETUP: v1.5 prereqs + v1.6 contracts + v1.5 lanes + v2.0 FeeQuoter =====

	// 1. Deploy with v1.5 prerequisites (PriceRegistry, RMN, etc.)
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

	// 2. Remove LinkToken (will be re-deployed by v1.6).
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

	// 4. Deploy v1.6 contracts (HomeChain, LinkToken, MCMS, Prerequisites, ChainContracts).
	deployV16Contracts(t, &tenv, e.HomeChainSel)

	// 5. Add v1.5 lanes (OnRamp + CommitStore + OffRamp per pair).
	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: source, DestChainSelector: dest},
		{SourceChainSelector: dest, DestChainSelector: source},
	}
	tenv = v1_5.AddLanes(t, tenv, state, pairs)

	// Reload state after v1.5 lane deployment.
	state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Verify v1.5 OnRamp is present on source chain.
	sourceState := state.MustGetEVMChainState(source)
	require.NotNil(t, sourceState.EVM2EVMOnRamp[dest], "v1.5 OnRamp must exist for source→dest")

	evmChain := tenv.BlockChains.EVMChains()[source]
	linkAddr, err := sourceState.LinkTokenAddress()
	require.NoError(t, err)
	wethAddr := sourceState.Weth9.Address()

	// 6. Deploy v2.0 FeeQuoter on source chain.
	fqV2Addr, fqV2 := deployV20FeeQuoter(t, evmChain, linkAddr)

	// 7. Add fee tokens (LINK + WETH) to v2.0 FeeQuoter.
	updateV20FeeQuoterFeeTokens(t, evmChain, fqV2Addr, []common.Address{linkAddr, wethAddr}, nil)

	connectedChains := []uint64{dest}

	// ===== SUBTEST 1: wrong values =====
	t.Run("wrong_values", func(t *testing.T) {
		// Set wrong v1.6 FeeQuoter dest config:
		// - Mismatches with v1.5 OnRamp (DestGasOverhead, MaxDataBytes, MaxPerMsgGasLimit, etc.)
		// - Business-rule violations (NetworkFeeUSDCents, DefaultTxGasLimit, ChainFamilySelector)
		badV16Cfg := fee_quoter.FeeQuoterDestChainConfig{
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           99,        // v1.5 has 5
			DestGasOverhead:                   999_999,   // v1.5 has 350_000
			DestDataAvailabilityOverheadGas:   1,         // v1.5 has 33_596
			DestGasPerDataAvailabilityByte:    1,         // v1.5 has 16
			DestDataAvailabilityMultiplierBps: 1,         // v1.5 has 6840
			MaxDataBytes:                      50_000,    // v1.5 has 100_000
			MaxPerMsgGasLimit:                 1_000,     // v1.5 has 4_000_000
			DefaultTokenDestGasOverhead:       1_000,     // v1.5 has 125_000
			DefaultTokenFeeUSDCents:           99,        // v1.5 has 50
			EnforceOutOfOrder:                 true,      // v1.5 has false
			DestGasPerPayloadByteBase:         99,        // v1.5 has 16
			DestGasPerPayloadByteHigh:         99,        // expected: CalldataGasPerByteHigh (40)
			DestGasPerPayloadByteThreshold:    99,        // expected: CalldataGasPerByteThreshold (3000)
			DefaultTxGasLimit:                 500_000,   // expected: 200_000
			NetworkFeeUSDCents:                77,        // expected: 10
			GasPriceStalenessThreshold:        0,         // must be non-zero
			GasMultiplierWeiPerEth:            99,        // v1.5 has 1e18
			ChainFamilySelector:               [4]byte{}, // must not be empty
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

		// Set wrong v2.0 FeeQuoter dest config:
		// - Different wrong values than v1.6 (triggers v1.6↔v2.0 cross-check failures)
		// - Business-rule violations (NetworkFeeUSDCents, DefaultTxGasLimit, LinkFeeMultiplierPercent)
		badV20Cfg := fqv2ops.DestChainConfig{
			IsEnabled:                   true,
			MaxDataBytes:                20_000,  // differs from v1.6 and v1.5
			MaxPerMsgGasLimit:           500,     // differs from v1.6 and v1.5
			DestGasOverhead:             111_111, // differs from v1.6 and v1.5
			DestGasPerPayloadByteBase:   50,      // differs from v1.6 and v1.5
			ChainFamilySelector:         [4]byte{0x28, 0x12, 0xd5, 0x2c},
			DefaultTokenFeeUSDCents:     11,      // differs from v1.6 and v1.5
			DefaultTokenDestGasOverhead: 500,     // differs from v1.6 and v1.5
			DefaultTxGasLimit:           300_000, // expected: 200_000
			NetworkFeeUSDCents:          55,      // expected: 10
			LinkFeeMultiplierPercent:    99,      // expected: fqv2seq.LinkFeeMultiplierPercent
		}
		updateV20FeeQuoterDestConfig(t, evmChain, fqV2Addr, dest, badV20Cfg)

		// Reload state and validate.
		state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)
		chainState := state.MustGetEVMChainState(source)

		err = chainState.ValidateFeeQuoter(tenv, source, connectedChains, fqV2, evmChain.Client)
		require.Error(t, err, "validation must fail with wrong values")
		errMsg := err.Error()

		// v1.5↔v1.6 cross-check fields
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
			assert.True(t, strings.Contains(errMsg, field),
				"expected v1.5↔v1.6 cross-check to catch %s, got: %s", field, errMsg)
		}

		// v1.6 business rules
		assert.Contains(t, errMsg, "NetworkFeeUSDCents", "v1.6 business rule")
		assert.Contains(t, errMsg, "DefaultTxGasLimit", "v1.6 business rule")
		assert.Contains(t, errMsg, "ChainFamilySelector", "v1.6 business rule")
		assert.Contains(t, errMsg, "GasPriceStalenessThreshold", "v1.6 business rule")

		// v1.6↔v2.0 cross-check (v2.0 values differ from v1.6 values)
		assert.Contains(t, errMsg, "v1.6<->v2.0", "v1.6↔v2.0 cross-check")

		// v2.0 business rules
		assert.Contains(t, errMsg, "LinkFeeMultiplierPercent", "v2.0 business rule")
	})

	// ===== SUBTEST 2: fixed values =====
	t.Run("fixed_values", func(t *testing.T) {
		// Read v1.5 OnRamp dynamic config to align v1.6 fields.
		callOpts := &bind.CallOpts{Context: t.Context()}
		state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)
		sourceState := state.MustGetEVMChainState(source)
		onRamp := sourceState.EVM2EVMOnRamp[dest]
		require.NotNil(t, onRamp)
		v15Cfg, err := onRamp.GetDynamicConfig(callOpts)
		require.NoError(t, err)

		// Fix v1.6 FeeQuoter dest config:
		// - All 11 cross-checked fields match v1.5 OnRamp dynamic config
		// - Business-rule fields set to expected values
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
			DefaultTokenFeeUSDCents:           v15Cfg.DefaultTokenFeeUSDCents,
			EnforceOutOfOrder:                 v15Cfg.EnforceOutOfOrder,
			DestGasPerPayloadByteBase:         uint8(v15Cfg.DestGasPerPayloadByte), //nolint:gosec // match v1.5 truncation
			// Business-rule fields
			DestGasPerPayloadByteHigh:      ccipevm.CalldataGasPerByteHigh,
			DestGasPerPayloadByteThreshold: ccipevm.CalldataGasPerByteThreshold,
			DefaultTxGasLimit:              200_000,
			NetworkFeeUSDCents:             10,
			GasPriceStalenessThreshold:     86400,
			ChainFamilySelector:            [4]byte{0x28, 0x12, 0xd5, 0x2c},
			GasMultiplierWeiPerEth:         1e18, // match v1.5 FeeTokenConfig
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

		// Fix v2.0 FeeQuoter dest config:
		// - All 10 v1.6↔v2.0 cross-checked fields match fixed v1.6 config
		// - Business-rule fields set to expected values
		goodV20Cfg := fqv2ops.DestChainConfig{
			IsEnabled:                   true,
			MaxDataBytes:                goodV16Cfg.MaxDataBytes,
			MaxPerMsgGasLimit:           goodV16Cfg.MaxPerMsgGasLimit,
			DestGasOverhead:             goodV16Cfg.DestGasOverhead,
			DestGasPerPayloadByteBase:   goodV16Cfg.DestGasPerPayloadByteBase,
			ChainFamilySelector:         goodV16Cfg.ChainFamilySelector,
			DefaultTokenFeeUSDCents:     goodV16Cfg.DefaultTokenFeeUSDCents,
			DefaultTokenDestGasOverhead: goodV16Cfg.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:           200_000,
			NetworkFeeUSDCents:          10,
			LinkFeeMultiplierPercent:    fqv2seq.LinkFeeMultiplierPercent,
		}
		updateV20FeeQuoterDestConfig(t, evmChain, fqV2Addr, dest, goodV20Cfg)

		// Fix v1.5 OnRamp FeeTokenConfig so v2.0↔v1.5 NetworkFeeUSDCents cross-check passes.
		// The v1.5 test default has NetworkFeeUSDCents=100 per-token, but expected per-dest is 10.
		updateV15OnRampFeeTokenConfig(t, evmChain, onRamp, linkAddr, wethAddr)

		// Reload state and validate.
		state, err = stateview.LoadOnchainState(tenv, stateview.WithLoadLegacyContracts(true))
		require.NoError(t, err)
		chainState := state.MustGetEVMChainState(source)

		err = chainState.ValidateFeeQuoter(tenv, source, connectedChains, fqV2, evmChain.Client)
		require.NoError(t, err, "validation must pass with correctly aligned configs")
	})
}

// --- Helpers ---

// deployV16Contracts deploys v1.6 HomeChain, LinkToken, MCMS, Prerequisites, and ChainContracts.
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

// deployV20FeeQuoter deploys a FeeQuoter v2.0 on the given chain.
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

// updateV20FeeQuoterFeeTokens calls applyFeeTokensUpdates on a v2.0 FeeQuoter.
func updateV20FeeQuoterFeeTokens(t *testing.T, evmChain cldf_evm.Chain, fqAddr common.Address, toAdd, toRemove []common.Address) {
	t.Helper()
	parsedABI, err := abi.JSON(strings.NewReader(fqv2ops.FeeQuoterABI))
	require.NoError(t, err)
	bc := bind.NewBoundContract(fqAddr, parsedABI, evmChain.Client, evmChain.Client, evmChain.Client)
	tx, err := bc.Transact(evmChain.DeployerKey, "applyFeeTokensUpdates", toRemove, toAdd)
	require.NoError(t, err, "applyFeeTokensUpdates")
	_, err = evmChain.Confirm(tx)
	require.NoError(t, err)
}

// updateV20FeeQuoterDestConfig calls applyDestChainConfigUpdates on a v2.0 FeeQuoter.
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

// updateV15OnRampFeeTokenConfig updates the v1.5 OnRamp fee token config so that
// NetworkFeeUSDCents=10 and GasMultiplierWeiPerEth=1e18 for all fee tokens.
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
