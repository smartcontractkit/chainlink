package ccipdeployment

import (
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
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

func DefaultTokenConfig() TokenConfig {
	symbolToInfo := make(map[TokenSymbol]pluginconfig.TokenInfo)
	// Add only enabled aggregates
	for _, symbol := range EnabledTokensSymbols {
		symbolToInfo[symbol] = TokenSymbolToTokenInfo[symbol]
	}
	return TokenConfig{
		TokenSymbolToInfo: symbolToInfo,
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
	destState CCIPChainState,
) map[ocrtypes.Account]pluginconfig.TokenInfo {
	tokenToAggregate := make(map[ocrtypes.Account]pluginconfig.TokenInfo)
	if _, ok := tc.TokenSymbolToInfo[LinkSymbol]; !ok {
		lggr.Debugw("Link aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping LinkToken to Link aggregator")
		acc := ocrtypes.Account(destState.LinkToken.Address().String())
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[LinkSymbol]
	}

	// TODO: Populate tokenInfo with weth and the token map in destState

	return tokenToAggregate
}

// These will be used for production values
var (
	LinkInfo = pluginconfig.TokenInfo{
		// Add real linkToken info
		AggregatorAddress: "", // Usually this will be already deployed on feed chain
		Decimals:          LinkDecimals,
		DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
	}
	// Enable feeds that have proper info from on-chain
	EnabledTokensSymbols   = []TokenSymbol{}
	TokenSymbolToTokenInfo = map[TokenSymbol]pluginconfig.TokenInfo{
		LinkSymbol: LinkInfo,
	}
)
