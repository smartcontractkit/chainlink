package ccipdeployment

import (
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
) map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo {
	tokenToAggregate := make(map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo)
	if _, ok := tc.TokenSymbolToInfo[LinkSymbol]; !ok {
		lggr.Debugw("Link aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping LinkToken to Link aggregator")
		acc := ccipocr3.UnknownEncodedAddress(linkToken.Address().String())
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[LinkSymbol]
	}

	// TODO: Populate tokenInfo with weth and the token map in destState

	return tokenToAggregate
}
