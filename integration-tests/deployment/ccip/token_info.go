package ccipdeployment

import (
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// TokenConfig mapping between token Descriptor (e.g. LinkDescriptor, WETHDescriptor)
// and the respective token info.
type TokenConfig struct {
	TokenDescriptorToInfo map[TokenDescriptor]pluginconfig.TokenInfo
}

func NewTokenConfig() TokenConfig {
	return TokenConfig{
		TokenDescriptorToInfo: make(map[TokenDescriptor]pluginconfig.TokenInfo),
	}
}

func DefaultTokenConfig() TokenConfig {
	descriptorToInfo := make(map[TokenDescriptor]pluginconfig.TokenInfo)
	// Add only enabled aggregates
	for _, descriptor := range EnabledTokensDescriptors {
		descriptorToInfo[descriptor] = TokenDescriptorToTokenInfo[descriptor]
	}
	return TokenConfig{
		TokenDescriptorToInfo: descriptorToInfo,
	}
}

func (tc *TokenConfig) UpsertTokenInfo(
	descriptor TokenDescriptor,
	info pluginconfig.TokenInfo,
) {
	tc.TokenDescriptorToInfo[descriptor] = info
}

// GetTokenInfo Adds mapping between dest chain tokens and their respective aggregators on feed chain.
func (tc *TokenConfig) GetTokenInfo(
	lggr logger.Logger,
	destState CCIPChainState,
) map[ocrtypes.Account]pluginconfig.TokenInfo {
	tokenToAggregate := make(map[ocrtypes.Account]pluginconfig.TokenInfo)
	if _, ok := tc.TokenDescriptorToInfo[LinkDescriptor]; !ok {
		lggr.Debugw("Link aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping LinkToken to Link aggregator")
		acc := ocrtypes.Account(destState.LinkToken.Address().String())
		tokenToAggregate[acc] = tc.TokenDescriptorToInfo[LinkDescriptor]
	}

	// TODO: Populate tokenInfo with weth and the token map in destState

	return tokenToAggregate
}

// These will be used for production values
var (
	LinkInfo = pluginconfig.TokenInfo{
		// Add real linkToken info
		AggregatorAddress: "", // Usually this will be already deployed on feed chain
		Decimals:          DECIMALS,
		DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
	}
	// Enable feeds that have proper info from on-chain
	EnabledTokensDescriptors   = []TokenDescriptor{}
	TokenDescriptorToTokenInfo = map[TokenDescriptor]pluginconfig.TokenInfo{
		LinkDescriptor: LinkInfo,
	}
)
