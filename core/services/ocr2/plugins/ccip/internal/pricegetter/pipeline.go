package pricegetter

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/parseutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

var _ PriceGetter = &PipelineGetter{}

// PipelineGetter is not supposed to be used but it seems that some JobSpecs are still using it.
// Should be removed after all JobSpecs migrate to the dynamic price getter. It uses a legacy pipeline component of
// chainlink core.
type PipelineGetter struct {
	source                string
	runner                pipeline.Runner
	jobID                 int32
	externalJobID         uuid.UUID
	name                  string
	lggr                  logger.Logger
	sourceNativeTokenAddr cciptypes.Address
	sourceChainSelector   uint64
	destChainSelector     uint64
}

func NewPipelineGetter(
	source string,
	runner pipeline.Runner,
	jobID int32,
	externalJobID uuid.UUID,
	name string,
	lggr logger.Logger,
	sourceNativeTokenAddr cciptypes.Address,
	sourceChainSelector uint64,
	destChainSelector uint64,
) (*PipelineGetter, error) {
	_, err := pipeline.Parse(source)
	if err != nil {
		return nil, err
	}

	return &PipelineGetter{
		source:                source,
		runner:                runner,
		jobID:                 jobID,
		externalJobID:         externalJobID,
		name:                  name,
		lggr:                  lggr,
		sourceNativeTokenAddr: sourceNativeTokenAddr,
		sourceChainSelector:   sourceChainSelector,
		destChainSelector:     destChainSelector,
	}, nil
}

// FilterForConfiguredTokens implements the PriceGetter interface.
// It filters a list of token addresses for only those that have a pipeline job configured on the TokenPricesUSDPipeline
func (d *PipelineGetter) FilterConfiguredTokens(ctx context.Context, tokens []cciptypes.Address) (configured []cciptypes.Address, unconfigured []cciptypes.Address, err error) {
	lcSource := strings.ToLower(d.source)
	for _, tk := range tokens {
		lcToken := strings.ToLower(string(tk))
		if strings.Contains(lcSource, lcToken) {
			configured = append(configured, tk)
		} else {
			unconfigured = append(unconfigured, tk)
		}
	}
	return configured, unconfigured, nil
}

func (d *PipelineGetter) GetJobSpecTokenPricesUSD(ctx context.Context) (map[ccipcommon.TokenID]*big.Int, error) {
	prices, err := d.getPricesFromRunner(ctx)
	if err != nil {
		return nil, err
	}

	// if the token address equals source native token then chain selector is source else it is dest

	tokenPrices := make(map[ccipcommon.TokenID]*big.Int)
	for tokenAddressStr, rawPrice := range prices {
		tokenAddr := ccipcalc.HexToAddress(tokenAddressStr)
		castedPrice, err1 := parseutil.ParseBigIntFromAny(rawPrice)
		if err1 != nil {
			return nil, fmt.Errorf("failed to parse price %s for token %s: %w", rawPrice, tokenAddr, err1)
		}

		tokenID := ccipcommon.TokenID{TokenAddress: tokenAddr, ChainSelector: d.destChainSelector}
		if tokenAddr == d.sourceNativeTokenAddr {
			tokenID.ChainSelector = d.sourceChainSelector
		}

		tokenPrices[tokenID] = castedPrice
	}

	return tokenPrices, nil
}

func (d *PipelineGetter) TokenPricesUSD(ctx context.Context, tokens []cciptypes.Address) (map[cciptypes.Address]*big.Int, error) {
	prices, err := d.getPricesFromRunner(ctx)
	if err != nil {
		return nil, err
	}

	providedTokensSet := mapset.NewSet(tokens...)
	tokenPrices := make(map[cciptypes.Address]*big.Int)
	for tokenAddressStr, rawPrice := range prices {
		tokenAddressStr := ccipcalc.HexToAddress(tokenAddressStr)
		castedPrice, err := parseutil.ParseBigIntFromAny(rawPrice)
		if err != nil {
			return nil, err
		}

		if providedTokensSet.Contains(tokenAddressStr) {
			tokenPrices[tokenAddressStr] = castedPrice
		}
	}

	// The mapping of token address to source of token price has to live offchain.
	// Best we can do is sanity check that the token price spec covers all our desired execution token prices.
	for _, token := range tokens {
		if _, ok := tokenPrices[token]; !ok {
			return nil, errors.Errorf("missing token %s from tokensForFeeCoin spec, got %v", token, prices)
		}
	}

	return tokenPrices, nil
}

func (d *PipelineGetter) GetTokenPricesUSD(ctx context.Context, tokens []ccipcommon.TokenID) (map[ccipcommon.TokenID]*big.Int, error) {
	tokenAddresses := make([]cciptypes.Address, len(tokens))
	for i, token := range tokens {
		tokenAddresses[i] = token.TokenAddress
	}

	tokenPrices, err := d.TokenPricesUSD(ctx, tokenAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get token prices %v: %w", tokenAddresses, err)
	}

	tokenPricesMap := make(map[ccipcommon.TokenID]*big.Int)
	for tokenAddr, price := range tokenPrices {
		tokenID := ccipcommon.TokenID{TokenAddress: tokenAddr, ChainSelector: d.destChainSelector}
		if tokenAddr == d.sourceNativeTokenAddr {
			tokenID.ChainSelector = d.sourceChainSelector
		}
		tokenPricesMap[tokenID] = price
	}

	return tokenPricesMap, nil
}

func (d *PipelineGetter) getPricesFromRunner(ctx context.Context) (map[string]interface{}, error) {
	_, trrs, err := d.runner.ExecuteRun(ctx, pipeline.Spec{
		ID:           d.jobID,
		DotDagSource: d.source,
		CreatedAt:    time.Now(),
		JobID:        d.jobID,
		JobName:      d.name,
		JobType:      "",
	}, pipeline.NewVarsFrom(map[string]interface{}{}))
	if err != nil {
		return nil, err
	}
	finalResult := trrs.FinalResult()
	if finalResult.HasErrors() {
		return nil, errors.Errorf("error getting prices %v", finalResult.AllErrors)
	}
	if len(finalResult.Values) != 1 {
		return nil, errors.Errorf("invalid number of price results, expected 1 got %v", len(finalResult.Values))
	}
	prices, ok := finalResult.Values[0].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("expected map output of price pipeline, got %T", finalResult.Values[0])
	}

	return prices, nil
}

func (d *PipelineGetter) Close() error {
	return d.runner.Close()
}
