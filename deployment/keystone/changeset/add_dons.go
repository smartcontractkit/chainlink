package changeset

import (
	"errors"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ deployment.ChangeSet[*AddDonsRequest] = AddDons

type RegisterableDon struct {
	Name              string             // the name of the DON
	F                 uint8              // the number of faulty nodes tolerated by the DON
	P2PIDs            []p2pkey.PeerID    // this is the unique identifier for the don
	CapabilityConfigs []CapabilityConfig // if Config subfield is nil, a default config is used
}

func (d *RegisterableDon) Validate() error {
	if d.Name == "" {
		return errors.New("name is required")
	}
	if d.F == 0 {
		return errors.New("F is required")
	}
	if len(d.P2PIDs) == 0 {
		return errors.New("P2PIDs is required")
	}
	return nil
}

func (d *RegisterableDon) ToDONToRegister(registry *kcr.CapabilitiesRegistry) (donToRegister internal.DONToRegister, nodeIDToP2PID map[string][32]byte, donCapabilities []internal.RegisteredCapability, err error) {
	nodeIDToP2PID = make(map[string][32]byte)
	for _, p2pid := range d.P2PIDs {
		nodeIDToP2PID[p2pid.Raw()] = p2pid
	}
	donCapabilities = make([]internal.RegisteredCapability, 0, len(d.CapabilityConfigs))
	donToRegister = internal.DONToRegister{
		Name: d.Name,
		F:    d.F,
	}
	for _, cap := range d.CapabilityConfigs {
		var cfg pb.CapabilityConfig
		if cap.Config != nil {
			err := proto.Unmarshal(cap.Config, &cfg)
			if err != nil {
				return donToRegister, nil, nil, fmt.Errorf("failed to unmarshal capability config: %w", err)
			}
		}
		r, err := internal.NewRegisteredCapability(registry, &cap.Capability, &cfg)
		if err != nil {
			return donToRegister, nil, nil, fmt.Errorf("failed to create registered capability: %w", err)
		}
		donCapabilities = append(donCapabilities, *r)
	}
	nodes := make([]deployment.Node, len(d.P2PIDs))
	for i, p2pid := range d.P2PIDs {
		nodes[i] = deployment.Node{
			NodeID: p2pid.Raw(),
			Name:   p2pid.Raw(),
			PeerID: p2pid,
		}
	}
	donToRegister.Nodes = nodes

	return
}

type AddDonsRequest struct {
	RegistryChainSel uint64

	DONs []*RegisterableDon
	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig

	RegistryRef datastore.AddressRefKey
}

func (r *AddDonsRequest) Validate(e deployment.Environment) error {
	for _, don := range r.DONs {
		if err := don.Validate(); err != nil {
			return fmt.Errorf("invalid DON to register: %w", err)
		}
	}
	_, exists := chainsel.ChainBySelector(r.RegistryChainSel)
	if !exists {
		return fmt.Errorf("invalid registry chain selector %d: selector does not exist", r.RegistryChainSel)
	}

	_, exists = e.Chains[r.RegistryChainSel]
	if !exists {
		return fmt.Errorf("invalid registry chain selector %d: chain does not exist in environment", r.RegistryChainSel)
	}
	if err := shouldUseDatastore(e, r.RegistryRef); err != nil {
		return fmt.Errorf("invalid registry reference: %w", err)
	}
	return nil
}

func (r AddDonsRequest) UseMCMS() bool {
	return r.MCMSConfig != nil
}

func (r AddDonsRequest) convertInternal(registry *kcr.CapabilitiesRegistry) (donsToRegister []internal.DONToRegister, nodeIDToP2PID map[string][32]byte, donToCapabilities map[string][]internal.RegisteredCapability, err error) {
	donsToRegister = make([]internal.DONToRegister, len(r.DONs))
	nodeIDToP2PID = make(map[string][32]byte)
	donToCapabilities = make(map[string][]internal.RegisteredCapability)
	for i, don := range r.DONs {
		dr, nodes, dc, err := don.ToDONToRegister(registry)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to convert DON to register: %w", err)
		}
		donsToRegister[i] = dr
		donToCapabilities[don.Name] = dc
		for id, node := range nodes {
			nodeIDToP2PID[id] = node
		}
	}
	return
}

// AddDons adds a DON to the capabilities registry
// The Nodes and Capabilities *must be registered before calling this function*
// See [AddNodes] and [AddCapabilities] for more detail
// If the DON already exists, it will do nothing. The DON is identified by the P2P addresses of the nodes.
// If you need to update the DON, use [UpdateDon].
func AddDons(env deployment.Environment, req *AddDonsRequest) (deployment.ChangesetOutput, error) {
	if err := req.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	// extract the registry contract and chain from the environment
	registryChain, ok := env.Chains[req.RegistryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	capReg, err := loadCapabilityRegistry(registryChain, env, req.RegistryRef)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load capability registry: %w", err)
	}
	donsToRegister, nodeIDToP2PID, donToCapabilities, err := req.convertInternal(capReg.Contract)

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert DON to register: %w", err)
	}
	resp, err := internal.RegisterDons(env.Logger, internal.RegisterDonsRequest{
		RegistryChainSelector: req.RegistryChainSel,
		RegistryChain:         &registryChain,
		Registry:              capReg.Contract,
		NodeIDToP2PID:         nodeIDToP2PID,
		DonToCapabilities:     donToCapabilities,
		DonsToRegister:        donsToRegister,
		UseMCMS:               req.UseMCMS(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to add dons: %w", err)
	}

	out := deployment.ChangesetOutput{}
	if req.UseMCMS() {
		if resp.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}

		timelocksPerChain := map[uint64]string{
			req.RegistryChainSel: capReg.McmsContracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			req.RegistryChainSel: capReg.McmsContracts.ProposerMcm.Address().Hex(),
		}
		inspector, err := proposalutils.McmsInspectorForChain(env, req.RegistryChainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		inspectorPerChain := map[uint64]sdk.Inspector{
			req.RegistryChainSel: inspector,
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocksPerChain,
			proposerMCMSes,
			inspectorPerChain,
			[]types.BatchOperation{*resp.Ops},
			"proposal to add dons",
			proposalutils.TimelockConfig{MinDelay: req.MCMSConfig.MinDuration},
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}
	return out, nil
}
