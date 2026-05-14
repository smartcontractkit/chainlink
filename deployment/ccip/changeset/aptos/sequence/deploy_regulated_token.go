package sequence

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-aptos/bindings/regulated_token"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

// DeployRegulatedTokenSeqInput is input for DeployRegulatedTokenSequence.
type DeployRegulatedTokenSeqInput struct {
	MCMSAddress          aptos.AccountAddress
	TokenParams          config.TokenParams
	TokenMint            *config.TokenMint
	RegistrarPreregister bool
}

// DeployRegulatedTokenSeqOutput contains token addresses and MCMS batch operations
// (accept_ownership + accept_admin) to be wrapped in a timelock proposal.
type DeployRegulatedTokenSeqOutput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	TokenMetadataAddress   aptos.AccountAddress
	MCMSOperations         []mcmstypes.BatchOperation
}

var DeployRegulatedTokenSequence = operations.NewSequence(
	"deploy-regulated-token-aptos-sequence",
	operation.Version1_0_0,
	"Deploy regulated token directly (regulated_token cannot be deployed via MCMS due to DFA re-entrancy) then build MCMS accept_ownership + accept_admin batch",
	deployRegulatedTokenSequence,
)

func deployRegulatedTokenSequence(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in DeployRegulatedTokenSeqInput,
) (DeployRegulatedTokenSeqOutput, error) {
	objReport, err := operations.ExecuteOperation(b, operation.DeployRegulatedTokenObjectOp, deps, operations.EmptyInput{})
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("DeployRegulatedTokenObjectOp: %w", err)
	}
	codeObj := objReport.Output

	_, err = operations.ExecuteOperation(b, operation.DeployRegulatedTokenMCMSRegistrarOp, deps, operation.DeployRegulatedTokenMCMSRegistrarInput{
		TokenCodeObjectAddress: codeObj,
		MCMSAddress:            in.MCMSAddress,
		RegistrarPreregister:   in.RegistrarPreregister,
	})
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("DeployRegulatedTokenMCMSRegistrarOp: %w", err)
	}

	_, err = operations.ExecuteOperation(b, operation.InitializeRegulatedTokenOp, deps, operation.InitializeRegulatedTokenInput{
		TokenCodeObjectAddress: codeObj,
		TokenParams:            in.TokenParams,
	})
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("InitializeRegulatedTokenOp: %w", err)
	}

	if in.TokenMint != nil {
		grantee := deps.AptosChain.DeployerSigner.AccountAddress()
		_, err = operations.ExecuteOperation(b, operation.GrantRegulatedTokenMinterRoleOp, deps, operation.GrantRegulatedTokenMinterRoleInput{
			TokenCodeObjectAddress: codeObj,
			Grantee:                grantee,
		})
		if err != nil {
			return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("GrantRegulatedTokenMinterRoleOp: %w", err)
		}
		_, err = operations.ExecuteOperation(b, operation.MintRegulatedTokenOp, deps, operation.MintRegulatedTokenInput{
			TokenCodeObjectAddress: codeObj,
			To:                     in.TokenMint.To,
			Amount:                 in.TokenMint.Amount,
		})
		if err != nil {
			return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("MintRegulatedTokenOp: %w", err)
		}
	}

	mcmsContract := mcmsbind.Bind(in.MCMSAddress, deps.AptosChain.Client)
	tokenOwnerAddress, err := mcmsContract.MCMSRegistry().GetPreexistingCodeObjectOwnerAddress(nil, codeObj)
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("GetPreexistingCodeObjectOwnerAddress: %w", err)
	}
	_, err = operations.ExecuteOperation(b, operation.TransferRegulatedTokenOwnershipOp, deps, operation.TransferRegulatedTokenOwnershipInput{
		TokenCodeObjectAddress: codeObj,
		To:                     tokenOwnerAddress,
	})
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("TransferRegulatedTokenOwnershipOp: %w", err)
	}

	_, err = operations.ExecuteOperation(b, operation.TransferRegulatedTokenAdminOp, deps, operation.TransferRegulatedTokenAdminInput{
		TokenCodeObjectAddress: codeObj,
		NewAdmin:               tokenOwnerAddress,
	})
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("TransferRegulatedTokenAdminOp: %w", err)
	}

	token := regulated_token.Bind(codeObj, deps.AptosChain.Client)
	tokenMetadata, err := token.RegulatedToken().TokenMetadata(nil)
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("TokenMetadata: %w", err)
	}

	acceptOwnershipReport, err := operations.ExecuteOperation(
		b,
		operation.AcceptTokenOwnershipOp,
		deps,
		operation.AcceptTokenOwnershipInput{
			TokenCodeObjectAddress: codeObj,
			TokenType:              shared.AptosRegulatedTokenType,
		},
	)
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("AcceptTokenOwnershipOp: %w", err)
	}

	acceptAdminReport, err := operations.ExecuteOperation(
		b,
		operation.AcceptTokenAdminOp,
		deps,
		operation.AcceptTokenAdminInput{
			TokenCodeObjectAddress: codeObj,
		},
	)
	if err != nil {
		return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("AcceptTokenAdminOp: %w", err)
	}

	mcmsOperations := []mcmstypes.BatchOperation{
		{
			ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
			Transactions:  []mcmstypes.Transaction{acceptOwnershipReport.Output, acceptAdminReport.Output},
		},
	}

	return DeployRegulatedTokenSeqOutput{
		TokenCodeObjectAddress: codeObj,
		TokenMetadataAddress:   tokenMetadata,
		MCMSOperations:         mcmsOperations,
	}, nil
}
