package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment"
	cciptypes "github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type MintERC677Config struct {
	To     common.Address
	Amount *big.Int
}

var MintERC677Op = opsutil.NewEVMCallOperation(
	"MintERC677Op",
	&deployment.Version1_0_0,
	"Mint ERC677 tokens on the specified evm chain",
	burn_mint_erc677.BurnMintERC677ABI,
	cciptypes.ERC677Token,
	burn_mint_erc677.NewBurnMintERC677,
	func(token *burn_mint_erc677.BurnMintERC677, opts *bind.TransactOpts, input MintERC677Config) (*types.Transaction, error) {
		return token.Mint(opts, input.To, input.Amount)
	},
)

var GrantMintAndBurnRolesERC677Op = opsutil.NewEVMCallOperation(
	"GrantMintAndBurnRolesERC677Op",
	&deployment.Version1_0_0,
	"Grant MINTER_ROLE and BURNER_ROLE to the specified ERC677 token address",
	burn_mint_erc677.BurnMintERC677ABI,
	cciptypes.ERC677Token,
	burn_mint_erc677.NewBurnMintERC677,
	func(token *burn_mint_erc677.BurnMintERC677, opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
		return token.GrantMintAndBurnRoles(opts, to)
	},
)

var RevokeMintRoleERC677Op = opsutil.NewEVMCallOperation(
	"RevokeMintRoleERC677Op",
	&deployment.Version1_0_0,
	"Revoke MINTER_ROLE from the specified ERC677 token address",
	burn_mint_erc677.BurnMintERC677ABI,
	cciptypes.ERC677Token,
	burn_mint_erc677.NewBurnMintERC677,
	func(token *burn_mint_erc677.BurnMintERC677, opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
		return token.RevokeMintRole(opts, minter)
	},
)

var RevokeBurnRoleERC677Op = opsutil.NewEVMCallOperation(
	"RevokeBurnRoleERC677Op",
	&deployment.Version1_0_0,
	"Revoke BURNER_ROLE from the specified ERC677 token address",
	burn_mint_erc677.BurnMintERC677ABI,
	cciptypes.ERC677Token,
	burn_mint_erc677.NewBurnMintERC677,
	func(token *burn_mint_erc677.BurnMintERC677, opts *bind.TransactOpts, minter common.Address) (*types.Transaction, error) {
		return token.RevokeBurnRole(opts, minter)
	},
)
