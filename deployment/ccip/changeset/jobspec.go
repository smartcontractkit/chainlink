package changeset

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[any] = CCIPCapabilityJobspec

// CCIPCapabilityJobspec returns the job specs for the CCIP capability.
// The caller needs to propose these job specs to the offchain system.
func CCIPCapabilityJobspec(env deployment.Environment, _ any) (deployment.ChangesetOutput, error) {
	js, err := NewCCIPJobSpecs(env.NodeIDs, env.Offchain)
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(err, "failed to create job specs")
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		JobSpecs:    js,
	}, nil
}
