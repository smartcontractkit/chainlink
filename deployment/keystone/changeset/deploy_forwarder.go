package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

var _ deployment.ChangeSet[uint64] = DeployForwarder

// DeployForwarder deploys the KeystoneForwarder contract to all chains in the environment
// callers must merge the output addressbook with the existing one
// TODO: add selectors to deploy only to specific chains
func DeployForwarder(env deployment.Environment, _ uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	for _, chain := range env.Chains {
		lggr.Infow("deploying forwarder", "chainSelector", chain.Selector)
		forwarderResp, err := kslib.DeployForwarder(chain, ab)
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
		return fmt.Errorf("WFNodeIDs must not be empty")
	}
	return nil
}

func (r ConfigureForwardContractsRequest) UseMCMS() bool {
	return r.MCMSConfig != nil
}

func ConfigureForwardContracts(env deployment.Environment, req ConfigureForwardContractsRequest) (deployment.ChangesetOutput, error) {
	wfDon, err := kslib.NewRegisteredDon(env, kslib.RegisteredDonConfig{
		NodeIDs:          req.WFNodeIDs,
		Name:             req.WFDonName,
		RegistryChainSel: req.RegistryChainSel,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create registered don: %w", err)
	}
	r, err := kslib.ConfigureForwardContracts(&env, kslib.ConfigureForwarderContractsRequest{
		Dons:    []kslib.RegisteredDon{*wfDon},
		UseMCMS: req.UseMCMS(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure forward contracts: %w", err)
	}

	cresp, err := kslib.GetContractSets(env.Logger, &kslib.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	var out deployment.ChangesetOutput
	if req.UseMCMS() {
		if len(r.OpsPerChain) == 0 {
			return out, fmt.Errorf("expected MCMS operation to be non-nil")
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
