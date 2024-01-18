package observability

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

type ObservedOffRampReader struct {
	ccipdata.OffRampReader
	metric metricDetails
}

func NewObservedOffRampReader(origin ccipdata.OffRampReader, chainID int64, pluginName string) *ObservedOffRampReader {
	return &ObservedOffRampReader{
		OffRampReader: origin,
		metric: metricDetails{
			interactionDuration: readerHistogram,
			resultSetSize:       readerDatasetSize,
			pluginName:          pluginName,
			readerName:          "OffRampReader",
			chainId:             chainID,
		},
	}
}

func (o *ObservedOffRampReader) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]ccipdata.Event[ccipdata.ExecutionStateChanged], error) {
	return withObservedInteraction(o.metric, "GetExecutionStateChangesBetweenSeqNums", func() ([]ccipdata.Event[ccipdata.ExecutionStateChanged], error) {
		return o.OffRampReader.GetExecutionStateChangesBetweenSeqNums(ctx, seqNumMin, seqNumMax, confs)
	})
}

func (o *ObservedOffRampReader) GetSenderNonce(ctx context.Context, sender common.Address) (uint64, error) {
	return withObservedInteraction(o.metric, "GetSenderNonce", func() (uint64, error) {
		return o.OffRampReader.GetSenderNonce(ctx, sender)
	})
}

func (o *ObservedOffRampReader) CurrentRateLimiterState(ctx context.Context) (ccipdata.TokenBucketRateLimit, error) {
	return withObservedInteraction(o.metric, "CurrentRateLimiterState", func() (ccipdata.TokenBucketRateLimit, error) {
		return o.OffRampReader.CurrentRateLimiterState(ctx)
	})
}

func (o *ObservedOffRampReader) GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error) {
	return withObservedInteraction(o.metric, "GetExecutionState", func() (uint8, error) {
		return o.OffRampReader.GetExecutionState(ctx, sequenceNumber)
	})
}

func (o *ObservedOffRampReader) GetStaticConfig(ctx context.Context) (ccipdata.OffRampStaticConfig, error) {
	return withObservedInteraction(o.metric, "GetStaticConfig", func() (ccipdata.OffRampStaticConfig, error) {
		return o.OffRampReader.GetStaticConfig(ctx)
	})
}

func (o *ObservedOffRampReader) GetSourceToDestTokensMapping(ctx context.Context) (map[common.Address]common.Address, error) {
	return withObservedInteraction(o.metric, "GetSourceToDestTokensMapping", func() (map[common.Address]common.Address, error) {
		return o.OffRampReader.GetSourceToDestTokensMapping(ctx)
	})
}

func (o *ObservedOffRampReader) GetTokens(ctx context.Context) (ccipdata.OffRampTokens, error) {
	return withObservedInteraction(o.metric, "GetTokens", func() (ccipdata.OffRampTokens, error) {
		return o.OffRampReader.GetTokens(ctx)
	})
}
