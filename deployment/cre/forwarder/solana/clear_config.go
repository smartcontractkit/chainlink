package solana

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana/sequence/operation"
)

type ClearForwarderConfigRequest struct {
	DonID         uint32 // the DON id whose config should be closed on the forwarder
	ConfigVersion uint32 // the config version to close. Must match the version the forwarder was configured with

	MCMS *cldfproposalutils.TimelockConfig // if set, assumes current ownership is the timelock

	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains    map[uint64]struct{}
	Qualifier string
	Version   string
}

var _ cldf.ChangeSetV2[*ClearForwarderConfigRequest] = ClearForwarderConfigs{}

// ClearForwarderConfigs closes the oracles config account of a DON on Keystone Forwarder contracts.
// Closing refunds the rent to the owner and frees the (DON id, config version) pair, so a later
// ConfigureForwarders with the same pair initializes it again from scratch.
type ClearForwarderConfigs struct{}

func (cs ClearForwarderConfigs) VerifyPreconditions(env cldf.Environment, req *ClearForwarderConfigRequest) error {
	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}

	if req.DonID == 0 {
		return errors.New("DON ID must be non-zero")
	}
	if req.ConfigVersion == 0 {
		return errors.New("config version must be non-zero")
	}

	return verifyForwarderChains(env, req.Chains, version, req.Qualifier, req.MCMS)
}

func (cs ClearForwarderConfigs) Apply(env cldf.Environment, req *ClearForwarderConfigRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	mcmsBatches, err := clearForwarderConfigs(env, req)
	if err != nil {
		return out, fmt.Errorf("failed to clear forwarder config: %w", err)
	}

	if req.MCMS == nil {
		return out, nil
	}

	out.MCMSTimelockProposals, err = buildTimelockProposals(env, mcmsBatches, *req.MCMS,
		"proposal to clear the config of keystone forwarder contract")
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return out, nil
}

func clearForwarderConfigs(env cldf.Environment, req *ClearForwarderConfigRequest) (map[uint64]mcmsTypes.BatchOperation, error) {
	version := semver.MustParse(req.Version)

	batches := make(map[uint64]mcmsTypes.BatchOperation)
	for chain := range forwarderChains(env, req.Chains) {
		target, err := resolveForwarderConfigTarget(env, chain, version, req.Qualifier, req.DonID, req.ConfigVersion, req.MCMS)
		if err != nil {
			return nil, fmt.Errorf("chain selector %d: %w", chain.Selector, err)
		}

		deps := operation.Deps{
			Datastore: env.DataStore,
			Env:       env,
			Chain:     chain,
		}

		opOut, err := operations.ExecuteOperation(env.OperationsBundle, operation.ClearForwarderConfigOp, deps, operation.ClearForwarderConfigInput{
			ForwarderConfigTarget: target,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to clear forwarder config for chain selector %d: %w", chain.Selector, err)
		}

		batches[chain.Selector] = opOut.Output.Batch
	}

	return batches, nil
}
