package operation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/operations"
)

const AcceptOwnershipProposalDescription = "Accept ownership of the contract to self"

var DeployMCMSOp = operations.NewOperation(
	"deploy-mcms-op",
	Version1_0_0,
	"Deploys MCMS Contract Operation for Aptos Chain",
	deployMCMS,
)

type DeployMCMSOutput struct {
	AddressMCMS  aptos.AccountAddress
	ContractMCMS *mcmsbind.MCMS // TODO: outputs should be serializable
}

func deployMCMS(b operations.Bundle, deps AptosDeps, input operations.EmptyInput) (DeployMCMSOutput, error) {
	mcmsSeed := mcmsbind.DefaultSeed + time.Now().String()
	addressMCMS, mcmsDeployTx, contractMCMS, err := mcmsbind.DeployToResourceAccount(deps.AptosChain.DeployerSigner, deps.AptosChain.Client, mcmsSeed)
	if err != nil {
		return DeployMCMSOutput{}, fmt.Errorf("failed to deploy MCMS contract: %v", err)
	}
	if err := utils.ConfirmTx(deps.AptosChain, mcmsDeployTx.Hash); err != nil {
		return DeployMCMSOutput{}, fmt.Errorf("failed to confirm MCMS deployment transaction: %v", err)
	}

	typeAndVersion := deployment.NewTypeAndVersion(changeset.AptosMCMSType, deployment.Version1_0_0)
	deps.AB.Save(deps.AptosChain.Selector, addressMCMS.String(), typeAndVersion)
	return DeployMCMSOutput{addressMCMS, &contractMCMS}, nil
}

type ConfigureMCMSInput struct {
	AddressMCMS aptos.AccountAddress
	MCMSConfigs mcmstypes.Config
}

var ConfigureMCMSOp = operations.NewOperation(
	"configure-mcms-op",
	Version1_0_0,
	"Configure MCMS Contract Operation for Aptos Chain",
	configureMCMS,
)

func configureMCMS(b operations.Bundle, deps AptosDeps, in ConfigureMCMSInput) (any, error) {
	configurer := aptosmcms.NewConfigurer(deps.AptosChain.Client, deps.AptosChain.DeployerSigner)
	setCfgTx, err := configurer.SetConfig(context.Background(), in.AddressMCMS.StringLong(), &in.MCMSConfigs, false)
	if err != nil {
		return nil, fmt.Errorf("failed to setConfig in MCMS contract: %w", err)
	}
	if err := utils.ConfirmTx(deps.AptosChain, setCfgTx.Hash); err != nil {
		return nil, fmt.Errorf("MCMS setConfig transaction failed: %w", err)
	}
	return nil, nil
}

var TransferOwnershipToSelfOp = operations.NewOperation(
	"transfer-ownership-to-self-op",
	Version1_0_0,
	"Transfer ownership to self",
	transferOwnershipToSelf,
)

func transferOwnershipToSelf(b operations.Bundle, deps AptosDeps, contractMCMS *mcmsbind.MCMS) (any, error) {
	opts := &bind.TransactOpts{Signer: deps.AptosChain.DeployerSigner}
	tx, err := (*contractMCMS).MCMSAccount().TransferOwnershipToSelf(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to TransferOwnershipToSelf in MCMS contract: %w", err)
	}
	_, err = deps.AptosChain.Client.WaitForTransaction(tx.Hash)
	if err != nil {
		return nil, fmt.Errorf("MCMS TransferOwnershipToSelf transaction failed: %w", err)
	}
	return nil, nil
}

type GenerateAcceptOwnershipProposalInput struct {
	AddressMCMS  aptos.AccountAddress
	ContractMCMS *mcmsbind.MCMS // TODO: outputs should be serializable
}

type GenerateAcceptOwnershipProposalOutput struct {
	MCMSProposal *mcms.Proposal
	NextOpCount  uint64
}

var GenerateAcceptOwnershipProposalOp = operations.NewOperation(
	"generate-accept-ownership-proposal-op",
	Version1_0_0,
	"Generate Accept Ownership Proposal for MCMS Contract",
	generateAcceptOwnershipProposal,
)

func generateAcceptOwnershipProposal(b operations.Bundle, deps AptosDeps, in GenerateAcceptOwnershipProposalInput) (GenerateAcceptOwnershipProposalOutput, error) {
	moduleInfo, function, _, args, err := (*in.ContractMCMS).MCMSAccount().Encoder().AcceptOwnership()
	if err != nil {
		return GenerateAcceptOwnershipProposalOutput{}, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	callOneAdditionalFields, err := json.Marshal(additionalFields)
	if err != nil {
		return GenerateAcceptOwnershipProposalOutput{}, fmt.Errorf("failed to marshal additionalFields: %w", err)
	}
	mcmsOps := []mcmstypes.Operation{{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transaction: mcmstypes.Transaction{
			To:               in.AddressMCMS.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: callOneAdditionalFields,
		},
	}}
	mcmsProposal, nextOpCount, err := utils.GenerateProposal(
		deps.AptosChain.Client,
		(*in.ContractMCMS).Address(),
		deps.AptosChain.Selector,
		mcmsOps,
		AcceptOwnershipProposalDescription,
		0,
	)

	return GenerateAcceptOwnershipProposalOutput{mcmsProposal, nextOpCount}, err
}
