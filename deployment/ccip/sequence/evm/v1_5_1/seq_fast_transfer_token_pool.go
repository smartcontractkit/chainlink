package v1_5_1

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

// FastTransferTokenPoolUpdateDestChainConfigSequenceInput defines inputs for updating destination chain configs across multiple chains
type FastTransferTokenPoolUpdateDestChainConfigSequenceInput struct {
	// ContractType specifies which type of fast transfer token pool to update
	ContractType cldf.ContractType
	// UpdatesByChain maps chain selector to the EVM call input for that chain
	UpdatesByChain map[uint64]opsutil.EVMCallInput[ccipops.UpdateDestChainConfigInput]
}

// FastTransferTokenPoolUpdateFillerAllowlistSequenceInput defines inputs for updating filler allowlists across multiple chains
type FastTransferTokenPoolUpdateFillerAllowlistSequenceInput struct {
	// ContractType specifies which type of fast transfer token pool to update
	ContractType cldf.ContractType
	// UpdatesByChain maps chain selector to the EVM call input for that chain
	UpdatesByChain map[uint64]opsutil.EVMCallInput[ccipops.UpdateFillerAllowlistInput]
}

var (
	// FastTransferTokenPoolUpdateDestChainConfigSequence updates destination chain configurations
	// on fast transfer token pool contracts across multiple EVM chains
	FastTransferTokenPoolUpdateDestChainConfigSequence = operations.NewSequence(
		"FastTransferTokenPoolUpdateDestChainConfigSequence",
		semver.MustParse("1.0.0"),
		"Update destination chain configurations on fast transfer token pool contracts across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FastTransferTokenPoolUpdateDestChainConfigSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))

			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}

				// Select the appropriate operation based on contract type
				var operation *operations.Operation[opsutil.EVMCallInput[ccipops.UpdateDestChainConfigInput], opsutil.EVMCallOutput, cldf_evm.Chain]
				switch input.ContractType {
				case shared.BurnMintFastTransferTokenPool:
					operation = ccipops.BurnMintFastTransferTokenPoolUpdateDestChainConfigOp
				case shared.BurnMintWithExternalMinterFastTransferTokenPool:
					operation = ccipops.BurnMintWithExternalMinterFastTransferTokenPoolUpdateDestChainConfigOp
				default:
					return nil, fmt.Errorf("unsupported contract type for fast transfer token pool: %s", input.ContractType)
				}

				report, err := operations.ExecuteOperation(b, operation, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute fast transfer token pool update dest chain config op on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})

	// FastTransferTokenPoolUpdateFillerAllowlistSequence updates filler allowlists
	// on fast transfer token pool contracts across multiple EVM chains
	FastTransferTokenPoolUpdateFillerAllowlistSequence = operations.NewSequence(
		"FastTransferTokenPoolUpdateFillerAllowlistSequence",
		semver.MustParse("1.0.0"),
		"Update filler allowlists on fast transfer token pool contracts across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FastTransferTokenPoolUpdateFillerAllowlistSequenceInput) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))

			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}

				// Select the appropriate operation based on contract type
				var operation *operations.Operation[opsutil.EVMCallInput[ccipops.UpdateFillerAllowlistInput], opsutil.EVMCallOutput, cldf_evm.Chain]
				switch input.ContractType {
				case shared.BurnMintFastTransferTokenPool:
					operation = ccipops.BurnMintFastTransferTokenPoolUpdateFillerAllowlistOp
				case shared.BurnMintWithExternalMinterFastTransferTokenPool:
					operation = ccipops.BurnMintWithExternalMinterFastTransferTokenPoolUpdateFillerAllowlistOp
				default:
					return nil, fmt.Errorf("unsupported contract type for fast transfer token pool: %s", input.ContractType)
				}

				report, err := operations.ExecuteOperation(b, operation, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute fast transfer token pool update filler allowlist op on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
