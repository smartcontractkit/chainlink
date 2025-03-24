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
)

const AcceptOwnershipProposalDescription = "Accept ownership of the contract to self"

type MCMSDeploymentOperations struct {
	Env          deployment.Environment
	Ab           *deployment.AddressBookMap
	AptosChain   deployment.AptosChain
	OnChainState changeset.AptosCCIPChainState
	MCMSConfigs  mcmstypes.Config
	Proposals    *[]mcms.Proposal
	MCMSOpCount  uint64
}

func (op *MCMSDeploymentOperations) DeployMCMS() (aptos.AccountAddress, mcmsbind.MCMS, error) {
	mcmsSeed := mcmsbind.DefaultSeed + time.Now().String()
	addressMCMS, mcmsDeployTx, contractMCMS, err := mcmsbind.DeployToResourceAccount(op.AptosChain.DeployerSigner, op.AptosChain.Client, mcmsSeed)
	if err != nil {
		return aptos.AccountZero, mcmsbind.MCMSContract{}, fmt.Errorf("failed to deploy MCMS contract: %v", err)
	}
	if err := utils.ConfirmTx(op.AptosChain, mcmsDeployTx.Hash); err != nil {
		return aptos.AccountZero, mcmsbind.MCMSContract{}, fmt.Errorf("failed to confirm MCMS deployment transaction: %v", err)
	}

	typeAndVersion := deployment.NewTypeAndVersion(changeset.AptosMCMSType, deployment.Version1_0_0)
	op.Ab.Save(op.AptosChain.Selector, addressMCMS.String(), typeAndVersion)
	op.OnChainState.MCMSAddress = addressMCMS
	return addressMCMS, contractMCMS, nil
}

func (op *MCMSDeploymentOperations) ConfigureMCMS(addressMCMS aptos.AccountAddress) error {
	configurer := aptosmcms.NewConfigurer(op.AptosChain.Client, op.AptosChain.DeployerSigner)
	setCfgTx, err := configurer.SetConfig(context.Background(), addressMCMS.StringLong(), &op.MCMSConfigs, false)
	if err != nil {
		return fmt.Errorf("failed to setConfig in MCMS contract: %w", err)
	}
	if err := utils.ConfirmTx(op.AptosChain, setCfgTx.Hash); err != nil {
		return fmt.Errorf("MCMS setConfig transaction failed: %w", err)
	}
	return nil
}

func (op *MCMSDeploymentOperations) TransferOwnershipToSelf(contractMCMS mcmsbind.MCMS) error {
	opts := &bind.TransactOpts{Signer: op.AptosChain.DeployerSigner}
	tx, err := contractMCMS.MCMSAccount().TransferOwnershipToSelf(opts)
	if err != nil {
		return fmt.Errorf("failed to TransferOwnershipToSelf in MCMS contract: %w", err)
	}
	_, err = op.AptosChain.Client.WaitForTransaction(tx.Hash)
	if err != nil {
		return fmt.Errorf("MCMS TransferOwnershipToSelf transaction failed: %w", err)
	}
	return nil
}

func (op *MCMSDeploymentOperations) GenerateAcceptOwnershipProposal(addressMCMS aptos.AccountAddress, contractMCMS mcmsbind.MCMS) (*mcms.Proposal, uint64, error) {
	var operations []mcmstypes.Operation
	moduleInfo, function, _, args, err := contractMCMS.MCMSAccount().Encoder().AcceptOwnership()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to encode AcceptOwnership: %w", err)
	}
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	callOneAdditionalFields, err := json.Marshal(additionalFields)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal additionalFields: %w", err)
	}
	operations = append(operations, mcmstypes.Operation{
		ChainSelector: mcmstypes.ChainSelector(op.AptosChain.Selector),
		Transaction: mcmstypes.Transaction{
			To:               addressMCMS.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: callOneAdditionalFields,
		},
	})

	return utils.GenerateProposal(op.AptosChain.Client, contractMCMS.Address(), op.AptosChain.Selector, operations, AcceptOwnershipProposalDescription, 0)
}
