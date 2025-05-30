package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[*AppendNodeCapabilitiesRequest] = AppendNodeCapabilities

// AppendNodeCapabilitiesRequest is a request to add capabilities to the existing capabilities of nodes in the registry
type AppendNodeCapabilitiesRequest = MutateNodeCapabilitiesRequest

// AppendNodeCapabilities adds any new capabilities to the registry, merges the new capabilities with the existing capabilities
// of the node, and updates the nodes in the registry host the union of the new and existing capabilities.
func AppendNodeCapabilities(env cldf.Environment, req *AppendNodeCapabilitiesRequest) (cldf.ChangesetOutput, error) {
	c, capReg, err := req.convert(env, req.RegistryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	r, err := internal.AppendNodeCapabilitiesImpl(env.Logger, c)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	out := cldf.ChangesetOutput{}
	if req.UseMCMS() {
		if r.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		if capReg.McmsContracts == nil {
			return out, fmt.Errorf("expected capabiity registry contract %s to be owned by MCMS", capReg.Contract.Address().String())
		}
		timelocksPerChain := map[uint64]string{
			c.Chain.Selector: capReg.McmsContracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			c.Chain.Selector: capReg.McmsContracts.ProposerMcm.Address().Hex(),
		}
		inspector, err := proposalutils.McmsInspectorForChain(env, req.RegistryChainSel)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		inspectorPerChain := map[uint64]sdk.Inspector{
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

func (req *AppendNodeCapabilitiesRequest) convert(e cldf.Environment, ref datastore.AddressRefKey) (*internal.AppendNodeCapabilitiesRequest, *OwnedContract[*kcr.CapabilitiesRegistry], error) {
	if err := req.Validate(e); err != nil {
		return nil, nil, fmt.Errorf("failed to validate UpdateNodeCapabilitiesRequest: %w", err)
	}
	registryChain := e.BlockChains.EVMChains()[req.RegistryChainSel] // exists because of the validation above
	cr, err := loadCapabilityRegistry(registryChain, e, ref)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load capability registry: %w", err)
	}
	return &internal.AppendNodeCapabilitiesRequest{
		Chain:                registryChain,
		CapabilitiesRegistry: cr.Contract,
		P2pToCapabilities:    req.P2pToCapabilities,
		UseMCMS:              req.UseMCMS(),
	}, cr, nil
}
