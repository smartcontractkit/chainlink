package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/mcms"
)

// Deploy CCIP Sequence
type DeployCCIPSeqInput struct {
	MCMSAddress aptos.AccountAddress
	MCMSOpCount uint64
	CCIPConfig  config.ChainContractParams
}

type DeployCCIPSeqOutput struct {
	CCIPAddress   aptos.AccountAddress
	MCMSProposals []*mcms.Proposal
	NextOpCount   uint64
}

var DeployCCIPSequence = operations.NewSequence(
	"deploy-aptos-ccip-sequence",
	operation.Version1_0_0,
	"Deploy Aptos CCIP contracts and initialize them",
	deployCCIPSequence,
)

func deployCCIPSequence(b operations.Bundle, deps operation.AptosDeps, in DeployCCIPSeqInput) (DeployCCIPSeqOutput, error) {
	var proposals []*mcms.Proposal
	mcmsOpCount := in.MCMSOpCount

	// Cleanup staging area
	cleanupInput := operation.CleanupStagingAreaInput{
		MCMSAddress: in.MCMSAddress,
		MCMSOpCount: mcmsOpCount,
	}
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, cleanupInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	if cleanupReport.Output.MCMSProposal != nil {
		proposals = append(proposals, cleanupReport.Output.MCMSProposal)
		mcmsOpCount = cleanupReport.Output.NextOpCount
	}

	// Generate proposal to deploy CCIP package
	deployCCIPInput := operation.DeployCCIPInput{
		MCMSAddress: in.MCMSAddress,
		MCMSOpCount: mcmsOpCount,
	}
	deployCCIPReport, err := operations.ExecuteOperation(b, operation.GenerateDeployCCIPProposalOp, deps, deployCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	ccipAddress := deployCCIPReport.Output.CCIPAddress
	if deployCCIPReport.Output.MCMSProposal != nil {
		proposals = append(proposals, deployCCIPReport.Output.MCMSProposal)
		mcmsOpCount = deployCCIPReport.Output.NextOpCount
	}

	// Generate proposal to deploy Router module
	deployRouterInput := operation.DeployRouterInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
		MCMSOpCount: mcmsOpCount,
	}
	deployRouterReport, err := operations.ExecuteOperation(b, operation.GenerateDeployRouterProposalOp, deps, deployRouterInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	if deployRouterReport.Output.MCMSProposal != nil {
		proposals = append(proposals, deployRouterReport.Output.MCMSProposal)
		mcmsOpCount = deployRouterReport.Output.NextOpCount
	}

	// Generate proposal to initialize CCIP
	initCCIPInput := operation.InitializeCCIPInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: ccipAddress,
		CCIPConfig:  in.CCIPConfig,
		MCMSOpCount: mcmsOpCount,
	}
	initCCIPReport, err := operations.ExecuteOperation(b, operation.InitializeCCIPOp, deps, initCCIPInput)
	if err != nil {
		return DeployCCIPSeqOutput{}, err
	}
	if initCCIPReport.Output.MCMSProposal != nil {
		proposals = append(proposals, initCCIPReport.Output.MCMSProposal)
		mcmsOpCount = initCCIPReport.Output.NextOpCount
	}

	return DeployCCIPSeqOutput{
		CCIPAddress:   ccipAddress,
		MCMSProposals: proposals,
		NextOpCount:   mcmsOpCount,
	}, nil
}
