package operation

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type PostOpsInput struct {
	MCMSConfig *proposalutils.TimelockConfig
	Proposals  []mcms.TimelockProposal
}

var (
	PostOpsAggregateProposals = operations.NewOperation(
		"postOpsToAggregateProposals",
		semver.MustParse("1.0.0"),
		"Post ops to aggregate proposals",
		func(b operations.Bundle, deps opsutil.ConfigureDependencies, input PostOpsInput) ([]mcms.TimelockProposal, error) {
			allProposals := input.Proposals
			evmState := deps.CurrentState.EVMMCMSStateByChain()
			solanaState := deps.CurrentState.SolanaMCMSStateByChain(deps.Env)
			proposal, err := proposalutils.AggregateProposals(
				deps.Env, evmState, solanaState, allProposals,
				"Aggregating all proposals", input.MCMSConfig)
			if err != nil {
				return nil, err
			}
			if proposal != nil {
				input.Proposals = []mcms.TimelockProposal{*proposal}
			}
			return input.Proposals, nil
		},
	)
)
