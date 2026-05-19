package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	opsevm "github.com/smartcontractkit/cld-changesets/pkg/family/evm/operations"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_2"
)

type RouterApplyRampUpdatesSequenceInput struct {
	UpdatesByChain map[uint64]opsevm.EVMCallInput[ccipops.RouterApplyRampUpdatesOpInput]
}

type RouterUpdateWrappedNativeSequenceInput struct {
	UpdatesByChain map[uint64]opsevm.EVMCallInput[common.Address]
}

var (
	RouterApplyRampUpdatesSequence = operations.NewSequence(
		"RouterApplyRampUpdatesSequence",
		semver.MustParse("1.0.0"),
		"Updates OnRamps and OffRamps on Router contracts across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input RouterApplyRampUpdatesSequenceInput) (map[uint64][]opsevm.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsevm.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.RouterApplyRampUpdatesOp, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute RouterApplyRampUpdatesOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsevm.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})

	RouterUpdateWrappedNativeSequence = operations.NewSequence(
		"RouterUpdateWrappedNativeSequence",
		semver.MustParse("1.0.0"),
		"Updates Wrapped Native token on Router contracts across multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input RouterUpdateWrappedNativeSequenceInput) (map[uint64][]opsevm.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsevm.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, newWrappedAddress := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.UpdateWrappedNativeAddressOnRouterOp, chain, newWrappedAddress)
				if err != nil {
					return nil, fmt.Errorf("failed to execute UpdateWrappedNativeAddressOnRouterOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsevm.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
