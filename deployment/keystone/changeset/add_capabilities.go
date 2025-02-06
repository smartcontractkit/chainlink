package changeset

import (
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

// AddCapabilitiesRequest is a request to add capabilities
type AddCapabilitiesRequest struct {
	RegistryChainSel uint64

	Capabilities []kcr.CapabilitiesRegistryCapability
	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
}

var _ deployment.ChangeSet[*AddCapabilitiesRequest] = AddCapabilities

// AddCapabilities is a deployment.ChangeSet that adds capabilities to the capabilities registry
//
// It is idempotent. It deduplicates the input capabilities.
//
// When using MCMS, the output will contain a single proposal with a single batch containing all capabilities to be added.
// When not using MCMS, each capability will be added in a separate transaction.
func AddCapabilities(env deployment.Environment, req *AddCapabilitiesRequest) (deployment.ChangesetOutput, error) {
	registryChain, ok := env.Chains[req.RegistryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	cs, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      map[uint64]deployment.Chain{req.RegistryChainSel: registryChain},
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contractSet, exists := cs.ContractSets[req.RegistryChainSel]
	if !exists {
		return deployment.ChangesetOutput{}, fmt.Errorf("contract set not found for chain %d", req.RegistryChainSel)
	}
	useMCMS := req.MCMSConfig != nil
	ops, err := internal.AddCapabilities(env.Logger, contractSet.CapabilitiesRegistry, env.Chains[req.RegistryChainSel], req.Capabilities, useMCMS)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to add capabilities: %w", err)
	}
	out := deployment.ChangesetOutput{}
	if useMCMS {
		if ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		timelocksPerChain := map[uint64]gethcommon.Address{
			registryChain.Selector: contractSet.Timelock.Address(),
		}
		proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
			registryChain.Selector: contractSet.ProposerMcm,
		}

		proposal, err := proposalutils.BuildProposalFromBatches(
			timelocksPerChain,
			proposerMCMSes,
			[]timelock.BatchChainOperation{*ops},
			"proposal to add capabilities",
			req.MCMSConfig.MinDuration,
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.Proposals = []timelock.MCMSWithTimelockProposal{*proposal}
	}
	return out, nil
}
