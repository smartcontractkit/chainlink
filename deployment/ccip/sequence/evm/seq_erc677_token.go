package v1_5_1

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

// GrantMintRoleSeqInp defines inputs for granting mint role across multiple chains
type GrantMintRoleSeqInp struct {
	// UpdatesByChain maps chain selector to the EVM call input for that chain
	UpdatesByChain map[uint64]opsutil.EVMCallInput[common.Address]
}

var (
	// GrantMintAndBurnRoleOnERC677Sequence grants mint & burn role on ERC677 token contracts across multiple EVM chains
	GrantMintAndBurnRoleOnERC677Sequence = operations.NewSequence(
		"GrantMintAndBurnRoleSequence",
		semver.MustParse("1.0.0"),
		"Grant mint & Burn role on ERC677 token contracts across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, inputs GrantMintRoleSeqInp) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(inputs.UpdatesByChain))
			for chainSel, input := range inputs.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.GrantMintAndBurnRolesERC677Op, chain, input)
				if err != nil {
					return nil, fmt.Errorf("failed to execute GrantMintAndBurnRolesERC677Op on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
