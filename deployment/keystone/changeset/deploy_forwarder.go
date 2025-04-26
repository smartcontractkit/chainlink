package changeset

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[DeployForwarderRequest] = DeployForwarder

type DeployForwarderRequest struct {
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
}

// DeployForwarder deploys the KeystoneForwarder contract to all chains in the environment
// callers must merge the output addressbook with the existing one
// TODO: add selectors to deploy only to specific chains
// Deprecated: use DeployForwarderV2 instead
func DeployForwarderX(env deployment.Environment, cfg DeployForwarderRequest) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	selectors := cfg.ChainSelectors
	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(env.Chains))
	}
	for _, sel := range selectors {
		chain, ok := env.Chains[sel]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", sel)
		}
		lggr.Infow("deploying forwarder", "chainSelector", chain.Selector)
		forwarderResp, err := internal.DeployForwarder(env.GetContext(), chain, ab)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", chain.Selector, err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", forwarderResp.Tv.String(), chain.Selector, forwarderResp.Address.String())
	}
	// convert all the addresses to t
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func DeployForwarder(env deployment.Environment, cfg DeployForwarderRequest) (deployment.ChangesetOutput, error) {
	var out deployment.ChangesetOutput
	out.AddressBook = deployment.NewMemoryAddressBook() //nolint:staticcheck // TODO CRE-400
	out.DataStore = datastore.NewMemoryDataStore[datastore.DefaultMetadata, datastore.DefaultMetadata]()

	selectors := cfg.ChainSelectors
	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(env.Chains))
	}

	for _, sel := range selectors {
		req := &DeployRequestV2{
			ChainSel: sel,
			deployFn: internal.DeployForwarder,
		}
		csOut, err := deploy(env, req)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy KeystoneForwarder to chain selector %d: %w", sel, err)
		}
		if err := out.AddressBook.Merge(csOut.AddressBook); err != nil { //nolint:staticcheck // TODO CRE-400
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge address book for chain selector %d: %w", sel, err)
		}
		if err := out.DataStore.Merge(csOut.DataStore.Seal()); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge datastore for chain selector %d: %w", sel, err)
		}
	}
	// convert all the addresses to t
	return out, nil
}

// DeployForwarderV2 deploys the KeystoneForwarder contract to the specified chain
func DeployForwarderV2(env deployment.Environment, req *DeployRequestV2) (deployment.ChangesetOutput, error) {
	req.deployFn = internal.DeployForwarder
	return deploy(env, req)
}

var _ deployment.ChangeSet[ConfigureForwardContractsRequest] = ConfigureForwardContracts

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

func ConfigureForwardContracts(env deployment.Environment, req ConfigureForwardContractsRequest) (deployment.ChangesetOutput, error) {
	wfDon, err := internal.NewRegisteredDon(env, internal.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create registered don: %w", err)
	}
	r, err := internal.ConfigureForwardContracts(&env, internal.ConfigureForwarderContractsRequest{
		Dons:    []internal.RegisteredDon{*wfDon},
		UseMCMS: req.UseMCMS(),
		Chains:  req.Chains,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure forward contracts: %w", err)
	}

	cresp, err := getContractSetsV2(env.Logger, getContractSetsRequestV2{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	var out deployment.ChangesetOutput
	if req.UseMCMS() {
		if len(r.OpsPerChain) == 0 {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		for chainSelector, op := range r.OpsPerChain {
			contracts := cresp.ContractSets[chainSelector]
			timelocksPerChain := map[uint64]string{
				chainSelector: contracts.Forwarder.McmsContracts.Timelock.Address().Hex(),
			}
			proposerMCMSes := map[uint64]string{
				chainSelector: contracts.Forwarder.McmsContracts.ProposerMcm.Address().Hex(),
			}
			inspector, err := proposalutils.McmsInspectorForChain(env, chainSelector)
			if err != nil {
				return deployment.ChangesetOutput{}, err
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
