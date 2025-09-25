package config

import (
	"errors"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type TokenParams struct {
	MaxSupply *big.Int
	Name      string
	Symbol    shared.TokenSymbol
	Decimals  byte
	Icon      string
	Project   string
}

func (tp TokenParams) Validate() error {
	if tp.MaxSupply != nil && tp.MaxSupply.Sign() <= 0 {
		return errors.New("maxSupply must be a positive integer or nil")
	}
	if tp.Name == "" {
		return errors.New("name cannot be empty")
	}
	if tp.Symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	if tp.Decimals < 1 || tp.Decimals > 8 {
		return errors.New("decimals must be between 1 and 8")
	}
	return nil
}

type TokenMint struct {
	Amount uint64
	To     aptos.AccountAddress
}

type DeployTokenFaucetInput struct {
	ChainSelector          uint64
	TokenCodeObjectAddress aptos.AccountAddress
	MCMSConfig             *proposalutils.TimelockConfig
}

type MintTokenInput struct {
	ChainSelector          uint64
	TokenCodeObjectAddress aptos.AccountAddress
	MCMSConfig             *proposalutils.TimelockConfig
	TokenMint
}

// ###################
// # Token Ownership #
// ###################

type TokenTransferInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	TokenType              deployment.ContractType
	To                     aptos.AccountAddress
}

type TransferTokenOwnershipInput struct {
	ChainSelector uint64
	Transfers     []TokenTransferInput
	MCMSConfig    *proposalutils.TimelockConfig
}

type TokenAcceptInput struct {
	TokenCodeObjectAddress aptos.AccountAddress
	TokenType              deployment.ContractType
}

type AcceptTokenOwnershipInput struct {
	ChainSelector uint64
	Accepts       []TokenAcceptInput
	MCMSConfig    *proposalutils.TimelockConfig
}

type ExecuteTokenOwnershipTransferInput struct {
	ChainSelector uint64
	Transfers     []TokenTransferInput
	MCMSConfig    *proposalutils.TimelockConfig
}

type TransferTokenAdminInput struct {
	ChainSelector uint64
	Transfers     []TokenTransferInput
	MCMSConfig    *proposalutils.TimelockConfig
}

type AcceptTokenAdminInput struct {
	ChainSelector uint64
	Accepts       []TokenAcceptInput
	MCMSConfig    *proposalutils.TimelockConfig
}
