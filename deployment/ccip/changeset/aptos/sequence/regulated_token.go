package sequence

import (
	"errors"
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

	// To be used for tests, staging and testnet
	if in.TokenMint != nil {
		deployerAddress := deps.AptosChain.DeployerSigner.AccountAddress()
		_, err = operations.ExecuteOperation(b, operation.GrantRegulatedTokenMinterRoleOp, deps, operation.GrantRegulatedTokenMinterRoleInput{
			TokenCodeObjectAddress: codeObj,
			Grantee:                deployerAddress,
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
		// Revoke the deployer's MINTER_ROLE so it doesn't retain mint authority once admin
		// is handed over to MCMS. Must happen while the deployer is still admin (i.e., before
		// MCMS executes accept_admin from the proposal generated below).
		_, err = operations.ExecuteOperation(b, operation.RevokeRegulatedTokenMinterRoleOp, deps, operation.RevokeRegulatedTokenMinterRoleInput{
			TokenCodeObjectAddress: codeObj,
			Account:                deployerAddress,
		})
		if err != nil {
			return DeployRegulatedTokenSeqOutput{}, fmt.Errorf("RevokeRegulatedTokenMinterRoleOp: %w", err)
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

// FinalizeRegulatedTokenOwnershipSequence runs regulated_token::execute_ownership_transfer
// (EOA) to finalize the 3-step ownable handoff after MCMS has accepted ownership. The full
// validation (transfer is to MCMS, accepted, deployer is current owner) lives in
// VerifyFinalizeRegulatedTokenOwnership and must be invoked by the calling changeset's
// VerifyPreconditions;
var FinalizeRegulatedTokenOwnershipSequence = operations.NewSequence(
	"finalize-regulated-token-ownership-sequence",
	operation.Version1_0_0,
	"Run regulated_token::execute_ownership_transfer (EOA) when a pending transfer exists",
	finalizeRegulatedTokenOwnershipSequence,
)

func finalizeRegulatedTokenOwnershipSequence(
	b operations.Bundle,
	deps dependency.AptosDeps,
	tokenCodeObjectAddress aptos.AccountAddress,
) (aptos.AccountAddress, error) {
	token := regulated_token.Bind(tokenCodeObjectAddress, deps.AptosChain.Client)
	hasPending, err := token.RegulatedToken().HasPendingTransfer(nil)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("query has_pending_transfer: %w", err)
	}
	if !hasPending {
		return aptos.AccountAddress{}, nil
	}
	pendingTo, err := token.RegulatedToken().PendingTransferTo(nil)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("query pending_transfer_to: %w", err)
	}
	if pendingTo == nil {
		return aptos.AccountAddress{}, errors.New("pending_transfer_to is empty despite has_pending_transfer=true")
	}
	_, err = operations.ExecuteOperation(
		b,
		operation.ExecuteRegulatedTokenOwnershipTransferOp,
		deps,
		operation.ExecuteRegulatedTokenOwnershipTransferInput{
			TokenCodeObjectAddress: tokenCodeObjectAddress,
			To:                     *pendingTo,
		},
	)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("ExecuteRegulatedTokenOwnershipTransferOp: %w", err)
	}
	return *pendingTo, nil
}

// VerifyFinalizeRegulatedTokenOwnership validates that the regulated token is in a state where
// the deployer EOA can run execute_ownership_transfer to finalize the 3-step handoff to MCMS.
func VerifyFinalizeRegulatedTokenOwnership(
	deps dependency.AptosDeps,
	tokenCodeObjectAddress aptos.AccountAddress,
	mcmsAddress aptos.AccountAddress,
) error {
	token := regulated_token.Bind(tokenCodeObjectAddress, deps.AptosChain.Client)

	hasPending, err := token.RegulatedToken().HasPendingTransfer(nil)
	if err != nil {
		return fmt.Errorf("query has_pending_transfer: %w", err)
	}
	if !hasPending {
		return nil
	}

	mcmsContract := mcmsbind.Bind(mcmsAddress, deps.AptosChain.Client)
	expectedOwner, err := mcmsContract.MCMSRegistry().GetPreexistingCodeObjectOwnerAddress(nil, tokenCodeObjectAddress)
	if err != nil {
		return fmt.Errorf("get mcms preexisting code object owner: %w", err)
	}

	pendingTo, err := token.RegulatedToken().PendingTransferTo(nil)
	if err != nil {
		return fmt.Errorf("query pending_transfer_to: %w", err)
	}
	if pendingTo == nil {
		return errors.New("pending_transfer_to is empty despite has_pending_transfer=true")
	}
	if *pendingTo != expectedOwner {
		return fmt.Errorf("pending_transfer_to %s does not match MCMS registry owner %s", pendingTo.StringLong(), expectedOwner.StringLong())
	}

	accepted, err := token.RegulatedToken().PendingTransferAccepted(nil)
	if err != nil {
		return fmt.Errorf("query pending_transfer_accepted: %w", err)
	}
	if accepted == nil || !*accepted {
		return errors.New("pending transfer has not been accepted by MCMS")
	}

	// The DeployRegulatedToken changeset emits accept_ownership and accept_admin in the same
	// MCMS proposal, so once ownership is accepted we expect admin to be MCMS as well. The
	// regulated pool deployment later emits an MCMS GrantRole for BRIDGE_MINTER_OR_BURNER_ROLE
	// that requires MCMS to be admin at execution time; finalizing ownership here without that
	// would produce a pool proposal that fails at execution.
	currentAdmin, err := token.RegulatedToken().Admin(nil)
	if err != nil {
		return fmt.Errorf("query admin: %w", err)
	}
	if currentAdmin != expectedOwner {
		return fmt.Errorf("admin %s has not been transferred to MCMS registry owner %s; accept_admin from the deploy proposal must execute before finalizing ownership", currentAdmin.StringLong(), expectedOwner.StringLong())
	}

	currentOwner, err := token.RegulatedToken().Owner(nil)
	if err != nil {
		return fmt.Errorf("query owner: %w", err)
	}
	deployerAddress := deps.AptosChain.DeployerSigner.AccountAddress()
	if currentOwner != deployerAddress {
		return fmt.Errorf("current owner %s is not the deployer EOA %s; only the original owner can execute the ownership transfer", currentOwner.StringLong(), deployerAddress.StringLong())
	}

	return nil
}
