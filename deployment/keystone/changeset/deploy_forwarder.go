package changeset

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[DeployForwarderRequest] = DeployForwarder

type DeployForwarderRequest struct {
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
}

// DeployForwarder deploys the KeystoneForwarder contract to all chains in the environment
// callers must merge the output addressbook with the existing one
// TODO: add selectors to deploy only to specific chains
// Deprecated: use DeployForwarderV2 instead
func DeployForwarderX(env cldf.Environment, cfg DeployForwarderRequest) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	ab := cldf.NewMemoryAddressBook()
	selectors := cfg.ChainSelectors
	evmChains := env.BlockChains.EVMChains()

	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(evmChains))
	}
	for _, sel := range selectors {
		chain, ok := evmChains[sel]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", sel)
		}
		lggr.Infow("deploying forwarder", "chainSelector", chain.Selector)
		forwarderResp, err := internal.DeployForwarder(env.GetContext(), chain, ab)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", chain.Selector, err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", forwarderResp.Tv.String(), chain.Selector, forwarderResp.Address.String())
	}
	// convert all the addresses to t
	return cldf.ChangesetOutput{AddressBook: ab}, nil
}

func DeployForwarder(env cldf.Environment, cfg DeployForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	out.AddressBook = cldf.NewMemoryAddressBook() //nolint:staticcheck // TODO CRE-400
	out.DataStore = datastore.NewMemoryDataStore()

	selectors := cfg.ChainSelectors
	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(env.BlockChains.EVMChains()))
	}

	for _, sel := range selectors {
		req := &DeployRequestV2{
			ChainSel: sel,
			deployFn: internal.DeployForwarder,
		}
		csOut, err := deploy(env, req)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", sel, err)
		}
		if err := out.AddressBook.Merge(csOut.AddressBook); err != nil { //nolint:staticcheck // TODO CRE-400
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge address book for chain selector %d: %w", sel, err)
		}
		if err := out.DataStore.Merge(csOut.DataStore.Seal()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge datastore for chain selector %d: %w", sel, err)
		}
	}
	// convert all the addresses to t
	return out, nil
}

// DeployForwarderV2 deploys the KeystoneForwarder contract to the specified chain
func DeployForwarderV2(env cldf.Environment, req *DeployRequestV2) (cldf.ChangesetOutput, error) {
	req.deployFn = internal.DeployForwarder
	return deploy(env, req)
}

var _ cldf.ChangeSet[ConfigureForwardContractsRequest] = ConfigureForwardContracts

type ConfigureForwardContractsRequest struct {
	WFDonName string
	// workflow don node ids in the offchain client. Used to fetch and derive the signer keys
	WFNodeIDs        []string
	RegistryChainSel uint64

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains map[uint64]struct{}
}

func (r ConfigureForwardContractsRequest) Validate() error {
	if len(r.WFNodeIDs) == 0 {
		return errors.New("WFNodeIDs must not be empty")
	}
	return nil
}

func (r ConfigureForwardContractsRequest) UseMCMS() bool {
	return r.MCMSConfig != nil
}

func ConfigureForwardContracts(env cldf.Environment, req ConfigureForwardContractsRequest) (cldf.ChangesetOutput, error) {
	wfDon, err := internal.NewRegisteredDon(env, internal.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create registered don: %w", err)
	}
	r, err := internal.ConfigureForwardContracts(&env, internal.ConfigureForwarderContractsRequest{
		Dons:    []internal.RegisteredDon{*wfDon},
		UseMCMS: req.UseMCMS(),
		Chains:  req.Chains,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure forward contracts: %w", err)
	}

	var out cldf.ChangesetOutput
	if req.UseMCMS() {
		if len(r.OpsPerChain) == 0 {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		for chainSelector, op := range r.OpsPerChain {
			fwrAddr, ok := r.ForwarderAddresses[chainSelector]
			if !ok {
				return out, fmt.Errorf("expected configured forwarder address for chain selector %d", chainSelector)
			}
			fwr, err := GetOwnedContractV2[*forwarder.KeystoneForwarder](env.DataStore.Addresses(), env.BlockChains.EVMChains()[chainSelector], fwrAddr.String())
			if err != nil {
				return out, fmt.Errorf("failed to get forwarder contract for chain selector %d: %w", chainSelector, err)
			}
			if fwr.McmsContracts == nil {
				return out, fmt.Errorf("expected forwarder contract %s to be owned by MCMS for chain selector %d", fwrAddr.String(), chainSelector)
			}
			timelocksPerChain := map[uint64]string{
				chainSelector: fwr.McmsContracts.Timelock.Address().Hex(),
			}
			proposerMCMSes := map[uint64]string{
				chainSelector: fwr.McmsContracts.ProposerMcm.Address().Hex(),
			}
			inspector, err := proposalutils.McmsInspectorForChain(env, chainSelector)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			inspectorPerChain := map[uint64]mcmssdk.Inspector{
				chainSelector: inspector,
			}

			proposal, err := proposalutils.BuildProposalFromBatchesV2(
				env,
				timelocksPerChain,
				proposerMCMSes,
				inspectorPerChain,
				[]mcmstypes.BatchOperation{op},
				"proposal to set forwarder config",
				proposalutils.TimelockConfig{
					MinDelay: req.MCMSConfig.MinDuration,
				},
			)
			if err != nil {
				return out, fmt.Errorf("failed to build proposal: %w", err)
			}
			out.MCMSTimelockProposals = append(out.MCMSTimelockProposals, *proposal)
		}
	}
	return out, nil
}
