package v1_6

import (
	"fmt"

	"dario.cat/mergo"
	"github.com/Masterminds/semver/v3"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type UpdateLanesSequenceInput struct {
	FeeQuoterApplyDestChainConfigUpdatesSequenceInput
	FeeQuoterUpdatePricesSequenceInput
	OffRampApplySourceChainConfigUpdatesSequenceInput
	OnRampApplyDestChainConfigUpdatesSequenceInput
	RouterApplyRampUpdatesSequenceInput
}

var UpdateLanesSequence = operations.NewSequence(
	"UpdateLanesSequence",
	semver.MustParse("1.0.0"),
	"Updates lanes on CCIP 1.6.0",
	func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input UpdateLanesSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
		result := make(map[uint64][]opsutil.EVMCallOutput)

		result, err := runAndMergeSequence(b, chains, FeeQuoterApplyDestChainConfigUpdatesSequence, input.FeeQuoterApplyDestChainConfigUpdatesSequenceInput, result)
		if err != nil {
			return nil, err
		}
		b.Logger.Info("Destination configs updated on FeeQuoters")

		result, err = runAndMergeSequence(b, chains, FeeQuoterUpdatePricesSequence, input.FeeQuoterUpdatePricesSequenceInput, result)
		if err != nil {
			return nil, err
		}
		b.Logger.Info("Gas prices updated on FeeQuoters")

		result, err = runAndMergeSequence(b, chains, OffRampApplySourceChainConfigUpdatesSequence, input.OffRampApplySourceChainConfigUpdatesSequenceInput, result)
		if err != nil {
			return nil, err
		}
		b.Logger.Info("Destination configs updated on OnRamps")

		result, err = runAndMergeSequence(b, chains, OnRampApplyDestChainConfigUpdatesSequence, input.OnRampApplyDestChainConfigUpdatesSequenceInput, result)
		if err != nil {
			return nil, err
		}
		b.Logger.Info("Source configs updated on OffRamps")

		result, err = runAndMergeSequence(b, chains, RouterApplyRampUpdatesSequence, input.RouterApplyRampUpdatesSequenceInput, result)
		if err != nil {
			return nil, err
		}
		b.Logger.Info("Ramps updated on Routers")

		return result, nil
	},
)

func runAndMergeSequence[IN any](
	b operations.Bundle,
	chains map[uint64]cldf_evm.Chain,
	seq *operations.Sequence[IN, map[uint64][]opsutil.EVMCallOutput, map[uint64]cldf_evm.Chain],
	input IN,
	agg map[uint64][]opsutil.EVMCallOutput,
) (map[uint64][]opsutil.EVMCallOutput, error) {
	if agg == nil {
		agg = make(map[uint64][]opsutil.EVMCallOutput)
	}
	report, err := operations.ExecuteSequence(b, seq, chains, input)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", seq.ID(), err)
	}
	err = mergo.Merge(&agg, report.Output, mergo.WithAppendSlice)
	if err != nil {
		return nil, fmt.Errorf("failed to merge output of %s: %w", seq.ID(), err)
	}
	return agg, nil
}
