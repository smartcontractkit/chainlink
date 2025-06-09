package shared

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
)

var CurrentTokenPoolVersion = deployment.Version1_5_1

var TokenTypes = map[cldf.ContractType]struct{}{
	BurnMintToken:     {},
	ERC20Token:        {},
	ERC677Token:       {},
	ERC677TokenHelper: {},
}

var TokenPoolTypes = map[cldf.ContractType]struct{}{
	BurnMintTokenPool:              {},
	BurnWithFromMintTokenPool:      {},
	BurnFromMintTokenPool:          {},
	LockReleaseTokenPool:           {},
	USDCTokenPool:                  {},
	HybridLockReleaseUSDCTokenPool: {},
}

var TokenPoolVersions = map[semver.Version]struct{}{
	deployment.Version1_5_1: {},
}

// tokenPool defines behavior common to all token pools.
type tokenPool interface {
	GetToken(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(*bind.CallOpts) (string, error)
}

// TokenPoolMetadata defines the token pool version version and symbol of the corresponding token.
type TokenPoolMetadata struct {
	Version semver.Version
	Symbol  TokenSymbol
}

// NewTokenPoolWithMetadata returns a token pool along with its metadata.
func NewTokenPoolWithMetadata[P tokenPool](
	ctx context.Context,
	newTokenPool func(address common.Address, backend bind.ContractBackend) (P, error),
	poolAddress common.Address,
	chainClient cldf_evm.OnchainClient,
) (P, TokenPoolMetadata, error) {
	pool, err := newTokenPool(poolAddress, chainClient)
	if err != nil {
		return pool, TokenPoolMetadata{}, fmt.Errorf("failed to connect address %s with token pool bindings: %w", poolAddress, err)
	}
	tokenAddress, err := pool.GetToken(&bind.CallOpts{Context: ctx})
	if err != nil {
		return pool, TokenPoolMetadata{}, fmt.Errorf("failed to get token address from pool with address %s: %w", poolAddress, err)
	}
	typeAndVersionStr, err := pool.TypeAndVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return pool, TokenPoolMetadata{}, fmt.Errorf("failed to get type and version from pool with address %s: %w", poolAddress, err)
	}
	_, versionStr, err := ccipconfig.ParseTypeAndVersion(typeAndVersionStr)
	if err != nil {
		return pool, TokenPoolMetadata{}, fmt.Errorf("failed to parse type and version of pool with address %s: %w", poolAddress, err)
	}
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return pool, TokenPoolMetadata{}, fmt.Errorf("failed parsing version %s of pool with address %s: %w", versionStr, poolAddress, err)
	}
	token, err := erc20.NewERC20(tokenAddress, chainClient)
	if err != nil {
		return pool, TokenPoolMetadata{}, fmt.Errorf("failed to connect address %s with ERC20 bindings: %w", tokenAddress, err)
	}
	symbol, err := token.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return pool, TokenPoolMetadata{}, fmt.Errorf("failed to fetch symbol from token with address %s: %w", tokenAddress, err)
	}
	return pool, TokenPoolMetadata{
		Symbol:  TokenSymbol(symbol),
		Version: *version,
	}, nil
}
