package v1_5_1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings"
	burn_mint_external "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/burn_mint_with_external_minter_fast_transfer_token_pool"
	hybrid_external "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/hybrid_with_external_minter_fast_transfer_token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

// UpdateDestChainConfigInput defines the input for updating destination chain configuration
type UpdateDestChainConfigInput struct {
	Updates []bindings.DestChainConfigUpdateArgs
}

// UpdateFillerAllowlistInput defines the input for updating filler allowlist
type UpdateFillerAllowlistInput struct {
	AddFillers    []common.Address
	RemoveFillers []common.Address
}

// WithdrawPoolFeesInput defines the input for withdrawing pool fees
type WithdrawPoolFeesInput struct {
	Recipient common.Address
}

var (
	// BurnMint Fast Transfer Token Pool Operations
	BurnMintFastTransferTokenPoolUpdateDestChainConfigOp = opsutil.NewEVMCallOperation(
		"BurnMintFastTransferTokenPoolUpdateDestChainConfigOp",
		semver.MustParse("1.0.0"),
		"Update destination chain configurations on BurnMint fast transfer token pool contract",
		burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolABI,
		shared.BurnMintFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input UpdateDestChainConfigInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.UpdateDestChainConfig(opts, input.Updates)
		},
	)

	BurnMintFastTransferTokenPoolUpdateFillerAllowlistOp = opsutil.NewEVMCallOperation(
		"BurnMintFastTransferTokenPoolUpdateFillerAllowlistOp",
		semver.MustParse("1.0.0"),
		"Update filler allowlist on BurnMint fast transfer token pool contract",
		burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolABI,
		shared.BurnMintFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input UpdateFillerAllowlistInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.UpdateFillerAllowList(opts, input.AddFillers, input.RemoveFillers)
		},
	)

	// BurnMintWithExternalMinter Fast Transfer Token Pool Operations
	BurnMintWithExternalMinterFastTransferTokenPoolUpdateDestChainConfigOp = opsutil.NewEVMCallOperation(
		"BurnMintWithExternalMinterFastTransferTokenPoolUpdateDestChainConfigOp",
		semver.MustParse("1.0.0"),
		"Update destination chain configurations on BurnMintWithExternalMinter fast transfer token pool contract",
		burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolABI,
		shared.BurnMintWithExternalMinterFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintWithExternalMinterFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input UpdateDestChainConfigInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.UpdateDestChainConfig(opts, input.Updates)
		},
	)

	BurnMintWithExternalMinterFastTransferTokenPoolUpdateFillerAllowlistOp = opsutil.NewEVMCallOperation(
		"BurnMintWithExternalMinterFastTransferTokenPoolUpdateFillerAllowlistOp",
		semver.MustParse("1.0.0"),
		"Update filler allowlist on BurnMintWithExternalMinter fast transfer token pool contract",
		burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolABI,
		shared.BurnMintWithExternalMinterFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintWithExternalMinterFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input UpdateFillerAllowlistInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.UpdateFillerAllowList(opts, input.AddFillers, input.RemoveFillers)
		},
	)

	// BurnMint Fast Transfer Token Pool Withdraw Operations
	BurnMintFastTransferTokenPoolWithdrawPoolFeesOp = opsutil.NewEVMCallOperation(
		"BurnMintFastTransferTokenPoolWithdrawPoolFeesOp",
		semver.MustParse("1.0.0"),
		"Withdraw pool fees from BurnMint fast transfer token pool contract",
		burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolABI,
		shared.BurnMintFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input WithdrawPoolFeesInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.WithdrawPoolFees(opts, input.Recipient)
		},
	)

	// BurnMintWithExternalMinter Fast Transfer Token Pool Withdraw Operations
	BurnMintWithExternalMinterFastTransferTokenPoolWithdrawPoolFeesOp = opsutil.NewEVMCallOperation(
		"BurnMintWithExternalMinterFastTransferTokenPoolWithdrawPoolFeesOp",
		semver.MustParse("1.0.0"),
		"Withdraw pool fees from BurnMintWithExternalMinter fast transfer token pool contract",
		burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolABI,
		shared.BurnMintWithExternalMinterFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintWithExternalMinterFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input WithdrawPoolFeesInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.WithdrawPoolFees(opts, input.Recipient)
		},
	)

	// HybridWithExternalMinter Fast Transfer Token Pool Operations
	HybridWithExternalMinterFastTransferTokenPoolUpdateDestChainConfigOp = opsutil.NewEVMCallOperation(
		"HybridWithExternalMinterFastTransferTokenPoolUpdateDestChainConfigOp",
		semver.MustParse("1.0.0"),
		"Update destination chain configurations on HybridWithExternalMinter fast transfer token pool contract",
		hybrid_external.HybridWithExternalMinterFastTransferTokenPoolABI,
		shared.HybridWithExternalMinterFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.HybridWithExternalMinterFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input UpdateDestChainConfigInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.UpdateDestChainConfig(opts, input.Updates)
		},
	)

	HybridWithExternalMinterFastTransferTokenPoolUpdateFillerAllowlistOp = opsutil.NewEVMCallOperation(
		"HybridWithExternalMinterFastTransferTokenPoolUpdateFillerAllowlistOp",
		semver.MustParse("1.0.0"),
		"Update filler allowlist on HybridWithExternalMinter fast transfer token pool contract",
		hybrid_external.HybridWithExternalMinterFastTransferTokenPoolABI,
		shared.HybridWithExternalMinterFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.HybridWithExternalMinterFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input UpdateFillerAllowlistInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.UpdateFillerAllowList(opts, input.AddFillers, input.RemoveFillers)
		},
	)

	HybridWithExternalMinterFastTransferTokenPoolWithdrawPoolFeesOp = opsutil.NewEVMCallOperation(
		"HybridWithExternalMinterFastTransferTokenPoolWithdrawPoolFeesOp",
		semver.MustParse("1.0.0"),
		"Withdraw pool fees from HybridWithExternalMinter fast transfer token pool contract",
		hybrid_external.HybridWithExternalMinterFastTransferTokenPoolABI,
		shared.HybridWithExternalMinterFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (any, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.HybridWithExternalMinterFastTransferTokenPool)
		},
		func(pool any, opts *bind.TransactOpts, input WithdrawPoolFeesInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.WithdrawPoolFees(opts, input.Recipient)
		},
	)
)
