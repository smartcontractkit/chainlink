package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

// AddCapabilitiesRequest is a request to add capabilities
type AddCapabilitiesRequest struct {
	RegistryChainSel uint64

	Capabilities []kcr.CapabilitiesRegistryCapability
	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *MCMSConfig

	RegistryRef datastore.AddressRefKey
}

func (r *AddCapabilitiesRequest) Validate(env cldf.Environment) error {
	if r.RegistryChainSel == 0 {
		return errors.New("registry chain selector must be set")
	}
	if len(r.Capabilities) == 0 {
		return errors.New("capabilities must be set")
	}

	if err := shouldUseDatastore(env, r.RegistryRef); err != nil {
		return fmt.Errorf("failed to check registry ref: %w", err)
	}
	return nil
}

// if the environment has a non-empty datastore, the registry ref must be set
// prevents accidental usage of the old address book
func shouldUseDatastore(env cldf.Environment, ref datastore.AddressRefKey) error {
	if addrs, err := env.DataStore.Addresses().Fetch(); err == nil {
		if len(addrs) != 0 && ref == nil {
			return errors.New("This environment has been migrated to DataStore: address ref key must not be nil")
		}
	}
	return nil
}

type AddCapabilitiesRequestV2 = struct {
	AddCapabilitiesRequest
	RegistryRef datastore.AddressRefKey
}

var _ cldf.ChangeSet[*AddCapabilitiesRequest] = AddCapabilities

// AddCapabilities is a deployment.ChangeSet that adds capabilities to the capabilities registry
//
// It is idempotent. It deduplicates the input capabilities.
//
// When using MCMS, the output will contain a single proposal with a single batch containing all capabilities to be added.
// When not using MCMS, each capability will be added in a separate transaction.
func AddCapabilities(env cldf.Environment, req *AddCapabilitiesRequest) (cldf.ChangesetOutput, error) {
	err := req.Validate(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate request: %w", err)
	}
	registryChain, ok := env.BlockChains.EVMChains()[req.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}

	cr, err := loadCapabilityRegistry(registryChain, env, req.RegistryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load capability registry: %w", err)
	}
	useMCMS := req.MCMSConfig != nil
	ops, err := internal.AddCapabilities(env.Logger, cr.Contract, env.BlockChains.EVMChains()[req.RegistryChainSel], req.Capabilities, useMCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to add capabilities: %w", err)
	}
	out := cldf.ChangesetOutput{}
	if useMCMS {
		if ops == nil {
			return out, errors.New("expected MCMS operation to be non-nil")
		}
		if cr.McmsContracts == nil {
			return out, fmt.Errorf("expected capabiity registry contract %s to be owned by MCMS", cr.Contract.Address().String())
		}
		timelocksPerChain := map[uint64]string{
			registryChain.Selector: cr.McmsContracts.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			registryChain.Selector: cr.McmsContracts.ProposerMcm.Address().Hex(),
		}
		inspector, err := proposalutils.McmsInspectorForChain(env, req.RegistryChainSel)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		inspectorPerChain := map[uint64]mcmssdk.Inspector{
			req.RegistryChainSel: inspector,
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocksPerChain,
			proposerMCMSes,
			inspectorPerChain,
			[]mcmstypes.BatchOperation{*ops},
			"proposal to add capabilities",
			proposalutils.TimelockConfig{MinDelay: req.MCMSConfig.MinDuration},
		)
		if err != nil {
			return out, fmt.Errorf("failed to build proposal: %w", err)
		}
		out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	}
	return out, nil
}
