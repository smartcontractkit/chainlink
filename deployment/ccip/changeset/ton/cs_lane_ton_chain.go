package ton

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type AddLaneCfg struct {
	FromChainSelector uint64
	ToChainSelector   uint64
	FromFamily        string
	ToFamily          string
}

type AddLane struct{}

var _ cldf.ChangeSetV2[AddLaneCfg] = AddLane{}

func (cs AddLane) VerifyPreconditions(_ cldf.Environment, _ AddLaneCfg) error {
	// TODO: Implement precondition checks for adding or updating a lane on Ton chain
	return nil
}

func (cs AddLane) Apply(_ cldf.Environment, _ AddLaneCfg) (cldf.ChangesetOutput, error) {
	// TODO: Implement logic of adding or updating a lane on Ton chain
	return cldf.ChangesetOutput{}, nil
}
