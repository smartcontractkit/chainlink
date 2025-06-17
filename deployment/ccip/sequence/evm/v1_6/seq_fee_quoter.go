package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type FeeQuoterApplyDestChainConfigUpdatesSequenceInput struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[[]fee_quoter.FeeQuoterDestChainConfigArgs]
}

type FeeQuoterUpdatePricesSequenceInput struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[fee_quoter.InternalPriceUpdates]
}

type FeeQuoterUpdateTokenTransferConfig struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[ccipops.ApplyTokenTransferFeeConfigUpdatesConfigPerChain]
}

type FeeQuoterUpdateFeeTokensConfig struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[ccipops.ApplyFeeTokensUpdatesInput]
}

type FeeQuoterUpdatePremiumMultiplierWeiPerEthConfig struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs]
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

	FeeQUpdateTransferTokenFeeCfgSeq = operations.NewSequence(
		"FeeQuoterUpdateTransferTokenFeeConfigSequence",
		semver.MustParse("1.0.0"),
		"Update token and gas prices on FeeQuoter 1.6.0 contracts on multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FeeQuoterUpdateTokenTransferConfig) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, update := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.FeeQuoterApplyTokenTransferFeeCfgOp, chain, update)
				if err != nil {
					return nil, fmt.Errorf("failed to execute FeeQuoterApplyTokenTransferFeeCfgOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})

	FeeQuoterApplyFeeTokensUpdatesSeq = operations.NewSequence(
		"FeeQuoterApplyFeeTokensUpdatesSeq",
		semver.MustParse("1.0.0"),
		"Add or Remove supported tokens on FeeQuoter 1.6.0 contracts on multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FeeQuoterUpdateFeeTokensConfig) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, input := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.FeeQuoterApplyFeeTokensUpdatesOp, chain, input)
				if err != nil {
					return nil, fmt.Errorf("failed to execute FeeQuoterApplyFeeTokensUpdatesOp on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})

	FeeQApplyPremiumMultiplierWeiPerEthUpdatesSeq = operations.NewSequence(
		"FeeQApplyPremiumMultiplierWeiPerEthUpdatesSeq",
		semver.MustParse("1.0.0"),
		"Applies premiumMultiplierWeiPerEth for tokens in FeeQuoter 1.6.0 contract on multiple EVM chains",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FeeQuoterUpdatePremiumMultiplierWeiPerEthConfig) (map[uint64][]opsutil.EVMCallOutput, error) {
			opOutputs := make(map[uint64][]opsutil.EVMCallOutput, len(input.UpdatesByChain))
			for chainSel, input := range input.UpdatesByChain {
				chain, ok := chains[chainSel]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				report, err := operations.ExecuteOperation(b, ccipops.FeeQApplyPremiumMultiplierWeiPerEthUpdateOp, chain, input)
				if err != nil {
					return nil, fmt.Errorf("failed to execute ApplyPremiumMultiplierWeiPerEthUpdates on %s: %w", chain, err)
				}
				opOutputs[chainSel] = []opsutil.EVMCallOutput{report.Output}
			}
			return opOutputs, nil
		})
)
