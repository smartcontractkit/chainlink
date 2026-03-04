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

// DeployCurseMCMSSeqOutput holds the deployed address and any main-MCMS
// operations that the caller must execute (e.g. registering the CurseMCMS
// signer as an allowed curser on RMN Remote).
type DeployCurseMCMSSeqOutput struct {
	CurseMCMSAddress aptos.AccountAddress
	MCMSOperation    mcmstypes.BatchOperation
}

var DeployCurseMCMSSequence = operations.NewSequence(
	"deploy-aptos-curse-mcms-sequence",
	operation.Version1_0_0,
	"Deploy Aptos CurseMCMS contract, configure it, and register as allowed curser",
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

	// Generate main-MCMS transaction to register the CurseMCMS signer as an
	// allowed curser on RMN Remote. The caller must include this in a proposal
	// targeting the main MCMS contract.
	initCursersReport, err := operations.ExecuteOperation(b, operation.InitializeAllowedCursersOp, deps, operation.InitializeAllowedCursersInput{
		CCIPAddress:      in.CCIPAddress,
		CurseMCMSAddress: curseMCMSAddr,
	})
	if err != nil {
		return DeployCurseMCMSSeqOutput{}, err
	}

	mcmsOp := mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{initCursersReport.Output},
	}

	return DeployCurseMCMSSeqOutput{
		CurseMCMSAddress: curseMCMSAddr,
		MCMSOperation:    mcmsOp,
	}, nil
}
