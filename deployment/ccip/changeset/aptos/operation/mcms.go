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
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
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
	MCMSRole    aptosmcms.TimelockRole
}

var ConfigureMCMSOp = operations.NewOperation(
	"configure-mcms-op",
	Version1_0_0,
	"Configure MCMS Contract Operation for Aptos Chain",
	configureMCMS,
)

func configureMCMS(b operations.Bundle, deps AptosDeps, in ConfigureMCMSInput) (any, error) {
	configurer := aptosmcms.NewConfigurer(deps.AptosChain.Client, deps.AptosChain.DeployerSigner, in.MCMSRole)
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

var GenerateAcceptOwnershipProposalOp = operations.NewOperation(
	"generate-accept-ownership-proposal-op",
	Version1_0_0,
	"Generate Accept Ownership Proposal for MCMS Contract",
	generateAcceptOwnershipProposal,
)

func generateAcceptOwnershipProposal(b operations.Bundle, deps AptosDeps, in GenerateAcceptOwnershipProposalInput) (mcmstypes.BatchOperation, error) {
	moduleInfo, function, _, args, err := (*in.ContractMCMS).MCMSAccount().Encoder().AcceptOwnership()
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	callOneAdditionalFields, err := json.Marshal(additionalFields)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to marshal additionalFields: %w", err)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions: []mcmstypes.Transaction{{
			To:               in.AddressMCMS.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: callOneAdditionalFields,
		}},
	}, err
}
