package observability

import (
	"context"
	"time"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

type ObservedPriceRegistryReader struct {
	ccipdata.PriceRegistryReader
	metric metricDetails
}

func NewPriceRegistryReader(origin ccipdata.PriceRegistryReader, chainID int64, pluginName string) *ObservedPriceRegistryReader {
	return &ObservedPriceRegistryReader{
		PriceRegistryReader: origin,
		metric: metricDetails{
			interactionDuration: readerHistogram,
			resultSetSize:       readerDatasetSize,
			pluginName:          pluginName,
			readerName:          "PriceRegistryReader",
			chainId:             chainID,
		},
	}
}

func (o *ObservedPriceRegistryReader) GetTokenPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]cciptypes.TokenPriceUpdateWithTxMeta, error) {
	return withObservedInteractionAndResults(o.metric, "GetTokenPriceUpdatesCreatedAfter", func() ([]cciptypes.TokenPriceUpdateWithTxMeta, error) {
		return o.PriceRegistryReader.GetTokenPriceUpdatesCreatedAfter(ctx, ts, confs)
	})
}

func (o *ObservedPriceRegistryReader) GetGasPriceUpdatesCreatedAfter(ctx context.Context, chainSelector uint64, ts time.Time, confs int) ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
	return withObservedInteractionAndResults(o.metric, "GetGasPriceUpdatesCreatedAfter", func() ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
		return o.PriceRegistryReader.GetGasPriceUpdatesCreatedAfter(ctx, chainSelector, ts, confs)
	})
}

func (o *ObservedPriceRegistryReader) GetAllGasPriceUpdatesCreatedAfter(ctx context.Context, ts time.Time, confs int) ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
	return withObservedInteractionAndResults(o.metric, "GetAllGasPriceUpdatesCreatedAfter", func() ([]cciptypes.GasPriceUpdateWithTxMeta, error) {
		return o.PriceRegistryReader.GetAllGasPriceUpdatesCreatedAfter(ctx, ts, confs)
	})
}

func (o *ObservedPriceRegistryReader) GetFeeTokens(ctx context.Context) ([]cciptypes.Address, error) {
	return withObservedInteraction(o.metric, "GetFeeTokens", func() ([]cciptypes.Address, error) {
		return o.PriceRegistryReader.GetFeeTokens(ctx)
	})
}

func (o *ObservedPriceRegistryReader) GetTokenPrices(ctx context.Context, wantedTokens []cciptypes.Address) ([]cciptypes.TokenPriceUpdate, error) {
	return withObservedInteractionAndResults(o.metric, "GetTokenPrices", func() ([]cciptypes.TokenPriceUpdate, error) {
		return o.PriceRegistryReader.GetTokenPrices(ctx, wantedTokens)
	})
}

func (o *ObservedPriceRegistryReader) GetTokensDecimals(ctx context.Context, tokenAddresses []cciptypes.Address) ([]uint8, error) {
	return withObservedInteractionAndResults(o.metric, "GetTokensDecimals", func() ([]uint8, error) {
		return o.PriceRegistryReader.GetTokensDecimals(ctx, tokenAddresses)
	})
}
