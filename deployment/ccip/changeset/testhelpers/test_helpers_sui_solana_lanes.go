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
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	cs_ccip "github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	_ "github.com/smartcontractkit/chainlink-sui/deployment/adapters"      // register Sui MCMS/curse/token/fee adapters (init)
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
	suiSel := to
	if fromFamily == chainsel.FamilySui {
		suiSel = from
	}
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
	solSel := to
	if fromFamily == chainsel.FamilySolana {
		solSel = from
	}
	var err error
	e.Env, err = reregisterSolanaCCIPRefsAtLanesVersion(e.Env, solSel)
	if err != nil {
		return fmt.Errorf("re-register Solana CCIP refs at lanes version 1.6.0: %w", err)
	}

	// SuiAdapter address getters read chain metadata from the env scope set below; run the
	// ConnectChains changeset application inside it.
	return suilanes.WithConnectChainsEnvironment(e.Env, func() error {
		e.Env, _, err = commoncs.ApplyChangesets(t, e.Env, changesets)
		return err
	})
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
	if _, err := sui_cs.ConfigureMCMS{}.Apply(e.Env, sui_cs.ConfigureMCMSConfig{
		ConfigureMCMSSeqInput: mcmsops.ConfigureMCMSSeqInput{
			ChainSelector: suiSel,
			Proposer:      &proposerCfg,
			Bypasser:      &bypasserCfg,
			Canceller:     &cancellerCfg,
		},
		// TimelockConfig nil -> SetConfig executes directly via the deployer signer.
	}); err != nil {
		return fmt.Errorf("configure MCMS (quorum) for Sui chain %d: %w", suiSel, err)
	}

	// 2. Accept CCIP ownership via an MCMS bypasser proposal (auto-executed by ApplyChangesets).
	var err error
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
