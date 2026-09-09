package testhelpers

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	_ "github.com/smartcontractkit/chainlink-ccip/chains/solana/deployment/v1_6_0/sequences" // register Solana deploy/lane/mcms adapters
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	cs_ccip "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	_ "github.com/smartcontractkit/chainlink-sui/deployment/adapters" // register Sui MCMS/curse/token/fee adapters (init)
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	suilanes "github.com/smartcontractkit/chainlink-sui/deployment/lanes" // registers SuiAdapter and provides WithConnectChainsEnvironment
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
	sui_utils "github.com/smartcontractkit/chainlink-sui/deployment/utils"

	ccipChangeSetSolanaV0_1_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// addSuiSolanaMixedLane wires a Sui<->Solana lane via ConnectChains. Unlike the Aptos mixed
// lane, there is no EVM leg to reconcile, so no ensureEVMRouterLaneConfig fallback is needed.
// ConnectChains is bidirectional: it configures both Sui->Solana and Solana->Sui legs using the
// SuiAdapter and SolanaAdapter registered by the chainlink-sui lanes package and the Solana
// v1_6_0 sequences package (blank-imported above).
//
// The SuiAdapter resolves Sui package IDs from the active CLDF environment rather than the
// datastore, so ConnectChains MUST run inside suilanes.WithConnectChainsEnvironment (see
// chainlink-sui/deployment/lanes/env.go). Other chain families are unaffected by that scope.
func addSuiSolanaMixedLane(
	t *testing.T,
	e *DeployedEnv,
	state stateview.CCIPOnChainState,
	from, to uint64,
	fromFamily, toFamily string,
	isTestRouter bool,
	gasPrices map[uint64]*big.Int,
	tokenPrices map[string]*big.Int,
	fqCfg fee_quoter.FeeQuoterDestChainConfig,
) error {
	// Resolve Sui/Solana selectors up front: the per-family price seeding below keys off them
	// and (for the Sui source) must run before the lane changesets are built.
	suiSel := to
	if fromFamily == chainsel.FamilySui {
		suiSel = from
	}
	solSel := to
	if fromFamily == chainsel.FamilySolana {
		solSel = from
	}

	// Sui source — seed the LINK source-token USD price EOA, BEFORE the ownership transfer to
	// MCMS consumes the deployer's CCIPOwnerCapObjectId. The Sui lanes adapter seeds the dest
	// GAS price during ConnectChains (connect_chains_source.go:145), but nobody seeds the LINK
	// usd_per_token, so get_fee for a Sui->Solana message (fee paid in LINK) would abort. EOA
	// mode (TimelockConfig nil) runs the same FeeQuoterUpdatePricesWithOwnerCapOp legacy
	// EVM<->Sui tests use; its VerifyPreconditions needs no dest config, so it works pre-lane.
	// The seeded price persists through the ownership transfer and ConnectChains (neither
	// touches the fee-quoter price tables). Token price only — gas is ConnectChains's job.
	//
	// ConnectChains seeds the Solana dest gas from gasPrices[solSel]. The FamilySui default
	// (1e17) is an untested placeholder that legacy tests silently overwrite via SendRequestSui's
	// EOA update (skipped here under MCMS ownership). Pin it to the known-good value
	// SendRequestSui uses (bigIntGasUsdPerUnitGas, test_sui_helpers.go:162) so the lane-seeded
	// gas yields the same payable fee as the legacy path.
	if fromFamily == chainsel.FamilySui {
		solDestGasUsd, ok := new(big.Int).SetString("41946474500", 10)
		require.True(t, ok, "parse Solana dest gas USD price")
		gasPrices[solSel] = solDestGasUsd
		if err := seedSuiSourceTokenPriceEOA(t, e, suiSel, state); err != nil {
			return fmt.Errorf("seed Sui source LINK token price (EOA): %w", err)
		}
	}

	changesets := addSuiSolanaLaneChangesets(t, from, to, isTestRouter, gasPrices, tokenPrices, fqCfg)

	// DeploySuiChain deploys MCMS with nil role configs (quorum=0) and only INITIATES the
	// ownership transfer of the four CCIP contracts (Router/CCIP state/OnRamp/OffRamp) to MCMS.
	// Unlike Aptos (whose deploy changeset configures MCMS quorum AND accepts+registers CCIP
	// ownership in one shot), the Sui deploy leaves CCIP unregistered in the MCMS registry. That
	// breaks lanes.ConnectChains, whose Sui ops are dispatched by MCMS via
	// mcms_registry::get_registered_proof_type -> abort 7 (EPackageNotRegistered) until each CCIP
	// package self-registers inside execute_ownership_transfer_to_mcms (ownable.move:244).
	//
	// Complete the MCMS governance of CCIP here, mirroring integration-tests/mcms/common.go:
	// ConfigureMCMS (quorum=1, EOA) -> AcceptOwnershipCCIP (MCMS bypasser proposal, auto-executed
	// by ApplyChangesets) -> ExecuteOwnershipTransferToMcms (EOA finalize + register_entrypoint).
	// This is scoped to Sui<->Solana lanes only; legacy EVM<->Sui tests never call this helper, so
	// their deployer-owned/quorum=0 EOA-executed Sui path is untouched.
	if err := completeSuiCCIPMCMSOwnership(t, e, suiSel); err != nil {
		return fmt.Errorf("complete Sui CCIP MCMS ownership: %w", err)
	}

	// The legacy solana_v0_1_1 env deploy registers Solana CCIP core contracts at
	// deployment.Version1_0_0 (1.0.0), but the chainlink-ccip lanes SolanaAdapter queries them at
	// 1.6.0 (router.Version / offramp.Version / fee_quoter.Version / rmn_remote.Version). Without
	// parallel 1.6.0 refs, ConnectChains -> populateAddresses -> GetOnRampAddress finds 0 refs
	// ("expected to find exactly 1 ref ... found 0"). Add 1.6.0 refs for the Solana leg here.
	// AddressRef keys include version, so the 1.6.0 refs coexist with the 1.0.0 refs (no
	// overwrite); stateview resolves newest-by-type and is unaffected. This mirrors the Aptos
	// mixed lane, whose legacy env deploy already saves its CCIP ref at the 1.6.0 version its
	// lanes adapter queries.
	var err error
	e.Env, err = reregisterSolanaCCIPRefsAtLanesVersion(e.Env, solSel)
	if err != nil {
		return fmt.Errorf("re-register Solana CCIP refs at lanes version 1.6.0: %w", err)
	}

	// SuiAdapter address getters read chain metadata from the env scope set below; run the
	// ConnectChains changeset application inside it.
	if err := suilanes.WithConnectChainsEnvironment(e.Env, func() error {
		e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, changesets)
		return err
	}); err != nil {
		return fmt.Errorf("connect Sui<->Solana lane: %w", err)
	}

	// Solana source — the Solana lanes adapter (ConfigureLaneLegAsSource) writes only the
	// FeeQuoter dest-chain config; it seeds NO gas prices, unlike the Sui/Aptos adapters. The
	// DON would normally push the Sui dest gas price via price-only OCR commits, but the test
	// OCR config has no Sui gas-price feed, so GetFee aborts StaleGasPrice (code 8024). Seed it
	// via the Solana UpdatePrices escape hatch AFTER ConnectChains (UpdatePrices.Validate
	// requires the dest config to exist). This smoke-test env deploys Solana CCIP with
	// preload=true, so mcmsCfg=nil at deploy and the FeeQuoter is DEPLOYER-owned (EOA), not
	// timelock-owned — there is no RBACTimelock PDA entry in the address book, so an MCMS
	// proposal fails ValidateSolana ("RBACTimelock not present on the chain"). Both Solana
	// billing changesets pick MCMS-vs-EOA from the FeeQuoter's ACTUAL on-chain ownership
	// (IsSolanaProgramOwnedByTimelock), not from cfg.MCMS; with MCMS=nil they Validate as EOA
	// (ValidateOwnershipSolana checks deployer ownership) and Apply via chain.Confirm with the
	// deployer signer. Mirrors the proven pattern in solana_v0_1_1/cs_chain_contracts_test.go
	// (add the timelock signer as a price updater, then push the gas price), but run EOA. The Sui
	// source leg needs no post-ConnectChains seeding: its dest gas was seeded by the Sui adapter
	// and its LINK token price was seeded EOA above.
	if fromFamily == chainsel.FamilySolana {
		if err := seedSolanaSourceSuiDestGasPriceEOA(t, e, solSel, suiSel, gasPrices); err != nil {
			return fmt.Errorf("seed Solana source Sui dest gas price (EOA): %w", err)
		}
	}

	return nil
}

// seedSuiSourceTokenPriceEOA seeds the Sui fee-quoter's LINK source-token USD price EOA, using
// the deployer's CCIPOwnerCapObjectId. It MUST run BEFORE completeSuiCCIPMCMSOwnership, which
// moves the OwnerCap into the MCMS registry (after which the EOA op aborts). The Sui lanes
// adapter seeds the dest gas price during ConnectChains but never the source-token price, so
// without this get_fee for a Sui-source message paid in LINK aborts. Token price only.
//
// This is the EOA mode of sui_cs.SeedDestChainPrices (TimelockConfig nil -> deployer signer); it
// calls the same FeeQuoterUpdatePricesWithOwnerCapOp legacy EVM<->Sui tests use. Apply is called
// DIRECTLY, not via ApplyChangesets: like ConfigureMCMS, SeedDestChainPrices returns a zero-value
// placeholder TimelockProposal even in EOA mode (cs_seed_dest_chain_prices.go:161), which
// ApplyChangesets would try to auto-execute and fail validation. The value matches
// SendSuiCCIPRequest's bigIntSourceUsdPerToken (test_sui_helpers.go:157) exactly.
func seedSuiSourceTokenPriceEOA(t *testing.T, e *DeployedEnv, suiSel uint64, state stateview.CCIPOnChainState) error {
	t.Helper()
	linkUsd, ok := new(big.Int).SetString("15377040000000000000000000000", 10)
	require.True(t, ok, "parse Sui source LINK USD price")
	suiState, ok := state.SuiChains[suiSel]
	require.True(t, ok, "no Sui chain state for selector %d", suiSel)
	out, err := sui_cs.SeedDestChainPrices{}.Apply(e.Env, sui_cs.SeedDestChainPricesConfig{
		SuiChainSelector:  suiSel,
		SourceTokens:      []string{suiState.LinkTokenCoinMetadataId},
		SourceUsdPerToken: []*big.Int{linkUsd},
	})
	if err != nil {
		return fmt.Errorf("seed Sui source LINK token price via EOA: %w", err)
	}
	_ = out
	return nil
}

// seedSolanaSourceSuiDestGasPriceEOA seeds the Sui dest gas price on the Solana source
// fee-quoter via the Solana UpdatePrices escape hatch. The Solana lanes adapter writes only the
// dest-chain config, not gas prices. This smoke-test env deploys Solana CCIP with preload=true, so
// the FeeQuoter is DEPLOYER-owned (EOA) and there is no RBACTimelock PDA entry in the address book
// — an MCMS proposal would fail ValidateSolana ("RBACTimelock not present on the chain"). Both
// Solana billing changesets dispatch MCMS-vs-EOA from the FeeQuoter's ACTUAL on-chain ownership
// (IsSolanaProgramOwnedByTimelock), not from cfg.MCMS; passing MCMS=nil makes Validate check EOA
// ownership (which holds) and Apply use chain.Confirm with the deployer signer. Must run AFTER
// ConnectChains (UpdatePrices.Validate requires the dest config to exist).
//
// PriceUpdater MUST equal the FeeQuoter's authority (owner). UpdatePrices carries no price-updater
// pubkey in its instruction args, so the on-chain program seeds allowed_price_updater from the
// authority signer; a mismatch aborts ConstraintSeeds 0x7d6. In the preload EOA env the authority
// is the deployer key — NOT FetchTimelockSigner, which is the authority only after
// TransferOwnership moves the FeeQuoter to the MCMS timelock a path the preload env never runs.
// Mirrors cs_chain_contracts_test.go:345 EOA branch: testPriceUpdater = DeployerKey.PublicKey().
func seedSolanaSourceSuiDestGasPriceEOA(t *testing.T, e *DeployedEnv, solSel, suiSel uint64, gasPrices map[uint64]*big.Int) error {
	t.Helper()
	priceUpdater := e.Env.BlockChains.SolanaChains()[solSel].DeployerKey.PublicKey()
	suiDestGasUsd := gasPrices[suiSel]
	if suiDestGasUsd == nil {
		return fmt.Errorf("no gas price provided for Sui dest chain %d", suiSel)
	}
	var err error
	e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.ModifyPriceUpdater),
			ccipChangeSetSolanaV0_1_1.ModifyPriceUpdaterConfig{
				ChainSelector:      solSel,
				PriceUpdater:       priceUpdater,
				PriceUpdaterAction: ccipChangeSetSolanaV0_1_1.AddUpdater,
				MCMS:               nil,
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.UpdatePrices),
			ccipChangeSetSolanaV0_1_1.UpdatePricesConfig{
				ChainSelector: solSel,
				GasPriceUpdates: []solFeeQuoter.GasPriceUpdate{
					{DestChainSelector: suiSel, UsdPerUnitGas: solCommonUtil.To28BytesBE(suiDestGasUsd.Uint64())},
				},
				PriceUpdater: priceUpdater,
				MCMS:         nil,
			},
		),
	})
	if err != nil {
		return fmt.Errorf("seed Solana source Sui dest gas price EOA: %w", err)
	}
	return nil
}

// completeSuiCCIPMCMSOwnership finishes the MCMS governance of the Sui CCIP contracts that
// DeploySuiChain only initiates. It mirrors integration-tests/mcms/common.go (the authoritative
// passing Sui MCMS-governance setup), in the same order:
//
//  1. ConfigureMCMS with quorum=1 for every timelock role (Proposer/Bypasser/Canceller). DeploySuiChain
//     deploys MCMS with nil role configs (quorum=0); without this, any MCMS proposal validated against
//     the on-chain Sui config fails "Quorum must be greater than 0". The changeset is applied directly
//     (not via ApplyChangesets) with a nil TimelockConfig so SetConfig runs EOA and the zero-value
//     placeholder timelock proposal it unconditionally returns is discarded (see the zero-proposal
//     gotcha: ApplyChangesets would try to auto-execute that placeholder and fail validation).
//
//  2. AcceptOwnershipCCIP, an MCMS bypasser timelock proposal (MinDelay=0) that makes MCMS accept the
//     pending ownership transfer DeploySuiChain initiated for all four CCIP contracts. Routed through
//     ApplyChangesets so the returned proposal is signed + auto-executed (SignMCMSTimelockProposal).
//     Acceptance is callable before CCIP registers because mcms_accept_ownership takes the MCMS address
//     as an argument rather than dispatching via the registry.
//
//  3. ExecuteOwnershipTransferToMcms, run EOA with the deployer's OwnerCaps, finalizes each transfer
//     and calls mcms_registry::register_entrypoint inside execute_ownership_transfer_to_mcms. This is
//     the step that actually registers CCIP in the MCMS registry, unblocking MCMS-dispatched lane ops.
func completeSuiCCIPMCMSOwnership(t *testing.T, e *DeployedEnv, suiSel uint64) error {
	t.Helper()

	// 1. Configure MCMS quorum=1 (EOA direct; discard the zero placeholder proposal).
	proposerCfg := cldftesthelpers.SingleGroupMCMS(t)
	bypasserCfg := cldftesthelpers.SingleGroupMCMS(t)
	cancellerCfg := cldftesthelpers.SingleGroupMCMS(t)
	// Applied directly (not via ApplyChangesets) with a nil TimelockConfig so SetConfig runs EOA
	// via the deployer signer and the zero-value placeholder timelock proposal ConfigureMCMS
	// unconditionally returns is discarded. (A composite literal in an if-init would need
	// parenthesizing, so assign to a variable first, matching the AcceptOwnershipEOA{} convention.)
	mcmsOut, err := sui_cs.ConfigureMCMS{}.Apply(e.Env, sui_cs.ConfigureMCMSConfig{
		ConfigureMCMSSeqInput: mcmsops.ConfigureMCMSSeqInput{
			ChainSelector: suiSel,
			Proposer:      &proposerCfg,
			Bypasser:      &bypasserCfg,
			Canceller:     &cancellerCfg,
		},
		// TimelockConfig nil -> SetConfig executes directly via the deployer signer.
	})
	if err != nil {
		return fmt.Errorf("configure MCMS (quorum) for Sui chain %d: %w", suiSel, err)
	}
	_ = mcmsOut

	// 2. Accept CCIP ownership via an MCMS bypasser proposal (auto-executed by ApplyChangesets).
	e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.AcceptOwnershipCCIP{}, sui_cs.AcceptOwnershipCCIPConfig{
			SuiChainSelector: suiSel,
			TimelockConfig: sui_utils.TimelockConfig{
				MCMSAction:   mcmstypes.TimelockActionBypass,
				MinDelay:     0,
				OverrideRoot: false,
			},
		}),
	})
	if err != nil {
		return fmt.Errorf("accept CCIP ownership for Sui chain %d: %w", suiSel, err)
	}

	// 3. Execute the ownership transfer to MCMS (EOA) so CCIP self-registers in the MCMS registry.
	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	if err != nil {
		return fmt.Errorf("load Sui onchain state: %w", err)
	}
	st, ok := suiState[suiSel]
	if !ok {
		return fmt.Errorf("no Sui state for chain selector %d", suiSel)
	}

	suiChain := e.Env.BlockChains.SuiChains()[suiSel]
	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			gasBudget := uint64(1_000_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &gasBudget,
			}
		},
		SuiRPC: suiChain.URL,
	}

	if _, err := operations.ExecuteSequence(e.Env.OperationsBundle, ownershipops.ExecuteOwnershipTransferToMcmsSequence, deps, ownershipops.ExecuteOwnershipTransferToMcmsSeqInput{
		StateObject: &ccipops.ExecuteOwnershipTransferToMcmsStateObjectInput{
			CCIPPackageId:         st.CCIPAddress,
			CCIPObjectRefObjectId: st.CCIPObjectRef,
			OwnerCapObjectId:      st.CCIPOwnerCapObjectId,
			RegistryObjectId:      st.MCMSRegistryObjectID,
			To:                    st.MCMSPackageID,
		},
		OnRamp: &onrampops.ExecuteOwnershipTransferToMcmsOnRampInput{
			OnRampPackageId:     st.OnRampAddress,
			OnRampRefObjectId:   st.CCIPObjectRef,
			OwnerCapObjectId:    st.OnRampOwnerCapObjectId,
			OnRampStateObjectId: st.OnRampStateObjectId,
			RegistryObjectId:    st.MCMSRegistryObjectID,
			To:                  st.MCMSPackageID,
		},
		OffRamp: &offrampops.ExecuteOwnershipTransferToMcmsOffRampInput{
			OffRampPackageId:     st.OffRampAddress,
			OffRampRefObjectId:   st.CCIPObjectRef,
			OwnerCapObjectId:     st.OffRampOwnerCapId,
			OffRampStateObjectId: st.OffRampStateObjectId,
			RegistryObjectId:     st.MCMSRegistryObjectID,
			To:                   st.MCMSPackageID,
		},
		Router: &routerops.ExecuteOwnershipTransferToMcmsRouterInput{
			RouterPackageId:     st.CCIPRouterAddress,
			OwnerCapObjectId:    st.CCIPRouterOwnerCapObjectId,
			RouterStateObjectId: st.CCIPRouterStateObjectID,
			RegistryObjectId:    st.MCMSRegistryObjectID,
			To:                  st.MCMSPackageID,
		},
	}); err != nil {
		return fmt.Errorf("execute ownership transfer to MCMS (register CCIP) for Sui chain %d: %w", suiSel, err)
	}

	return nil
}

// reregisterSolanaCCIPRefsAtLanesVersion copies the Solana chain's CCIP core contract refs
// (Router, OffRamp, FeeQuoter, RMNRemote) into the env datastore at version 1.6.0, the version
// the chainlink-ccip lanes SolanaAdapter queries them at. The legacy solana_v0_1_1 env deploy
// registers them at 1.0.0; this adds parallel 1.6.0 refs (AddressRef keys include version, so
// the 1.0.0 refs are untouched) so lanes.ConnectChains can resolve them.
//
// It is a test-only shim: it does not mutate ExistingAddresses, does not touch the Sui leg, and
// the merged datastore carries forward every prior ref. ApplyChangesets later unions its own
// output over this datastore, so the 1.6.0 refs survive into ConnectChains.
func reregisterSolanaCCIPRefsAtLanesVersion(env cldf.Environment, solSel uint64) (cldf.Environment, error) {
	lanesVersion := semver.MustParse("1.6.0")
	coreTypes := map[cldf.ContractType]bool{
		shared.Router:    true,
		shared.OffRamp:   true,
		shared.FeeQuoter: true,
		shared.RMNRemote: true,
	}

	ds := datastore.NewMemoryDataStore()
	if err := ds.Merge(env.DataStore); err != nil {
		return env, fmt.Errorf("merge env datastore: %w", err)
	}

	refs, err := env.ExistingAddresses.AddressesForChain(solSel)
	if err != nil {
		return env, fmt.Errorf("read Solana chain %d refs: %w", solSel, err)
	}
	for addr, tv := range refs {
		if !coreTypes[tv.Type] {
			continue
		}
		ref := datastore.AddressRef{
			ChainSelector: solSel,
			Address:       addr,
			Type:          datastore.ContractType(tv.Type),
			Version:       lanesVersion,
		}
		if err := ds.Addresses().Add(ref); err != nil {
			return env, fmt.Errorf("add 1.6.0 ref for %s %s: %w", tv.Type, addr, err)
		}
	}

	env.DataStore = ds.Seal()
	return env, nil
}

// addSuiSolanaLaneChangesets connects a Sui chain with its Solana peer via the generic
// lanes.ConnectChains changeset, which configures both legs of the lane using the adapters
// registered by the chainlink-sui lanes package and the chainlink-ccip Solana v1_6_0 sequences
// package. Gas prices are keyed by chain selector (a missing entry falls back to the adapter
// default). Token prices are attached to the source chain definition; the source leg uses them
// when seeding FeeQuoter prices. The FeeQuoterDestChainConfig override is attached to the
// destination definition, mirroring addLaneAptosChangesets.
func addSuiSolanaLaneChangesets(
	t *testing.T,
	srcChainSelector, destChainSelector uint64,
	isTestRouter bool,
	gasPrices map[uint64]*big.Int,
	tokenPrices map[string]*big.Int,
	fqCfg fee_quoter.FeeQuoterDestChainConfig,
) []commoncs.ConfiguredChangeSet {
	srcFamily, err := chainsel.GetSelectorFamily(srcChainSelector)
	require.NoError(t, err)
	destFamily, err := chainsel.GetSelectorFamily(destChainSelector)
	require.NoError(t, err)

	if (srcFamily != chainsel.FamilySui && srcFamily != chainsel.FamilySolana) ||
		(destFamily != chainsel.FamilySui && destFamily != chainsel.FamilySolana) {
		t.Fatalf("addSuiSolanaLaneChangesets requires a Sui<->Solana lane. srcFamily: %v destFamily: %v", srcFamily, destFamily)
	}

	destFQOverride := feeQuoterDestChainConfigOverride(fqCfg)

	makeDefinition := func(selector uint64) lanes.ChainDefinition {
		definition := lanes.ChainDefinition{
			Selector:               selector,
			GasPrice:               gasPrices[selector],
			RMNVerificationEnabled: false,
			AllowListEnabled:       false,
		}
		if selector == srcChainSelector {
			definition.TokenPrices = tokenPrices
		}
		if selector == destChainSelector {
			definition.FeeQuoterDestChainConfigOverrides = &destFQOverride
		}
		return definition
	}

	validUntil, err := mcmsValidUntil(time.Now().Add(24 * time.Hour))
	require.NoError(t, err)

	return []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			lanes.ConnectChains(lanes.GetLaneAdapterRegistry(), cs_ccip.GetRegistry()),
			lanes.ConnectChainsConfig{
				MCMS: ccipmcms.Input{
					ValidUntil:     validUntil,
					TimelockDelay:  mcmstypes.NewDuration(time.Second),
					TimelockAction: mcmstypes.TimelockActionSchedule,
				},
				Lanes: []lanes.LaneConfig{
					{
						Version:    semver.MustParse("1.6.0"),
						ChainA:     makeDefinition(srcChainSelector),
						ChainB:     makeDefinition(destChainSelector),
						TestRouter: isTestRouter,
					},
				},
			},
		),
	}
}
