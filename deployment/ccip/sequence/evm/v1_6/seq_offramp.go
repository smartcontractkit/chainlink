package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type OffRampApplySourceChainConfigUpdatesSequenceInput struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[[]offramp.OffRampSourceChainConfigArgs]
}

var (
	OffRampApplySourceChainConfigUpdatesSequence = operations.NewSequence(
		"OffRampApplySourceChainConfigUpdatesSequence",
		semver.MustParse("1.0.0"),
		"Applies updates to source chain configurations stored on OffRamp contracts on multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input OffRampApplySourceChainConfigUpdatesSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.OffRampApplySourceChainConfigUpdatesOp, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute OffRampApplySourceChainConfigUpdatesOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
