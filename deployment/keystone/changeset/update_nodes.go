package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"

	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ deployment.ChangeSet[*UpdateNodesRequest] = UpdateNodes

type UpdateNodesRequest struct {
	RegistryChainSel uint64
	P2pToUpdates     map[p2pkey.PeerID]NodeUpdate

	UseMCMS bool
}
type NodeUpdate = internal.NodeUpdate

// UpdateNodes updates the a set of nodes.
// This a complex action in practice that involves registering missing capabilities, adding the nodes, and updating
// the capabilities of the DON
func UpdateNodes(env deployment.Environment, req *UpdateNodesRequest) (deployment.ChangesetOutput, error) {
	// extract the registry contract and chain from the environment
	registryChain, ok := env.Chains[req.RegistryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	contracts, err := kslib.GetContractSets(env.Logger, &kslib.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	_, err = internal.UpdateNodes(env.Logger, &internal.UpdateNodesRequest{
		Chain:        registryChain,
		Registry:     contracts.ContractSets[req.RegistryChainSel].CapabilitiesRegistry,
		P2pToUpdates: req.P2pToUpdates,
		UseMCMS:      req.UseMCMS,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}
