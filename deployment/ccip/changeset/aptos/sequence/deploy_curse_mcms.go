package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// DeployCurseMCMSSeqInput holds the configuration for deploying CurseMCMS.
type DeployCurseMCMSSeqInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPAddress aptos.AccountAddress
	CurseMCMS   types.MCMSWithTimelockConfigV2
}

// DeployCurseMCMSSeqOutput holds the deployed address and a CurseMCMS
// self-governance batch (AcceptOwnership + SetMinDelay) that must be submitted
// as a proposal targeting the CurseMCMS contract.
type DeployCurseMCMSSeqOutput struct {
	CurseMCMSAddress   aptos.AccountAddress
	CurseMCMSOperation mcmstypes.BatchOperation
}

var DeployCurseMCMSSequence = operations.NewSequence(
	"deploy-aptos-curse-mcms-sequence",
	operation.Version1_0_0,
	"Deploy and configure Aptos CurseMCMS contract",
	deployCurseMCMSSequence,
)

func deployCurseMCMSSequence(b operations.Bundle, deps dependency.AptosDeps, in DeployCurseMCMSSeqInput) (DeployCurseMCMSSeqOutput, error) {
	// Check if CurseMCMS is already deployed
	onChainState := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector]
	if onChainState.CurseMCMSAddress != (aptos.AccountAddress{}) {
		b.Logger.Infow("CurseMCMS already deployed", "addr", onChainState.CurseMCMSAddress.StringLong())
		return DeployCurseMCMSSeqOutput{}, nil
	}

	// Deploy CurseMCMS
	deployReport, err := operations.ExecuteOperation(b, operation.DeployCurseMCMSOp, deps, operation.DeployCurseMCMSInput{
		MCMSAddress: in.MCMSAddress,
		CCIPAddress: in.CCIPAddress,
	})
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}
	curseMCMSAddr := deployReport.Output

	// Configure CurseMCMS – bypasser
	_, err = operations.ExecuteOperation(b, operation.ConfigureCurseMCMSOp, deps, operation.ConfigureCurseMCMSInput{
		CurseMCMSAddress: curseMCMSAddr,
		MCMSConfigs:      in.CurseMCMS.Bypasser,
		MCMSRole:         aptosmcms.TimelockRoleBypasser,
	})
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}

	// Configure CurseMCMS – canceller
	_, err = operations.ExecuteOperation(b, operation.ConfigureCurseMCMSOp, deps, operation.ConfigureCurseMCMSInput{
		CurseMCMSAddress: curseMCMSAddr,
		MCMSConfigs:      in.CurseMCMS.Canceller,
		MCMSRole:         aptosmcms.TimelockRoleCanceller,
	})
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}

	// Configure CurseMCMS – proposer
	_, err = operations.ExecuteOperation(b, operation.ConfigureCurseMCMSOp, deps, operation.ConfigureCurseMCMSInput{
		CurseMCMSAddress: curseMCMSAddr,
		MCMSConfigs:      in.CurseMCMS.Proposer,
		MCMSRole:         aptosmcms.TimelockRoleProposer,
	})
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}

	// Transfer ownership to self (deployer-signed)
	_, err = operations.ExecuteOperation(b, operation.TransferCurseMCMSOwnershipToSelfOp, deps, curseMCMSAddr)
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}

	// Encode AcceptOwnership as a CurseMCMS proposal transaction
	aoReport, err := operations.ExecuteOperation(b, operation.AcceptCurseMCMSOwnershipOp, deps, curseMCMSAddr)
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}

	// Encode SetMinDelay as a CurseMCMS proposal transaction
	mdReport, err := operations.ExecuteOperation(b, operation.SetCurseMCMSMinDelayOp, deps, operation.CurseMCMSMinDelayInput{
		CurseMCMSAddress: curseMCMSAddr,
		TimelockMinDelay: (*in.CurseMCMS.TimelockMinDelay).Uint64(),
	})
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}

	curseMCMSOp := mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{aoReport.Output, mdReport.Output},
	}

	return DeployCurseMCMSSeqOutput{
		CurseMCMSAddress:   curseMCMSAddr,
		CurseMCMSOperation: curseMCMSOp,
	}, nil
}
