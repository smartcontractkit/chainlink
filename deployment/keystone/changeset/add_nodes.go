package changeset

import (
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-deployments/offchain"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

type NOPWithID struct {
	Admin   gethcommon.Address
	Name    string
	ID      uint32
	Deleted bool
}

type AddNodesRequest struct {
	RegistryChainSel uint64

	Nops []NOPWithID
	Dons []offchain.DONConfig

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig
}

var _ deployment.ChangeSet[*AddNodesRequest] = AddNodes

func AddNodes(env deployment.Environment, req *AddNodesRequest) (deployment.ChangesetOutput, error) {
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

	req2, err := buildRegisterNodesRequest(env, req.RegistryChainSel, req.Nops, req.Dons)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	req2.UseMCMS = req.MCMSConfig != nil

	resp, err := RegisterNodes(env.Logger, req2)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	env.Logger.Debugw("registered nodes", "resp", resp)
	out := deployment.ChangesetOutput{}
	if req2.UseMCMS {
		if resp.Ops == nil {
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
			[]timelock.BatchChainOperation{*resp.Ops},
			"proposal to add nodes",
			req.MCMSConfig.MinDuration,
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.Proposals = []timelock.MCMSWithTimelockProposal{*proposal}
	}
	return out, nil
}

func buildRegisterNodesRequest(e deployment.Environment, registryChainSel uint64, nops []NOPWithID, dons []offchain.DONConfig) (*RegisterNodesRequest, error) {
	var nopsA []*kcr.CapabilitiesRegistryNodeOperatorAdded
	nopNameToAddress := map[string]gethcommon.Address{}
	for _, nop := range nops {
		nopNameToAddress[nop.Name] = nop.Admin
		nopsA = append(nopsA, &kcr.CapabilitiesRegistryNodeOperatorAdded{
			Admin:          nop.Admin,
			Name:           nop.Name,
			NodeOperatorId: nop.ID,
		})
	}

	var allCSAKeys []string
	for _, don := range dons {
		for _, node := range don.Nodes {
			allCSAKeys = append(allCSAKeys, node.CSAKey)
		}
	}

	nodeInfos, err := deployment.NodeInfoByCSA(e, allCSAKeys)
	if err != nil {
		return nil, err
	}
	csaToNodeInfo := map[string]deployment.Node{}
	for _, nodeInfo := range nodeInfos {
		csaToNodeInfo[nodeInfo.CSAKey] = nodeInfo
	}

	nopToNodeIDs := map[kcr.CapabilitiesRegistryNodeOperator][]string{}
	for _, don := range dons {
		for _, node := range don.Nodes {
			nop := kcr.CapabilitiesRegistryNodeOperator{
				Name:  node.NOP,
				Admin: nopNameToAddress[node.NOP],
			}
			nopToNodeIDs[nop] = append(nopToNodeIDs[nop], csaToNodeInfo[node.CSAKey].NodeID)
		}
	}

	donToNodes := map[string][]deployment.Node{}
	for _, don := range dons {
		donToNodes[don.Name] = []deployment.Node{}
		for _, node := range don.Nodes {
			donToNodes[don.Name] = append(donToNodes[don.Name], csaToNodeInfo[node.CSAKey])
		}
	}

	donToCaps := map[string][]RegisteredCapability{}
	for _, don := range dons {
		donToCaps[don.Name] = []RegisteredCapability{}
		for i := range don.Capabilities {
			donCap := &don.Capabilities[i]
			rCap, capErr := internal.FromCapabilitiesRegistryCapability(&donCap.Capability, &donCap.Config.CapabilityConfig, e, registryChainSel)
			if capErr != nil {
				return nil, capErr
			}
			donToCaps[don.Name] = append(donToCaps[don.Name], *rCap)
		}
	}

	return &RegisterNodesRequest{
		Env:                   &e,
		RegistryChainSelector: registryChainSel,
		NopToNodeIDs:          nopToNodeIDs,
		DonToNodes:            donToNodes,
		DonToCapabilities:     donToCaps,
		Nops:                  nopsA,
	}, nil
}
