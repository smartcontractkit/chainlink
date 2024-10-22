package ccipdeployment

import (
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// TokenConfig mapping between token Symbol (e.g. LinkSymbol, WethSymbol)
// and the respective token info.
type TokenConfig struct {
	TokenSymbolToInfo map[TokenSymbol]pluginconfig.TokenInfo
}

func NewTokenConfig() TokenConfig {
	return TokenConfig{
		TokenSymbolToInfo: make(map[TokenSymbol]pluginconfig.TokenInfo),
	}
}

func (tc *TokenConfig) UpsertTokenInfo(
	symbol TokenSymbol,
	info pluginconfig.TokenInfo,
) {
	tc.TokenSymbolToInfo[symbol] = info
}

// GetTokenInfo Adds mapping between dest chain tokens and their respective aggregators on feed chain.
func (tc *TokenConfig) GetTokenInfo(
	lggr logger.Logger,
	linkToken *burn_mint_erc677.BurnMintERC677,
	wethToken *weth9.WETH9,
) map[ocrtypes.Account]pluginconfig.TokenInfo {
	tokenToAggregate := make(map[ocrtypes.Account]pluginconfig.TokenInfo)
	if _, ok := tc.TokenSymbolToInfo[LinkSymbol]; !ok {
		lggr.Debugw("Link aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping LinkToken to Link aggregator")
		acc := ocrtypes.Account(linkToken.Address().String())
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[LinkSymbol]
	}

	if _, ok := tc.TokenSymbolToInfo[WethSymbol]; !ok {
		lggr.Debugw("Weth aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping WethToken to Weth aggregator")
		acc := ocrtypes.Account(wethToken.Address().String())
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[WethSymbol]
	}

	return tokenToAggregate
}
