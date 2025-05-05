package changeset

import (
	"errors"
	"fmt"
	"strconv"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ deployment.ChangeSet[*MutateNodeCapabilitiesRequest] = UpdateNodeCapabilities

type P2PSignerEnc = internal.P2PSignerEnc

func NewP2PSignerEnc(n *deployment.Node, registryChainSel uint64) (*P2PSignerEnc, error) {
	// TODO: deduplicate everywhere
	registryChainID, err := chainsel.ChainIdFromSelector(registryChainSel)
	if err != nil {
		return nil, err
	}
	registryChainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.Itoa(int(registryChainID)), chainsel.FamilyEVM)
	if err != nil {
		return nil, err
	}
	evmCC, exists := n.SelToOCRConfig[registryChainDetails]
	if !exists {
		return nil, fmt.Errorf("NewP2PSignerEnc: registryChainSel not found on node: %v", registryChainSel)
	}
	var signer [32]byte
	copy(signer[:], evmCC.OnchainPublicKey)
	var csakey [32]byte
	copy(csakey[:], evmCC.ConfigEncryptionPublicKey[:])

	return &P2PSignerEnc{
		Signer:              signer,
		P2PKey:              n.PeerID,
		EncryptionPublicKey: csakey,
	}, nil
}

// UpdateNodeCapabilitiesRequest is a request to set the capabilities of nodes in the registry
type UpdateNodeCapabilitiesRequest = MutateNodeCapabilitiesRequest

// MutateNodeCapabilitiesRequest is a request to change the capabilities of nodes in the registry
type MutateNodeCapabilitiesRequest struct {
	RegistryChainSel  uint64
	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig

	// RegistryRef is the reference to the registry contract in the datastore.
	RegistryRef datastore.AddressRefKey
}

func (req *MutateNodeCapabilitiesRequest) Validate(e deployment.Environment) error {
	if len(req.P2pToCapabilities) == 0 {
		return errors.New("p2pToCapabilities is empty")
	}
	_, exists := chainsel.ChainBySelector(req.RegistryChainSel)
	if !exists {
		return fmt.Errorf("invalid registry chain selector %d: selector does not exist", req.RegistryChainSel)
	}

	_, exists = e.Chains[req.RegistryChainSel]
	if !exists {
		return fmt.Errorf("invalid registry chain selector %d: chain does not exist in environment", req.RegistryChainSel)
	}
	if err := shouldUseDatastore(e, req.RegistryRef); err != nil {
		return fmt.Errorf("invalid registry reference: %w", err)
	}
	return nil
}

func (req *MutateNodeCapabilitiesRequest) UseMCMS() bool {
	return req.MCMSConfig != nil
}

func (req *MutateNodeCapabilitiesRequest) updateNodeCapabilitiesImplRequest(e deployment.Environment) (*internal.UpdateNodeCapabilitiesImplRequest, *OwnedContract[*kcr.CapabilitiesRegistry], error) {
	if err := req.Validate(e); err != nil {
		return nil, nil, fmt.Errorf("failed to validate UpdateNodeCapabilitiesRequest: %w", err)
	}
	registryChain := e.Chains[req.RegistryChainSel] // exists because of the validation above
	capReg, err := loadCapabilityRegistry(registryChain, e, req.RegistryRef)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load capability registry: %w", err)
	}

	return &internal.UpdateNodeCapabilitiesImplRequest{
		Chain:                registryChain,
		CapabilitiesRegistry: capReg.Contract,
		P2pToCapabilities:    req.P2pToCapabilities,
		UseMCMS:              req.UseMCMS(),
	}, capReg, nil
}

// UpdateNodeCapabilities updates the capabilities of nodes in the registry
func UpdateNodeCapabilities(env deployment.Environment, req *UpdateNodeCapabilitiesRequest) (deployment.ChangesetOutput, error) {
	c, capReg, err := req.updateNodeCapabilitiesImplRequest(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert request: %w", err)
	}

	r, err := internal.UpdateNodeCapabilitiesImpl(env.Logger, c)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update nodes: %w", err)
	}

	out := deployment.ChangesetOutput{}
	if req.UseMCMS() {
		if r.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		timelocksPerChain := map[uint64]string{
			c.Chain.Selector: capReg.McmsContracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			c.Chain.Selector: capReg.McmsContracts.ProposerMcm.Address().Hex(),
		}
		inspector, err := proposalutils.McmsInspectorForChain(env, req.RegistryChainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		inspectorPerChain := map[uint64]mcmssdk.Inspector{
			req.RegistryChainSel: inspector,
		}
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocksPerChain,
			proposerMCMSes,
			inspectorPerChain,
			[]types.BatchOperation{*r.Ops},
			"proposal to set update node capabilities",
			proposalutils.TimelockConfig{MinDelay: req.MCMSConfig.MinDuration},
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}
	return out, nil
}
