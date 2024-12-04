package changeset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
)

type ChangesetApplication struct {
	Changeset deployment.ChangeSet[any]
	Config    any
}

func WrapChangeSet[C any](fn deployment.ChangeSet[C]) func(e deployment.Environment, config any) (deployment.ChangesetOutput, error) {
	return func(e deployment.Environment, config any) (deployment.ChangesetOutput, error) {
		var zeroC C
		if config != nil {
			c, ok := config.(C)
			if !ok {
				return deployment.ChangesetOutput{}, fmt.Errorf("invalid config type, expected %T", c)
			}
			return fn(e, config.(C))
		}

		return fn(e, zeroC)
	}
}

// ApplyChangesets applies the changeset applications to the environment and returns the updated environment.
func ApplyChangesets(t *testing.T, e deployment.Environment, timelocksPerChain map[uint64]*gethwrappers.RBACTimelock, changesetApplications []ChangesetApplication) (deployment.Environment, error) {
	currentEnv := e
	for i, csa := range changesetApplications {
		out, err := csa.Changeset(currentEnv, csa.Config)
		if err != nil {
			return e, fmt.Errorf("failed to apply changeset at index %d: %w", i, err)
		}
		var addresses deployment.AddressBook
		if out.AddressBook != nil {
			addresses = out.AddressBook
			err := addresses.Merge(currentEnv.ExistingAddresses)
			if err != nil {
				return e, fmt.Errorf("failed to merge address book: %w", err)
			}
		} else {
			addresses = currentEnv.ExistingAddresses
		}
		if out.JobSpecs != nil {
			ctx := testcontext.Get(t)
			for nodeID, jobs := range out.JobSpecs {
				for _, job := range jobs {
					// Note these auto-accept
					_, err := currentEnv.Offchain.ProposeJob(ctx,
						&jobv1.ProposeJobRequest{
							NodeId: nodeID,
							Spec:   job,
						})
					if err != nil {
						return e, fmt.Errorf("failed to propose job: %w", err)
					}
				}
			}
		}
		if out.Proposals != nil {
			for _, prop := range out.Proposals {
				chains := mapset.NewSet[uint64]()
				for _, op := range prop.Transactions {
					chains.Add(uint64(op.ChainIdentifier))
				}

				signed := SignProposal(t, e, &prop)
				for _, sel := range chains.ToSlice() {
					timelock, ok := timelocksPerChain[sel]
					if !ok || timelock == nil {
						return deployment.Environment{}, fmt.Errorf("timelock not found for chain %d", sel)
					}
					ExecuteProposal(t, e, signed, timelock, sel)
				}
			}
		}
		currentEnv = deployment.Environment{
			Name:              e.Name,
			Logger:            e.Logger,
			ExistingAddresses: addresses,
			Chains:            e.Chains,
			NodeIDs:           e.NodeIDs,
			Offchain:          e.Offchain,
		}
	}
	return currentEnv, nil
}
