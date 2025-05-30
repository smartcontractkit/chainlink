package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
)

type FeeQuoterApplyDestChainConfigUpdatesSequenceInput struct {
	UpdatesByChain map[uint64]ccipops.FeeQuoterApplyDestChainConfigUpdatesOpInput
}

var (
	FeeQuoterApplyDestChainConfigUpdatesSequence = operations.NewSequence(
		"FeeQuoterApplyDestChainConfigUpdatesSequence",
		semver.MustParse("1.0.0"),
		"Apply updates to destination chain configs on the FeeQuoter 1.6.0 contract across multiple EVM chains",
		func(b operations.Bundle, deps map[uint64]ccipops.EVMCallDeps[*fee_quoter.FeeQuoter], input FeeQuoterApplyDestChainConfigUpdatesSequenceInput) (map[uint64]ccipops.EVMCallOutput, error) {
			opOutputs := make(map[uint64]ccipops.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, update := range input.UpdatesByChain {
				chainDeps, ok := deps[chainSel]
				if !ok {
					return map[uint64]ccipops.EVMCallOutput{}, fmt.Errorf("no dependencies defined for chain with selector %d", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.FeeQuoterApplyDestChainConfigUpdatesOp, chainDeps, update)
				if err != nil {
					return map[uint64]ccipops.EVMCallOutput{}, fmt.Errorf("failed to execute FeeQuoterApplyDestChainConfigUpdatesOp on chain with selector %d: %w", chainSel, err)
				}
				opOutputs[chainSel] = report.Output
			}
			return opOutputs, nil
		})
)
