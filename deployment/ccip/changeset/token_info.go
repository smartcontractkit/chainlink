package changeset

import (
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type TokenSymbol string

const (
	LinkSymbol   TokenSymbol = "LINK"
	WethSymbol   TokenSymbol = "WETH"
	USDCSymbol   TokenSymbol = "USDC"
	USDCName     string      = "USD Coin"
	LinkDecimals             = 18
	WethDecimals             = 18
	UsdcDecimals             = 6
)

var (
	TestDeviationPPB = ccipocr3.NewBigIntFromInt64(1e9)
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

func NewTestTokenConfig(linkSymbolAddress, wethSymbolAddress string, chainSelector uint64) TokenConfig {
	tc := NewTokenConfig()
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return tc
	}
	tc.UpsertTokenInfo(LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: ccipocr3.UnknownEncodedAddress(linkSymbolAddress),
			Decimals:          LinkDecimals,
			DeviationPPB:      TestDeviationPPB,
			ChainFamily:       family,
		},
	)
	tc.UpsertTokenInfo(WethSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: ccipocr3.UnknownEncodedAddress(wethSymbolAddress),
			Decimals:          WethDecimals,
			DeviationPPB:      TestDeviationPPB,
			ChainFamily:       family,
		},
	)
	return tc
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
	linkTokenAddress string,
	wethTokenAddress string,
) map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo {
	tokenToAggregate := make(map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo)
	if _, ok := tc.TokenSymbolToInfo[LinkSymbol]; !ok {
		lggr.Debugw("Link aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping LinkToken to Link aggregator")
		acc := ccipocr3.UnknownEncodedAddress(linkTokenAddress)
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[LinkSymbol]
	}

	if _, ok := tc.TokenSymbolToInfo[WethSymbol]; !ok {
		lggr.Debugw("Weth aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping WethToken to Weth aggregator")
		acc := ccipocr3.UnknownEncodedAddress(wethTokenAddress)
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[WethSymbol]
	}

	return tokenToAggregate
}
