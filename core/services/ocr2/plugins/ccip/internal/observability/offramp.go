package observability

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

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

func (o *ObservedOffRampReader) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]cciptypes.ExecutionStateChangedWithTxMeta, error) {
	return withObservedInteraction(o.metric, "GetExecutionStateChangesBetweenSeqNums", func() ([]cciptypes.ExecutionStateChangedWithTxMeta, error) {
		return o.OffRampReader.GetExecutionStateChangesBetweenSeqNums(ctx, seqNumMin, seqNumMax, confs)
	})
}

func (o *ObservedOffRampReader) CurrentRateLimiterState(ctx context.Context) (cciptypes.TokenBucketRateLimit, error) {
	return withObservedInteraction(o.metric, "CurrentRateLimiterState", func() (cciptypes.TokenBucketRateLimit, error) {
		return o.OffRampReader.CurrentRateLimiterState(ctx)
	})
}

func (o *ObservedOffRampReader) GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error) {
	return withObservedInteraction(o.metric, "GetExecutionState", func() (uint8, error) {
		return o.OffRampReader.GetExecutionState(ctx, sequenceNumber)
	})
}

func (o *ObservedOffRampReader) GetStaticConfig(ctx context.Context) (cciptypes.OffRampStaticConfig, error) {
	return withObservedInteraction(o.metric, "GetStaticConfig", func() (cciptypes.OffRampStaticConfig, error) {
		return o.OffRampReader.GetStaticConfig(ctx)
	})
}

func (o *ObservedOffRampReader) GetSourceToDestTokensMapping(ctx context.Context) (map[cciptypes.Address]cciptypes.Address, error) {
	return withObservedInteraction(o.metric, "GetSourceToDestTokensMapping", func() (map[cciptypes.Address]cciptypes.Address, error) {
		return o.OffRampReader.GetSourceToDestTokensMapping(ctx)
	})
}

func (o *ObservedOffRampReader) GetTokens(ctx context.Context) (cciptypes.OffRampTokens, error) {
	return withObservedInteraction(o.metric, "GetTokens", func() (cciptypes.OffRampTokens, error) {
		return o.OffRampReader.GetTokens(ctx)
	})
}

func (o *ObservedOffRampReader) GetSendersNonce(ctx context.Context, senders []cciptypes.Address) (map[cciptypes.Address]uint64, error) {
	return withObservedInteraction(o.metric, "ListSenderNonces", func() (map[cciptypes.Address]uint64, error) {
		return o.OffRampReader.ListSenderNonces(ctx, senders)
	})
}
