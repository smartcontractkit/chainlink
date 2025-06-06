package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type FeeQuoterApplyDestChainConfigUpdatesSequenceInput struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[[]fee_quoter.FeeQuoterDestChainConfigArgs]
}

type FeeQuoterUpdatePricesSequenceInput struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[fee_quoter.InternalPriceUpdates]
}

var (
	FeeQuoterApplyDestChainConfigUpdatesSequence = operations.NewSequence(
		"FeeQuoterApplyDestChainConfigUpdatesSequence",
		semver.MustParse("1.0.0"),
		"Apply updates to destination chain configs on the FeeQuoter 1.6.0 contract across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FeeQuoterApplyDestChainConfigUpdatesSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.FeeQuoterApplyDestChainConfigUpdatesOp, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute FeeQuoterApplyDestChainConfigUpdatesOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})

	FeeQuoterUpdatePricesSequence = operations.NewSequence(
		"FeeQuoterUpdatePricesSequence",
		semver.MustParse("1.0.0"),
		"Update token and gas prices on FeeQuoter 1.6.0 contracts on multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FeeQuoterUpdatePricesSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.FeeQuoterUpdatePricesOp, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute FeeQuoterUpdatePricesOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
