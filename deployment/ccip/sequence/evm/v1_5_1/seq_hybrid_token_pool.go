package v1_5_1

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

// HybridTokenPoolUpdateGroupsSequenceInput defines inputs for updating groups across multiple chains
type HybridTokenPoolUpdateGroupsSequenceInput struct {
	// ContractType specifies which type of hybrid token pool to update
	ContractType cldf.ContractType
	// UpdatesByChain maps chain selector to the EVM call input for that chain
	UpdatesByChain map[uint64]opsutil.EVMCallInput[ccipops.UpdateGroupsInput]
}

var (
	// HybridTokenPoolUpdateGroupsSequence updates groups on hybrid token pool contracts across multiple EVM chains
	HybridTokenPoolUpdateGroupsSequence = operations.NewSequence(
		"HybridTokenPoolUpdateGroupsSequence",
		semver.MustParse("1.0.0"),
		"Update groups on hybrid token pool contracts across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input HybridTokenPoolUpdateGroupsSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))

			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}

				// Select the appropriate operation based on contract type
				var operation *operations.Operation[opsutil.EVMCallInput[ccipops.UpdateGroupsInput], opsutil.EVMCallOutput, cldf_evm.Chain]
				switch input.ContractType {
				case shared.HybridWithExternalMinterFastTransferTokenPool:
					operation = ccipops.HybridWithExternalMinterTokenPoolUpdateGroupsOp
				default:
					return nil, fmt.Errorf("unsupported contract type for hybrid token pool group updates: %s", input.ContractType)
				}

				report, err := operations.ExecuteOperation(b, operation, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute hybrid token pool update groups op on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
