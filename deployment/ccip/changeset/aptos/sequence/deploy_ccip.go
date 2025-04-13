package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Deploy CCIP Sequence
type DeployCCIPSeqInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPConfig  config.ChainContractParams
}

type DeployCCIPSeqOutput struct {
	CCIPAddress    aptos.AccountAddress
	MCMSOperations []mcmstypes.Operation
}

var DeployCCIPSequence = operations.NewSequence(
	"deploy-aptos-ccip-sequence",
	operation.Version1_0_0,
	"Deploy Aptos CCIP contracts and initialize them",
	deployCCIPSequence,
)

func deployCCIPSequence(b operations.Bundle, deps operation.AptosDeps, in DeployCCIPSeqInput) (DeployCCIPSeqOutput, error) {
	var mcmsOperations []mcmstypes.Operation

	// Cleanup staging area
	cleanupInput := operation.CleanupStagingAreaInput{
		MCMSAddress: in.MCMSAddress,
	}
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, cleanupInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	if cleanupReport.Output != nil {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output...)
	}

	// Generate proposal to deploy CCIP package
	deployCCIPInput := operation.DeployCCIPInput{
		MCMSAddress: in.MCMSAddress,
	}
	deployCCIPReport, err := operations.ExecuteOperation(b, operation.GenerateDeployCCIPProposalOp, deps, deployCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	ccipAddress := deployCCIPReport.Output.CCIPAddress
	mcmsOperations = append(mcmsOperations, deployCCIPReport.Output.MCMSOperations...)

	// Generate proposal to deploy Router module
	deployRouterInput := operation.DeployRouterInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
	}
	deployRouterReport, err := operations.ExecuteOperation(b, operation.GenerateDeployRouterProposalOp, deps, deployRouterInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, deployRouterReport.Output...)

	// Generate proposal to initialize CCIP
	initCCIPInput := operation.InitializeCCIPInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
		CCIPConfig:  in.CCIPConfig,
	}
	initCCIPReport, err := operations.ExecuteOperation(b, operation.InitializeCCIPOp, deps, initCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, initCCIPReport.Output...)

	return DeployCCIPSeqOutput{
		CCIPAddress:    ccipAddress,
		MCMSOperations: mcmsOperations,
	}, nil
}
