package shared

import (
	"maps"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/aggregator_v3_interface"
)

type TokenRegistry interface {
	GetSymbol(desc string) (TokenSymbol, bool)
}

// Default implementation
type defaultRegistry struct{}

var registry TokenRegistry = defaultRegistry{}

// SetRegistry sets the global token registry implementation.
func SetRegistry(r TokenRegistry) {
	registry = r
}

// GetSymbolFromDescription retrieves the TokenSymbol associated with the given description.
// Delegates the lookup to the current registry implementation.
func GetSymbolFromDescription(desc string) (TokenSymbol, bool) {
	return registry.GetSymbol(desc)
}

// GetSymbol implements the default registry's lookup logic.
// It returns the TokenSymbol corresponding to a description, if found.
func (defaultRegistry) GetSymbol(desc string) (TokenSymbol, bool) {
	symbol, ok := DescriptionToTokenSymbol[desc]
	return symbol, ok
}

// NewMergedRegistry combines the defaultPriceFeed with new priceFeeds retrieved from CLD or elsewhere.
func NewMergedRegistry(tokens map[string]TokenSymbol) TokenRegistry {
	combined := make(map[string]TokenSymbol)

	// Add core defaults from Chainlink
	maps.Copy(combined, DescriptionToTokenSymbol)

	// Override or extend with CLD-provided tokens
	maps.Copy(combined, tokens)

	return mergedRegistry{entries: combined}
}

// mergedRegistry is a local wrapper
type mergedRegistry struct {
	entries map[string]TokenSymbol
}

func (r mergedRegistry) GetSymbol(desc string) (TokenSymbol, bool) {
	sym, ok := r.entries[desc]
	return sym, ok
}

type TokenSymbol string

func (ts TokenSymbol) String() string {
	return string(ts)
}

const (
	LinkSymbol                 TokenSymbol = "LINK"
	WethSymbol                 TokenSymbol = "WETH"
	WAVAXSymbol                TokenSymbol = "WAVAX"
	WBNBSymbol                 TokenSymbol = "WBNB"
	WPOLSymbol                 TokenSymbol = "WPOL"
	WSSymbol                   TokenSymbol = "WS"
	USDCSymbol                 TokenSymbol = "USDC"
	LBTCSymbol                 TokenSymbol = "LBTC"
	FactoryBurnMintERC20Symbol TokenSymbol = "Factory-BnM-ERC20"
	CCIPBnMSymbol              TokenSymbol = "CCIP-BnM"
	CCIPLnMSymbol              TokenSymbol = "CCIP-LnM"
	CLCCIPLnMSymbol            TokenSymbol = "clCCIP-LnM"
	USDCName                   string      = "USD Coin"
	LinkDecimals                           = 18
	WethDecimals                           = 18
	UsdcDecimals                           = 6

	// Price Feed Descriptions
	AvaxUSD  = "AVAX / USD"
	LinkUSD  = "LINK / USD"
	EthUSD   = "ETH / USD"
	MaticUSD = "MATIC / USD"
	BNBUSD   = "BNB / USD"
	FTMUSD   = "FTM / USD" // S token uses FTM / USD price feed under the hood

	// MockLinkAggregatorDescription is the description of the MockV3Aggregator.sol contract
	// https://github.com/smartcontractkit/chainlink/blob/a348b98e90527520049c580000a86fb8ceff7fa7/contracts/src/v0.8/tests/MockV3Aggregator.sol#L76-L76
	MockLinkAggregatorDescription = "v0.8/tests/MockV3Aggregator.sol"
	// MockWETHAggregatorDescription is the description from MockETHUSDAggregator.sol
	// https://github.com/smartcontractkit/chainlink/blob/a348b98e90527520049c580000a86fb8ceff7fa7/contracts/src/v0.8/automation/testhelpers/MockETHUSDAggregator.sol#L19-L19
	MockWETHAggregatorDescription = "MockETHUSDAggregator"
)

var (
	MockLinkPrice = deployment.E18Mult(500)
	MockWethPrice = big.NewInt(9e8)
	// DescriptionToTokenSymbol maps price feed description to token descriptor
	DescriptionToTokenSymbol = map[string]TokenSymbol{
		MockLinkAggregatorDescription: LinkSymbol,
		MockWETHAggregatorDescription: WethSymbol,
		LinkUSD:                       LinkSymbol,
		AvaxUSD:                       WAVAXSymbol,
		EthUSD:                        WethSymbol,
		MaticUSD:                      WPOLSymbol,
		BNBUSD:                        WBNBSymbol,
		FTMUSD:                        WSSymbol,
	}
	MockSymbolToDescription = map[TokenSymbol]string{
		LinkSymbol: MockLinkAggregatorDescription,
		WethSymbol: MockWETHAggregatorDescription,
	}
	TestDeviationPPB = ccipocr3.NewBigIntFromInt64(1e9)

	TokenSymbolSubstitute = map[string]string{
		"wS": WSSymbol.String(),
	}
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

func NewTestTokenConfig(feeds map[TokenSymbol]*aggregator_v3_interface.AggregatorV3Interface) TokenConfig {
	tc := NewTokenConfig()
	tc.UpsertTokenInfo(LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: ccipocr3.UnknownEncodedAddress(feeds[LinkSymbol].Address().String()),
			Decimals:          LinkDecimals,
			DeviationPPB:      TestDeviationPPB,
		},
	)
	tc.UpsertTokenInfo(WethSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: ccipocr3.UnknownEncodedAddress(feeds[WethSymbol].Address().String()),
			Decimals:          WethDecimals,
			DeviationPPB:      TestDeviationPPB,
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
	linkTokenAddr,
	wethTokenAddr common.Address,
) map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo {
	tokenToAggregate := make(map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo)
	if _, ok := tc.TokenSymbolToInfo[LinkSymbol]; !ok {
		lggr.Debugw("Link aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping LinkToken to Link aggregator")
		acc := ccipocr3.UnknownEncodedAddress(linkTokenAddr.String())
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[LinkSymbol]
	}

	if _, ok := tc.TokenSymbolToInfo[WethSymbol]; !ok {
		lggr.Debugw("Weth aggregator not found, deploy without mapping link token")
	} else {
		lggr.Debugw("Mapping WethToken to Weth aggregator")
		acc := ccipocr3.UnknownEncodedAddress(wethTokenAddr.String())
		tokenToAggregate[acc] = tc.TokenSymbolToInfo[WethSymbol]
	}

	return tokenToAggregate
}

type TokenDetails interface {
	Address() common.Address
	Symbol(opts *bind.CallOpts) (string, error)
	Decimals(opts *bind.CallOpts) (uint8, error)
}
