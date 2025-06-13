package v1_5_1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings"
	burn_mint_external "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/burn_mint_with_external_minter_fast_transfer_token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
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

var (
	// BurnMint Fast Transfer Token Pool Operations
	BurnMintFastTransferTokenPoolUpdateDestChainConfigOp = opsutil.NewEVMCallOperation(
		"BurnMintFastTransferTokenPoolUpdateDestChainConfigOp",
		semver.MustParse("1.0.0"),
		"Update destination chain configurations on BurnMint fast transfer token pool contract",
		burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolABI,
		shared.BurnMintFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (interface{}, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintFastTransferTokenPool)
		},
		func(pool interface{}, opts *bind.TransactOpts, input UpdateDestChainConfigInput) (*types.Transaction, error) {
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
		func(address common.Address, backend bind.ContractBackend) (interface{}, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintFastTransferTokenPool)
		},
		func(pool interface{}, opts *bind.TransactOpts, input UpdateFillerAllowlistInput) (*types.Transaction, error) {
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
		func(address common.Address, backend bind.ContractBackend) (interface{}, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintWithExternalMinterFastTransferTokenPool)
		},
		func(pool interface{}, opts *bind.TransactOpts, input UpdateDestChainConfigInput) (*types.Transaction, error) {
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
		func(address common.Address, backend bind.ContractBackend) (interface{}, error) {
			return bindings.NewFastTransferTokenPoolWrapper(address, backend, shared.BurnMintWithExternalMinterFastTransferTokenPool)
		},
		func(pool interface{}, opts *bind.TransactOpts, input UpdateFillerAllowlistInput) (*types.Transaction, error) {
			wrapper := pool.(*bindings.FastTransferTokenPoolWrapper)
			return wrapper.UpdateFillerAllowList(opts, input.AddFillers, input.RemoveFillers)
		},
	)
)
