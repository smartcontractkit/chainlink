package pricegetter

import (
	"context"
	"math/big"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/parseutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

var _ PriceGetter = &PipelineGetter{}

type PipelineGetter struct {
	source        string
	runner        pipeline.Runner
	jobID         int32
	externalJobID uuid.UUID
	name          string
	lggr          logger.Logger
}

func NewPipelineGetter(source string, runner pipeline.Runner, jobID int32, externalJobID uuid.UUID, name string, lggr logger.Logger) (*PipelineGetter, error) {
	_, err := pipeline.Parse(source)
	if err != nil {
		return nil, err
	}

	return &PipelineGetter{
		source:        source,
		runner:        runner,
		jobID:         jobID,
		externalJobID: externalJobID,
		name:          name,
		lggr:          lggr,
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

func (d *PipelineGetter) GetJobSpecTokenPricesUSD(ctx context.Context) (map[cciptypes.Address]*big.Int, error) {
	prices, err := d.getPricesFromRunner(ctx)
	if err != nil {
		return nil, err
	}

	tokenPrices := make(map[cciptypes.Address]*big.Int)
	for tokenAddressStr, rawPrice := range prices {
		tokenAddressStr := ccipcalc.HexToAddress(tokenAddressStr)
		castedPrice, err := parseutil.ParseBigIntFromAny(rawPrice)
		if err != nil {
			return nil, err
		}

		tokenPrices[tokenAddressStr] = castedPrice
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
