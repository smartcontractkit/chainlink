package dione

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/integration-tests/client"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	integrationtesthelpers "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/integration"
)

// NewCCIPJobSpecParams returns set of parameters needed for setting up ccip jobs for sourceClient --> destClient
func NewCCIPJobSpecParams(sourceChainConfig *rhea.EVMChainConfig, sourceLane rhea.EVMLaneConfig, destChainConfig *rhea.EVMChainConfig, destLane rhea.EVMLaneConfig, version string) integrationtesthelpers.CCIPJobSpecParams {
	return integrationtesthelpers.CCIPJobSpecParams{
		OffRamp:                destLane.OffRamp,
		CommitStore:            destLane.CommitStore,
		SourceChainName:        ccip.ChainName(int64(sourceChainConfig.EvmChainId)),
		DestChainName:          ccip.ChainName(int64(destChainConfig.EvmChainId)),
		TokenPricesUSDPipeline: GetTokenPricesUSDPipeline(getPipelineTokens(sourceChainConfig, destChainConfig)),
		SourceStartBlock:       sourceLane.DeploySettings.DeployedAtBlock,
		DestStartBlock:         destLane.DeploySettings.DeployedAtBlock,
		P2PV2Bootstrappers:     []string{}, // Set in env vars
		SourceEvmChainId:       sourceChainConfig.EvmChainId,
		DestEvmChainId:         destChainConfig.EvmChainId,
		Version:                version,
	}
}

// Gathers all tokens needed for TokenPricesUSDPipeline
func getPipelineTokens(sourceChainConfig *rhea.EVMChainConfig, destChainConfig *rhea.EVMChainConfig) []rhea.EVMBridgedToken {
	var pipelineTokens []rhea.EVMBridgedToken

	for _, token := range destChainConfig.SupportedTokens {
		token.ChainId = destChainConfig.EvmChainId
		pipelineTokens = append(pipelineTokens, token)
	}
	for _, feeTokenName := range destChainConfig.FeeTokens {
		if _, ok := destChainConfig.SupportedTokens[feeTokenName]; !ok {
			panic(fmt.Errorf("FeeToken is not a supported token for chain: %d", sourceChainConfig.EvmChainId))
		}
	}
	if sourceChainConfig.WrappedNative == "" {
		panic(fmt.Errorf("WrappedNative is nil for chain: %d", sourceChainConfig.EvmChainId))
	}
	if _, ok := sourceChainConfig.SupportedTokens[sourceChainConfig.WrappedNative]; !ok {
		panic(fmt.Errorf("SupportedTokens does not contain WrappedNative: %s on chain: %d", sourceChainConfig.WrappedNative, sourceChainConfig.EvmChainId))
	}
	sourceWrappedNative := sourceChainConfig.SupportedTokens[sourceChainConfig.WrappedNative]
	sourceWrappedNative.ChainId = sourceChainConfig.EvmChainId
	pipelineTokens = append(pipelineTokens, sourceWrappedNative)

	return pipelineTokens
}

func GetTokenPricesUSDPipeline(tokens []rhea.EVMBridgedToken) string {
	var tokenPricesUSDPipeline string
	for i, token := range tokens {
		if token.TokenPriceType == rhea.PriceFeeds {
			err := ValidatePriceFeedToken(token)
			if err != nil {
				panic(err)
			}
			tokenPricesUSDPipeline += fmt.Sprintf(`
encode_call_token%[1]d_usd  [type="ethabiencode" abi="latestRoundData()"]

call_token%[1]d_usd [type="ethcall"
evmChainId=%[2]d
contract="%s"
data="$(encode_call_token%[1]d_usd)"]

decode_result_token%[1]d_usd [type="ethabidecode"
abi="uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound"
data="$(call_token%[1]d_usd)"]

multiply_token%[1]d_usd [type="multiply" input="$(decode_result_token%[1]d_usd.answer)" times=%[4]d]

encode_call_token%[1]d_usd -> call_token%[1]d_usd -> decode_result_token%[1]d_usd -> multiply_token%[1]d_usd
`, i+1, token.ChainId, token.PriceFeed.Aggregator, token.PriceFeed.Multiplier)
		}
	}
	tokenPricesUSDPipeline += "merge [type=merge left=\"{}\" right=\"{"
	for i, token := range tokens {
		if token.TokenPriceType == rhea.PriceFeeds {
			tokenPricesUSDPipeline += fmt.Sprintf(`\\\"%s\\\":$(multiply_token%d_usd),`, token.Token, i+1)
		} else if token.TokenPriceType == rhea.TokenPrices || token.TokenPriceType == "" {
			// to be consistent with prices in real price feeds
			// job spec prices use price per token
			// committing don will translate this value into price per 1e18 token nits
			tokenPricesUSDPipeline += fmt.Sprintf(`\\\"%s\\\":\\\"%s\\\",`, token.Token.Hex(), token.Price)
		}
	}
	tokenPricesUSDPipeline = strings.TrimSuffix(tokenPricesUSDPipeline, ",")
	tokenPricesUSDPipeline += "}\"];"

	return tokenPricesUSDPipeline
}

func GetOCRkeysForChainType(OCRKeys client.OCR2Keys, chainType string) client.OCR2KeyData {
	for _, key := range OCRKeys.Data {
		if key.Attributes.ChainType == chainType {
			return key
		}
	}

	panic("Keys not found for chain")
}

func ValidatePriceFeedToken(token rhea.EVMBridgedToken) error {
	if token.PriceFeed.Aggregator == (common.Address{}) {
		return fmt.Errorf("must set PriceFeed Aggregator address for token: %s", token.Token.Hex())
	}
	if token.PriceFeed.Multiplier == nil || token.PriceFeed.Multiplier.BitLen() == 0 {
		return fmt.Errorf("must set PriceFeed Multiplier for token %s", token.Token.Hex())
	}
	if token.ChainId == 0 {
		return fmt.Errorf("invalid token chain id for token %s", token.Token.Hex())
	}
	return nil
}
