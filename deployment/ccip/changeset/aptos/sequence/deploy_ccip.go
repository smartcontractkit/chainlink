package sequence

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

const aptosTokenAddress = "0xa"

type DeployCCIPSeqInput struct {
	MCMSAddress      aptos.AccountAddress
	CCIPConfig       config.ChainContractParams
	LinkTokenAddress aptos.AccountAddress
}

type DeployCCIPSeqOutput struct {
	CCIPAddress    aptos.AccountAddress
	MCMSOperations []mcmstypes.BatchOperation
}

var DeployCCIPSequence = operations.NewSequence(
	"deploy-aptos-ccip-sequence",
	operation.Version1_0_0,
	"Deploy Aptos CCIP contracts and initialize them",
	deployCCIPSequence,
)

func deployCCIPSequence(b operations.Bundle, deps operation.AptosDeps, in DeployCCIPSeqInput) (DeployCCIPSeqOutput, error) {
	var mcmsOperations []mcmstypes.BatchOperation

	// Cleanup MCMS staging area if not clear
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, in.MCMSAddress)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// Generate batch operations to deploy CCIP package
	deployCCIPInput := operation.DeployCCIPInput{
		MCMSAddress: in.MCMSAddress,
		IsUpdate:    false,
	}
	deployCCIPReport, err := operations.ExecuteOperation(b, operation.DeployCCIPOp, deps, deployCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	ccipAddress := deployCCIPReport.Output.CCIPAddress
	// For CCIP deployment the txs cannot be batched - it'd exceed Aptos API limits
	// so they're converted to batch operations with single transactions in each batch
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployCCIPReport.Output.MCMSOperations)...)

	// Generate batch operations to deploy CCIP modules
	deployModulesInput := operation.DeployModulesInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
	}
	// OnRamp module
	deployOnRampReport, err := operations.ExecuteOperation(b, operation.DeployOnRampOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOnRampReport.Output)...)
	// OffRamp module
	deployOffRampReport, err := operations.ExecuteOperation(b, operation.DeployOffRampOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployOffRampReport.Output)...)
	// Router module
	deployRouterReport, err := operations.ExecuteOperation(b, operation.DeployRouterOp, deps, deployModulesInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployRouterReport.Output)...)

	var txs []mcmstypes.Transaction
	// Generate txs to Initialize CCIP
	initCCIPInput := operation.InitializeCCIPInput{
		MCMSAddress:      in.MCMSAddress,
		CCIPAddress:      ccipAddress,
		CCIPConfig:       in.CCIPConfig,
		LinkTokenAddress: in.LinkTokenAddress,
	}
	initCCIPReport, err := operations.ExecuteOperation(b, operation.InitializeCCIPOp, deps, initCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	txs = append(txs, initCCIPReport.Output...)

	// Apply Premium multiplier on fee tokens
	multiplierConfig, err := getMultiplierConfig(in.CCIPConfig.FeeQuoterParams.PremiumMultiplierWeiPerEthByFeeToken, in.LinkTokenAddress)
	if err != nil {
		return DeployCCIPSeqOutput{}, fmt.Errorf("failed to get multiplier config: %w", err)
	}
	apmInput := operation.ApplyPremiumMultiplierInput{
		CCIPAddress:             ccipAddress,
		MultiplierBySourceToken: multiplierConfig,
	}
	applyPMultReport, err := operations.ExecuteOperation(b, operation.ApplyPremiumMultiplierOp, deps, apmInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	txs = append(txs, applyPMultReport.Output...)

	//  Generate batch operation
	mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	})

	return DeployCCIPSeqOutput{
		CCIPAddress:    ccipAddress,
		MCMSOperations: mcmsOperations,
	}, nil
}

func getMultiplierConfig(multiplierBySymbol map[shared.TokenSymbol]uint64, linkTokenAddress aptos.AccountAddress) (map[string]uint64, error) {
	multiplierByToken := make(map[string]uint64)
	for symbol, multiplier := range multiplierBySymbol {
		switch symbol {
		case shared.APTSymbol:
			address := aptos.AccountAddress{}
			err := address.ParseStringRelaxed(aptosTokenAddress)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Aptos token address %s: %w", aptosTokenAddress, err)
			}
			multiplierByToken[address.StringLong()] = multiplier
		case shared.LinkSymbol:
			multiplierByToken[linkTokenAddress.StringLong()] = multiplier
		default:
			return nil, fmt.Errorf("unsupported fee token symbol %s for Aptos CCIP", symbol)
		}
	}
	return multiplierByToken, nil
}
