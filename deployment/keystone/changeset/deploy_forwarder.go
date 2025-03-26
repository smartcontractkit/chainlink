package changeset

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

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
func DeployForwarder(env deployment.Environment, cfg DeployForwarderRequest) (deployment.ChangesetOutput, error) {
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

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

var _ deployment.ChangeSet[ConfigureForwardContractsRequest] = ConfigureForwardContracts

type ConfigureForwardContractsRequest struct {
	WFDonName string
	// workflow don node ids in the offchain client. Used to fetch and derive the signer keys
	WFNodeIDs        []string
	RegistryChainSel uint64

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
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
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure forward contracts: %w", err)
	}

	cresp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
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
			timelocksPerChain := map[uint64]common.Address{
				chainSelector: contracts.Timelock.Address(),
			}
			proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
				chainSelector: contracts.ProposerMcm,
			}

			proposal, err := proposalutils.BuildProposalFromBatches(
				timelocksPerChain,
				proposerMCMSes,
				[]timelock.BatchChainOperation{op},
				"proposal to set forwarder config",
				req.MCMSConfig.MinDuration,
			)
			if err != nil {
				return out, fmt.Errorf("failed to build proposal: %w", err)
			}
			out.Proposals = append(out.Proposals, *proposal)
		}
	}
	return out, nil
}
