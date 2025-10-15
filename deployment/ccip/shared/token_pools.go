package shared

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc20"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
)

var CurrentTokenPoolVersion = deployment.Version1_5_1
var FastTransferTokenPoolVersion = deployment.Version1_6_3Dev
var BurnMintWithExternalMinterFastTransferTokenPoolVersion = deployment.Version1_6_0
var HybridWithExternalMinterFastTransferTokenPoolVersion = deployment.Version1_6_0

var TokenTypes = map[cldf.ContractType]struct{}{
	BurnMintToken:      {},
	ERC20Token:         {},
	ERC677Token:        {},
	ERC677TokenHelper:  {},
	BurnMintERC20Token: {},
}

var TokenPoolTypes = map[cldf.ContractType]struct{}{
	BurnMintFastTransferTokenPool:                   {},
	BurnMintTokenPool:                               {},
	BurnWithFromMintTokenPool:                       {},
	BurnFromMintTokenPool:                           {},
	LockReleaseTokenPool:                            {},
	USDCTokenPool:                                   {},
	HybridLockReleaseUSDCTokenPool:                  {},
	BurnMintWithExternalMinterFastTransferTokenPool: {},
	HybridWithExternalMinterFastTransferTokenPool:   {},
	BurnMintWithExternalMinterTokenPool:             {},
	HybridWithExternalMinterTokenPool:               {},
}

var TokenPoolVersions = map[semver.Version]struct{}{
	deployment.Version1_5_0:      {},
	deployment.Version1_5_1:      {},
	FastTransferTokenPoolVersion: {},
	deployment.Version1_6_0:      {},
	deployment.Version1_6_2:      {},
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
		// fallback: try to normalize invalid semver like 1.6.x-dev -> 1.6.0-dev
		safeVersion := strings.ReplaceAll(versionStr, "x", "3")
		version, err = semver.NewVersion(safeVersion)
		if err != nil {
			return pool, TokenPoolMetadata{}, fmt.Errorf("failed parsing version %s (normalized as %s): %w", versionStr, safeVersion, err)
		}
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
