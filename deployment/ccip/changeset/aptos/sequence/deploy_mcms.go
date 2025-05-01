package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Deploy MCMS Sequence
type DeployMCMSSeqOutput struct {
	MCMSAddress   aptos.AccountAddress
	MCMSOperation mcmstypes.BatchOperation
}

var DeployMCMSSequence = operations.NewSequence(
	"deploy-aptos-mcms-sequence",
	operation.Version1_0_0,
	"Deploy Aptos MCMS contract and configure it",
	deployMCMSSequence,
)

func deployMCMSSequence(b operations.Bundle, deps operation.AptosDeps, configMCMS types.MCMSWithTimelockConfigV2) (DeployMCMSSeqOutput, error) {
	// Check if MCMS package is already deployed
	if deps.OnChainState.MCMSAddress != (aptos.AccountAddress{}) {
		b.Logger.Infow("MCMS Package already deployed", "addr", deps.OnChainState.MCMSAddress.String())
		return DeployMCMSSeqOutput{}, nil
	}
	// Deploy MCMS
	deployMCMSReport, err := operations.ExecuteOperation(b, operation.DeployMCMSOp, deps, operations.EmptyInput{})
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Configure MCMS
	configureMCMSBypassers := operation.ConfigureMCMSInput{
		AddressMCMS: deployMCMSReport.Output.AddressMCMS,
		MCMSConfigs: configMCMS.Bypasser,
		MCMSRole:    aptosmcms.TimelockRoleBypasser,
	}
	_, err = operations.ExecuteOperation(b, operation.ConfigureMCMSOp, deps, configureMCMSBypassers)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	configureMCMSCancellers := operation.ConfigureMCMSInput{
		AddressMCMS: deployMCMSReport.Output.AddressMCMS,
		MCMSConfigs: configMCMS.Canceller,
		MCMSRole:    aptosmcms.TimelockRoleCanceller,
	}
	_, err = operations.ExecuteOperation(b, operation.ConfigureMCMSOp, deps, configureMCMSCancellers)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	configureMCMSProposers := operation.ConfigureMCMSInput{
		AddressMCMS: deployMCMSReport.Output.AddressMCMS,
		MCMSConfigs: configMCMS.Proposer,
		MCMSRole:    aptosmcms.TimelockRoleProposer,
	}
	_, err = operations.ExecuteOperation(b, operation.ConfigureMCMSOp, deps, configureMCMSProposers)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// TODO: Should set MinDelay to timelock
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
		MCMSAddress:   deployMCMSReport.Output.AddressMCMS,
		MCMSOperation: gaopReport.Output,
	}, nil
}
