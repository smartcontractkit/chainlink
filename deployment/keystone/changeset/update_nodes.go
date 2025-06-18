package changeset

import (
	"errors"
	"fmt"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type MCMSConfig struct {
	MinDuration time.Duration
}

var _ cldf.ChangeSet[*UpdateNodesRequest] = UpdateNodes

type UpdateNodesRequest struct {
	RegistryChainSel uint64
	P2pToUpdates     map[p2pkey.PeerID]NodeUpdate

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig

	RegistryRef datastore.AddressRefKey
}

func (r *UpdateNodesRequest) Validate(e cldf.Environment) error {
	if r.P2pToUpdates == nil {
		return errors.New("P2pToUpdates must be non-nil")
	}

	_, exists := chainsel.ChainBySelector(r.RegistryChainSel)
	if !exists {
		return fmt.Errorf("invalid registry chain selector %d: selector does not exist", r.RegistryChainSel)
	}

	_, exists = e.BlockChains.EVMChains()[r.RegistryChainSel]
	if !exists {
		return fmt.Errorf("invalid registry chain selector %d: chain does not exist in environment", r.RegistryChainSel)
	}

	if err := shouldUseDatastore(e, r.RegistryRef); err != nil {
		return fmt.Errorf("invalid registry reference: %w", err)
	}
	return nil
}

func (r UpdateNodesRequest) UseMCMS() bool {
	return r.MCMSConfig != nil
}

type NodeUpdate = internal.NodeUpdate

// UpdateNodes updates a set of nodes.
// The nodes and capabilities in the request must already exist in the registry contract.
func UpdateNodes(env cldf.Environment, req *UpdateNodesRequest) (cldf.ChangesetOutput, error) {
	if err := req.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid request: %w", err)
	}
	// extract the registry contract and chain from the environment
	registryChain, ok := env.BlockChains.EVMChains()[req.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	capReg, err := loadCapabilityRegistry(registryChain, env, req.RegistryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load capability registry: %w", err)
	}

	resp, err := internal.UpdateNodes(env.Logger, &internal.UpdateNodesRequest{
		Chain:                registryChain,
		CapabilitiesRegistry: capReg.Contract,
		P2pToUpdates:         req.P2pToUpdates,
		UseMCMS:              req.UseMCMS(),
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}

	out := cldf.ChangesetOutput{}
	if req.UseMCMS() {
		if resp.Ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		if capReg.McmsContracts == nil {
			return out, fmt.Errorf("expected capabiity registry contract %s to be owned by MCMS", capReg.Contract.Address().String())
		}
		timelocksPerChain := map[uint64]string{
			req.RegistryChainSel: capReg.McmsContracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			req.RegistryChainSel: capReg.McmsContracts.ProposerMcm.Address().Hex(),
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
			[]types.BatchOperation{*resp.Ops},
			"proposal to set update nodes",
			proposalutils.TimelockConfig{MinDelay: req.MCMSConfig.MinDuration},
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}

	return out, nil
}
