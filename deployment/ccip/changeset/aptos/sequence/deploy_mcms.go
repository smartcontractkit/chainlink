package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Deploy MCMS Sequence
type DeployMCMSSeqOutput struct {
	MCMSAddress  aptos.AccountAddress
	MCMSProposal *mcms.Proposal
	NextOpCount  uint64
}

var DeployMCMSSequence = operations.NewSequence(
	"deploy-aptos-mcms-sequence",
	operation.Version1_0_0,
	"Deploy Aptos MCMS contract and configure it",
	deployMCMSSequence,
)

func deployMCMSSequence(b operations.Bundle, deps operation.AptosDeps, configMCMS mcmstypes.Config) (DeployMCMSSeqOutput, error) {
	// Check if MCMS package is already deployed
	if deps.OnChainState.MCMSAddress != aptos.AccountZero {
		b.Logger.Infow("MCMS Package already deployed", "addr", deps.OnChainState.MCMSAddress.String())
		return DeployMCMSSeqOutput{}, nil
	}
	// Deploy MCMS
	deployMCMSReport, err := operations.ExecuteOperation(b, operation.DeployMCMSOp, deps, operations.EmptyInput{})
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Configure MCMS
	configureMCMSInput := operation.ConfigureMCMSInput{
		AddressMCMS: deployMCMSReport.Output.AddressMCMS,
		MCMSConfigs: configMCMS,
	}
	_, err = operations.ExecuteOperation(b, operation.ConfigureMCMSOp, deps, configureMCMSInput)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Transfer ownership to self
	_, err = operations.ExecuteOperation(b, operation.TransferOwnershipToSelfOp, deps, deployMCMSReport.Output.ContractMCMS)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Generate proposal to accept ownership
	generateAcceptOwnershipProposalInput := operation.GenerateAcceptOwnershipProposalInput{
		AddressMCMS:  deployMCMSReport.Output.AddressMCMS,
		ContractMCMS: deployMCMSReport.Output.ContractMCMS,
	}
	gaopReport, err := operations.ExecuteOperation(b, operation.GenerateAcceptOwnershipProposalOp, deps, generateAcceptOwnershipProposalInput)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}

	return DeployMCMSSeqOutput{
		MCMSAddress:  deployMCMSReport.Output.AddressMCMS,
		MCMSProposal: gaopReport.Output.MCMSProposal,
		NextOpCount:  gaopReport.Output.NextOpCount,
	}, nil
}
