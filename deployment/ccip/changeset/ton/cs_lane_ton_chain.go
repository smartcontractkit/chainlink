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

type AddLane struct {
}

var _ cldf.ChangeSetV2[AddLaneCfg] = AddLane{}

func (c AddLane) VerifyPreconditions(env cldf.Environment, config AddLaneCfg) error {
	// TODO: Implement precondition checks for adding or updating a lane on Ton chain
	return nil
}

func (cs AddLane) Apply(env cldf.Environment, config AddLaneCfg) (cldf.ChangesetOutput, error) {
	// TODO: Implement logic of adding or updating a lane on Ton chain
	return cldf.ChangesetOutput{}, nil
}
