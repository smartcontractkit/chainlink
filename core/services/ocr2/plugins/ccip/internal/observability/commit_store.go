package observability

import (
	"context"
	"time"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

type ObservedCommitStoreReader struct {
	ccipdata.CommitStoreReader
	metric metricDetails
}

func NewObservedCommitStoreReader(origin ccipdata.CommitStoreReader, chainID int64, pluginName string) *ObservedCommitStoreReader {
	return &ObservedCommitStoreReader{
		CommitStoreReader: origin,
		metric: metricDetails{
			interactionDuration: readerHistogram,
			resultSetSize:       readerDatasetSize,
			pluginName:          pluginName,
			readerName:          "CommitStoreReader",
			chainId:             chainID,
		},
	}
}

func (o *ObservedCommitStoreReader) GetExpectedNextSequenceNumber(context context.Context) (uint64, error) {
	return withObservedInteraction(o.metric, "GetExpectedNextSequenceNumber", func() (uint64, error) {
		return o.CommitStoreReader.GetExpectedNextSequenceNumber(context)
	})
}

func (o *ObservedCommitStoreReader) GetLatestPriceEpochAndRound(context context.Context) (uint64, error) {
	return withObservedInteraction(o.metric, "GetLatestPriceEpochAndRound", func() (uint64, error) {
		return o.CommitStoreReader.GetLatestPriceEpochAndRound(context)
	})
}

func (o *ObservedCommitStoreReader) GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confs int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return withObservedInteractionAndResults(o.metric, "GetCommitReportMatchingSeqNum", func() ([]cciptypes.CommitStoreReportWithTxMeta, error) {
		return o.CommitStoreReader.GetCommitReportMatchingSeqNum(ctx, seqNum, confs)
	})
}

func (o *ObservedCommitStoreReader) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confs int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return withObservedInteractionAndResults(o.metric, "GetAcceptedCommitReportsGteTimestamp", func() ([]cciptypes.CommitStoreReportWithTxMeta, error) {
		return o.CommitStoreReader.GetAcceptedCommitReportsGteTimestamp(ctx, ts, confs)
	})
}

func (o *ObservedCommitStoreReader) IsDown(ctx context.Context) (bool, error) {
	return withObservedInteraction(o.metric, "IsDown", func() (bool, error) {
		return o.CommitStoreReader.IsDown(ctx)
	})
}

func (o *ObservedCommitStoreReader) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return withObservedInteraction(o.metric, "IsBlessed", func() (bool, error) {
		return o.CommitStoreReader.IsBlessed(ctx, root)
	})
}

func (o *ObservedCommitStoreReader) VerifyExecutionReport(ctx context.Context, report cciptypes.ExecReport) (bool, error) {
	return withObservedInteraction(o.metric, "VerifyExecutionReport", func() (bool, error) {
		return o.CommitStoreReader.VerifyExecutionReport(ctx, report)
	})
}

func (o *ObservedCommitStoreReader) GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error) {
	return withObservedInteraction(o.metric, "GetCommitStoreStaticConfig", func() (cciptypes.CommitStoreStaticConfig, error) {
		return o.CommitStoreReader.GetCommitStoreStaticConfig(ctx)
	})
}
