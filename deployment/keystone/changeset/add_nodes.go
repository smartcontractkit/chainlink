package changeset

import (
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

// NOPIdentity is a node operator identity
//
// either by operator or registration ID must be non empty
// it is an error to have both set
type NOPIdentity struct {
	Operator       kcr.CapabilitiesRegistryNodeOperator
	RegistrationID uint32 // onchain registration ID; 1-indexed
}

func (i NOPIdentity) Validate() error {
	dflt := kcr.CapabilitiesRegistryNodeOperator{}
	if i.Operator == dflt && i.RegistrationID == 0 {
		return errors.New("NOPIdentity must have either Operator or RegistrationID set, both empty")
	}
	if i.Operator != dflt && i.RegistrationID != 0 {
		return fmt.Errorf("NOPIdentity must have either Operator (%v) or RegistrationID (%d) set, both set", i.Operator, i.RegistrationID)
	}
	return nil
}

// resolve returns the registration ID of the NOP
func (i NOPIdentity) resolve(registry *kcr.CapabilitiesRegistry) (uint32, error) {
	if i.RegistrationID != 0 {
		_, err := registry.GetNodeOperator(nil, i.RegistrationID)
		if err != nil {
			return 0, fmt.Errorf("failed to get node operator %d: %w", i.RegistrationID, err)
		}
		return i.RegistrationID, nil
	}
	nops, err := registry.GetNodeOperators(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get node operators: %w", err)
	}
	id := uint32(0)
	for _, nop := range nops {
		if nop.Name == i.Operator.Name && nop.Admin == i.Operator.Admin {
			id++ // 1-indexed; ordered
			break
		}
	}
	if id == 0 {
		return 0, fmt.Errorf("NOP %v not found in capabilities registry", i.Operator)
	}
	return id, nil
}

type CapabilityIdentity struct {
	Capability     kcr.CapabilitiesRegistryCapability
	RegistrationID [32]byte
}

func (c CapabilityIdentity) Validate() error {
	if c.Capability == (kcr.CapabilitiesRegistryCapability{}) && c.RegistrationID == [32]byte{} {
		return errors.New("CapabilityIdentity must have either Capability or RegistrationID set, both empty")
	}
	if c.Capability != (kcr.CapabilitiesRegistryCapability{}) && c.RegistrationID != [32]byte{} {
		return fmt.Errorf("CapabilityIdentity must have either Capability (%v) or RegistrationID (%x) set, both set", c.Capability, c.RegistrationID)
	}
	return nil
}

func (c CapabilityIdentity) resolve(registry *kcr.CapabilitiesRegistry) ([32]byte, error) {
	if c.RegistrationID != [32]byte{} {
		_, err := registry.GetCapability(nil, c.RegistrationID)
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to get capability %x: %w", c.RegistrationID, err)
		}
		return c.RegistrationID, nil
	}
	caps, err := registry.GetCapabilities(nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to get capabilities: %w", err)
	}
	wantID, err := registry.GetHashedCapabilityId(nil, c.Capability.LabelledName, c.Capability.Version)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to get capability ID for capability %v: %w", c, err)
	}
	for _, c := range caps {
		if c.HashedId == wantID {
			return c.HashedId, nil
		}
	}
	return [32]byte{}, fmt.Errorf("capability %v not found in capabilities registry", c)
}

type CapabilityIdentities []CapabilityIdentity

func (capabilities CapabilityIdentities) Validate() error {
	for i, c := range capabilities {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("invalid CapabilityIdentity at %d: %w", i, err)
		}
	}
	return nil
}

func (capabilities CapabilityIdentities) resolve(registry *kcr.CapabilitiesRegistry) ([][32]byte, error) {
	out := make([][32]byte, len(capabilities))
	for i, c := range capabilities {
		id, err := c.resolve(registry)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve CapabilityIdentity %d: %w", i, err)
		}
		out[i] = id
	}
	return out, nil
}

type CreateNodeRequest struct {
	NOPIdentity
	Signer               [32]byte // signer address of the NOP
	P2PID                [32]byte // p2p ID of the node
	EncryptionPublicKey  [32]byte // encryption public key of the node
	CapabilityIdentities          // the capabilities of the node; must all exist in the capabilities registry
}

func (r *CreateNodeRequest) Validate() error {
	if err := r.NOPIdentity.Validate(); err != nil {
		return fmt.Errorf("invalid NOPIdentity: %w", err)
	}
	if err := r.CapabilityIdentities.Validate(); err != nil {
		return fmt.Errorf("invalid CapabilityIdentities: %w", err)
	}
	if r.Signer == [32]byte{} {
		return errors.New("signer address is required")
	}
	if r.P2PID == [32]byte{} {
		return errors.New("p2p ID is required")
	}
	if r.EncryptionPublicKey == [32]byte{} {
		return errors.New("encryption public key is required")
	}
	return nil
}

func (r *CreateNodeRequest) Resolve(registry *kcr.CapabilitiesRegistry) (kcr.CapabilitiesRegistryNodeParams, error) {
	id, err := r.NOPIdentity.resolve(registry)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to resolve NOPIdentity: %w", err)
	}
	capIDs, err := r.CapabilityIdentities.resolve(registry)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to resolve CapabilityIdentities: %w", err)
	}
	return kcr.CapabilitiesRegistryNodeParams{
		NodeOperatorId:      id,
		P2pId:               r.P2PID,
		EncryptionPublicKey: r.EncryptionPublicKey,
		Signer:              r.Signer,
		HashedCapabilityIds: capIDs,
	}, nil
}

type AddNodesRequest struct {
	RegistryChainSel uint64

	CreateNodeRequests map[string]CreateNodeRequest
	// MCMS is the configuration for the Multi-Chain Manager Service
	// Required if the registry contract has be delegated to MCMS
	// If nil, the registry contract will be used directly
	MCMSConfig *MCMSConfig
}

func (r *AddNodesRequest) Validate() error {
	if len(r.CreateNodeRequests) == 0 {
		return errors.New("must provide create node requests")
	}
	for nodeName, cr := range r.CreateNodeRequests {
		if err := cr.Validate(); err != nil {
			return fmt.Errorf("invalid create node request for node %s: %w", nodeName, err)
		}
	}
	return nil
}

var _ deployment.ChangeSet[*AddNodesRequest] = AddNodes

func AddNodes(env deployment.Environment, req *AddNodesRequest) (deployment.ChangesetOutput, error) {
	err := req.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid request: %w", err)
	}

	contractSetResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	nodeParams := make(map[string]kcr.CapabilitiesRegistryNodeParams)
	for nodeName, cr := range req.CreateNodeRequests {
		params, err := cr.Resolve(contractSetResp.ContractSets[req.RegistryChainSel].CapabilitiesRegistry)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to resolve node params for node %s: %w", nodeName, err)
		}
		p2p := string(params.P2pId[:])
		if _, exists := nodeParams[p2p]; exists {
			return deployment.ChangesetOutput{}, fmt.Errorf("duplicate p2pid %s at node %s", p2p, nodeName)
		}
		nodeParams[p2p] = params
	}

	var (
		useMCMS                = req.MCMSConfig != nil
		registryChain          = env.Chains[req.RegistryChainSel]
		registry               = contractSetResp.ContractSets[req.RegistryChainSel].CapabilitiesRegistry
		registryChainContracts = contractSetResp.ContractSets[req.RegistryChainSel]
	)
	resp, err := internal.AddNodes(env.Logger, &internal.AddNodesRequest{
		RegistryChain:        registryChain,
		CapabilitiesRegistry: registry,
		NodeParams:           nodeParams,
		UseMCMS:              useMCMS,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to add nodes: %w", err)
	}
	// create mcms proposal if needed
	out := deployment.ChangesetOutput{}
	if useMCMS {
		if resp.Ops == nil || len(resp.Ops.Batch) == 0 {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		timelocksPerChain := map[uint64]gethcommon.Address{
			registryChain.Selector: registryChainContracts.Timelock.Address(),
		}
		proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
			registryChain.Selector: registryChainContracts.ProposerMcm,
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

		out.Proposals = []timelock.MCMSWithTimelockProposal{*proposal} //nolint:staticcheck //SA1019 ignoring deprecated field for compatibility; we don't have tools to generate the new field
	}
	return out, nil
}
