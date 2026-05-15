package operation

import (
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/regulated_token"
	module_regulated_token "github.com/smartcontractkit/chainlink-aptos/bindings/regulated_token/regulated_token"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
)

// regulated_token cannot be deployed via MCMS: MCMS dispatches calls through
// dispatchable_fungible_asset (DFA), and regulated_token::initialize itself calls
// dispatchable_fungible_asset::register_dispatch_functions, which Aptos rejects as
// nested DFA dispatch (re-entrancy). Every operation in this file is therefore
// signed directly by the deployer and confirmed on chain — they do not produce
// MCMS transactions and must not be wrapped in an MCMS proposal.

// DeployRegulatedTokenObjectOp publishes the regulated_token package to a new code object (one tx).
var DeployRegulatedTokenObjectOp = operations.NewOperation(
	"deploy-regulated-token-object-op",
	Version1_0_0,
	"Deploy regulated_token package to a new object",
	deployRegulatedTokenObject,
)

func deployRegulatedTokenObject(
	b operations.Bundle,
	deps dependency.AptosDeps,
	_ operations.EmptyInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	adminAddress := signer.AccountAddress()

	tokenAddress, tx, _, err := regulated_token.DeployToObject(signer, client, adminAddress)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("DeployToObject: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm DeployToObject: %w", err)
	}
	return tokenAddress, nil
}

// DeployRegulatedTokenMCMSRegistrarInput is input for DeployRegulatedTokenMCMSRegistrarOp.
type DeployRegulatedTokenMCMSRegistrarInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	MCMSAddress            aptos.AccountAddress
	RegistrarPreregister   bool
}

// DeployRegulatedTokenMCMSRegistrarOp attaches the MCMS registrar to the token code object (one tx).
var DeployRegulatedTokenMCMSRegistrarOp = operations.NewOperation(
	"deploy-regulated-token-mcms-registrar-op",
	Version1_0_0,
	"Deploy MCMS registrar on existing regulated token code object",
	deployRegulatedTokenMCMSRegistrar,
)

func deployRegulatedTokenMCMSRegistrar(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in DeployRegulatedTokenMCMSRegistrarInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	adminAddress := signer.AccountAddress()

	tx, _, err := regulated_token.DeployMCMSRegistrarToExistingObject(
		signer, client, in.TokenCodeObjectAddress, adminAddress, in.MCMSAddress, in.RegistrarPreregister,
	)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("DeployMCMSRegistrarToExistingObject: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm DeployMCMSRegistrarToExistingObject: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}

// InitializeRegulatedTokenInput is input for InitializeRegulatedTokenOp.
type InitializeRegulatedTokenInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	TokenParams            config.TokenParams
}

// InitializeRegulatedTokenOp runs regulated_token::initialize (one tx).
var InitializeRegulatedTokenOp = operations.NewOperation(
	"initialize-regulated-token-op",
	Version1_0_0,
	"Initialize regulated token metadata and hooks",
	initializeRegulatedToken,
)

func initializeRegulatedToken(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in InitializeRegulatedTokenInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	opts := &bind.TransactOpts{Signer: signer}

	token := regulated_token.Bind(in.TokenCodeObjectAddress, client)

	var maxSupply **big.Int
	if in.TokenParams.MaxSupply != nil {
		p := in.TokenParams.MaxSupply
		maxSupply = &p
	}
	tx, err := token.RegulatedToken().Initialize(
		opts,
		maxSupply,
		in.TokenParams.Name,
		string(in.TokenParams.Symbol),
		in.TokenParams.Decimals,
		in.TokenParams.Icon,
		in.TokenParams.Project,
	)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("initialize: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm initialize: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}

// GrantRegulatedTokenMinterRoleInput is input for GrantRegulatedTokenMinterRoleOp.
type GrantRegulatedTokenMinterRoleInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	Grantee                aptos.AccountAddress
}

// GrantRegulatedTokenMinterRoleOp grants MINTER_ROLE to one account (one tx).
var GrantRegulatedTokenMinterRoleOp = operations.NewOperation(
	"grant-regulated-token-minter-role-op",
	Version1_0_0,
	"Grant regulated token minter role",
	grantRegulatedTokenMinterRole,
)

func grantRegulatedTokenMinterRole(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in GrantRegulatedTokenMinterRoleInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	opts := &bind.TransactOpts{Signer: signer}
	token := regulated_token.Bind(in.TokenCodeObjectAddress, client)

	tx, err := token.RegulatedToken().GrantRole(opts, module_regulated_token.MINTER_ROLE, in.Grantee)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("GrantRole minter: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm GrantRole: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}

// RevokeRegulatedTokenMinterRoleInput is input for RevokeRegulatedTokenMinterRoleOp.
type RevokeRegulatedTokenMinterRoleInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	Account                aptos.AccountAddress
}

// RevokeRegulatedTokenMinterRoleOp revokes MINTER_ROLE from one account (one tx).
var RevokeRegulatedTokenMinterRoleOp = operations.NewOperation(
	"revoke-regulated-token-minter-role-op",
	Version1_0_0,
	"Revoke regulated token minter role from an account",
	revokeRegulatedTokenMinterRole,
)

func revokeRegulatedTokenMinterRole(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in RevokeRegulatedTokenMinterRoleInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	opts := &bind.TransactOpts{Signer: signer}
	token := regulated_token.Bind(in.TokenCodeObjectAddress, client)

	tx, err := token.RegulatedToken().RevokeRole(opts, module_regulated_token.MINTER_ROLE, in.Account)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("revoke role minter: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm revoke role: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}

// MintRegulatedTokenInput is input for MintRegulatedTokenOp.
type MintRegulatedTokenInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	To                     aptos.AccountAddress
	Amount                 uint64
}

// MintRegulatedTokenOp mints to an address (one tx).
var MintRegulatedTokenOp = operations.NewOperation(
	"mint-regulated-token-op",
	Version1_0_0,
	"Mint regulated token to an account",
	mintRegulatedToken,
)

func mintRegulatedToken(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in MintRegulatedTokenInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	opts := &bind.TransactOpts{Signer: signer}
	token := regulated_token.Bind(in.TokenCodeObjectAddress, client)

	tx, err := token.RegulatedToken().Mint(opts, in.To, in.Amount)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("mint: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm mint: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}

// TransferRegulatedTokenOwnershipInput is input for TransferRegulatedTokenOwnershipOp.
type TransferRegulatedTokenOwnershipInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	To                     aptos.AccountAddress
}

// TransferRegulatedTokenOwnershipOp starts ownership transfer to the given address (one tx).
var TransferRegulatedTokenOwnershipOp = operations.NewOperation(
	"transfer-regulated-token-ownership-op",
	Version1_0_0,
	"Transfer regulated token ownership to the given address",
	transferRegulatedTokenOwnership,
)

func transferRegulatedTokenOwnership(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in TransferRegulatedTokenOwnershipInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	opts := &bind.TransactOpts{Signer: signer}
	token := regulated_token.Bind(in.TokenCodeObjectAddress, client)

	tx, err := token.RegulatedToken().TransferOwnership(opts, in.To)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("TransferOwnership: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm TransferOwnership: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}

// ExecuteRegulatedTokenOwnershipTransferInput is input for ExecuteRegulatedTokenOwnershipTransferOp.
type ExecuteRegulatedTokenOwnershipTransferInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	To                     aptos.AccountAddress
}

// ExecuteRegulatedTokenOwnershipTransferOp finalizes the 3-step ownable handoff by calling
// regulated_token::execute_ownership_transfer with the original deployer (one tx). It must
// be called after the new owner has accepted ownership.
var ExecuteRegulatedTokenOwnershipTransferOp = operations.NewOperation(
	"execute-regulated-token-ownership-transfer-op",
	Version1_0_0,
	"Execute regulated token ownership transfer (finalize 3-step handoff)",
	executeRegulatedTokenOwnershipTransfer,
)

func executeRegulatedTokenOwnershipTransfer(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in ExecuteRegulatedTokenOwnershipTransferInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	opts := &bind.TransactOpts{Signer: signer}
	token := regulated_token.Bind(in.TokenCodeObjectAddress, client)

	tx, err := token.RegulatedToken().ExecuteOwnershipTransfer(opts, in.To)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("execute ownership transfer: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm execute ownership transfer: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}

// TransferRegulatedTokenAdminInput is input for TransferRegulatedTokenAdminOp.
type TransferRegulatedTokenAdminInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	NewAdmin               aptos.AccountAddress
}

// TransferRegulatedTokenAdminOp proposes admin role transfer to the given address (one tx).
var TransferRegulatedTokenAdminOp = operations.NewOperation(
	"transfer-regulated-token-admin-op",
	Version1_0_0,
	"Transfer regulated token admin role to the given address",
	transferRegulatedTokenAdmin,
)

func transferRegulatedTokenAdmin(
	b operations.Bundle,
	deps dependency.AptosDeps,
	in TransferRegulatedTokenAdminInput,
) (aptos.AccountAddress, error) {
	_ = b
	signer := deps.AptosChain.DeployerSigner
	client := deps.AptosChain.Client
	opts := &bind.TransactOpts{Signer: signer}
	token := regulated_token.Bind(in.TokenCodeObjectAddress, client)

	tx, err := token.RegulatedToken().TransferAdmin(opts, in.NewAdmin)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("TransferAdmin: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("confirm TransferAdmin: %w", err)
	}
	return in.TokenCodeObjectAddress, nil
}
