package changeset

import (
	"errors"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[*RemoveDONsRequest] = RemoveDONs

type RemoveDONsRequest struct {
	RegistryChainSel uint64
	DONs             []uint32

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig

	RegistryRef datastore.AddressRefKey
}

func (r *RemoveDONsRequest) Validate(e deployment.Environment) error {
	if len(r.DONs) == 0 {
		return errors.New("dons is required")
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

func (r RemoveDONsRequest) UseMCMS() bool {
	return r.MCMSConfig != nil
}

// RemoveDONs removes a DON from the capabilities registry
func RemoveDONs(env deployment.Environment, req *RemoveDONsRequest) (deployment.ChangesetOutput, error) {
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

	resp, err := internal.RemoveDONs(env.Logger, &internal.RemoveDONsRequest{
		Chain:                registryChain,
		CapabilitiesRegistry: capReg.Contract,
		DONs:                 req.DONs,
		UseMCMS:              req.UseMCMS(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to remove don: %w", err)
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
			"proposal to remove dons",
			proposalutils.TimelockConfig{MinDelay: req.MCMSConfig.MinDuration},
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}

	return out, nil
}
