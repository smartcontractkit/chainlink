package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"

	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
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
	onChainState := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector]
	if onChainState.MCMSAddress != (aptos.AccountAddress{}) {
		b.Logger.Infow("MCMS Package already deployed", "addr", onChainState.MCMSAddress.String())
		return DeployMCMSSeqOutput{}, nil
	}
	// Deploy MCMS
	deployMCMSReport, err := operations.ExecuteOperation(b, operation.DeployMCMSOp, deps, operations.EmptyInput{})
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Configure MCMS
	configureMCMSBypassers := operation.ConfigureMCMSInput{
		MCMSAddress: deployMCMSReport.Output,
		MCMSConfigs: configMCMS.Bypasser,
		MCMSRole:    aptosmcms.TimelockRoleBypasser,
	}
	_, err = operations.ExecuteOperation(b, operation.ConfigureMCMSOp, deps, configureMCMSBypassers)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	configureMCMSCancellers := operation.ConfigureMCMSInput{
		MCMSAddress: deployMCMSReport.Output,
		MCMSConfigs: configMCMS.Canceller,
		MCMSRole:    aptosmcms.TimelockRoleCanceller,
	}
	_, err = operations.ExecuteOperation(b, operation.ConfigureMCMSOp, deps, configureMCMSCancellers)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	configureMCMSProposers := operation.ConfigureMCMSInput{
		MCMSAddress: deployMCMSReport.Output,
		MCMSConfigs: configMCMS.Proposer,
		MCMSRole:    aptosmcms.TimelockRoleProposer,
	}
	_, err = operations.ExecuteOperation(b, operation.ConfigureMCMSOp, deps, configureMCMSProposers)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Transfer ownership to self
	_, err = operations.ExecuteOperation(b, operation.TransferOwnershipToSelfOp, deps, deployMCMSReport.Output)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Accept ownership
	aoReport, err := operations.ExecuteOperation(b, operation.AcceptOwnershipOp, deps, deployMCMSReport.Output)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}
	// Set MinDelay
	timelockMinDelayInput := operation.TimelockMinDelayInput{
		MCMSAddress:      deployMCMSReport.Output,
		TimelockMinDelay: (*configMCMS.TimelockMinDelay).Uint64(),
	}
	mdReport, err := operations.ExecuteOperation(b, operation.SetMinDelayOP, deps, timelockMinDelayInput)
	if err != nil {
		return DeployMCMSSeqOutput{}, err
	}

	// Generate MCMS Batch Operation
	mcmsOps := mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{aoReport.Output, mdReport.Output},
	}

	return DeployMCMSSeqOutput{
		MCMSAddress:   deployMCMSReport.Output,
		MCMSOperation: mcmsOps,
	}, nil
}
