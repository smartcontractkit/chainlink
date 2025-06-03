package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type RouterApplyRampUpdatesSequenceInput struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[ccipops.RouterApplyRampUpdatesOpInput]
}

var (
	RouterApplyRampUpdatesSequence = operations.NewSequence(
		"RouterApplyRampUpdatesSequence",
		semver.MustParse("1.0.0"),
		"Updates OnRamps and OffRamps on Router contracts across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input RouterApplyRampUpdatesSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.RouterApplyRampUpdatesOp, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute RouterApplyRampUpdatesOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
