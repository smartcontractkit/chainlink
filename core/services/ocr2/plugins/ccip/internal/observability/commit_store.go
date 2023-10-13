package observability

import (
	"time"

	"golang.org/x/net/context"

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
			histogram:  commitStoreHistogram,
			pluginName: pluginName,
			chainId:    chainID,
		},
	}
}

func (o *ObservedCommitStoreReader) GetExpectedNextSequenceNumber(context context.Context) (uint64, error) {
	return withObservedContract(o.metric, "GetExpectedNextSequenceNumber", func() (uint64, error) {
		return o.CommitStoreReader.GetExpectedNextSequenceNumber(context)
	})
}

func (o *ObservedCommitStoreReader) GetLatestPriceEpochAndRound(context context.Context) (uint64, error) {
	return withObservedContract(o.metric, "GetLatestPriceEpochAndRound", func() (uint64, error) {
		return o.CommitStoreReader.GetLatestPriceEpochAndRound(context)
	})
}

func (o *ObservedCommitStoreReader) GetAcceptedCommitReportsGteSeqNum(ctx context.Context, seqNum uint64, confs int) ([]ccipdata.Event[ccipdata.CommitStoreReport], error) {
	return withObservedContract(o.metric, "GetAcceptedCommitReportsGteSeqNum", func() ([]ccipdata.Event[ccipdata.CommitStoreReport], error) {
		return o.CommitStoreReader.GetAcceptedCommitReportsGteSeqNum(ctx, seqNum, confs)
	})
}

func (o *ObservedCommitStoreReader) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confs int) ([]ccipdata.Event[ccipdata.CommitStoreReport], error) {
	return withObservedContract(o.metric, "GetAcceptedCommitReportsGteTimestamp", func() ([]ccipdata.Event[ccipdata.CommitStoreReport], error) {
		return o.CommitStoreReader.GetAcceptedCommitReportsGteTimestamp(ctx, ts, confs)
	})
}

func (o *ObservedCommitStoreReader) IsDown(ctx context.Context) (bool, error) {
	return withObservedContract(o.metric, "IsDown", func() (bool, error) {
		return o.CommitStoreReader.IsDown(ctx)
	})
}

func (o *ObservedCommitStoreReader) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return withObservedContract(o.metric, "IsBlessed", func() (bool, error) {
		return o.CommitStoreReader.IsBlessed(ctx, root)
	})
}

func (o *ObservedCommitStoreReader) EncodeCommitReport(report ccipdata.CommitStoreReport) ([]byte, error) {
	return withObservedContract(o.metric, "EncodeCommitReport", func() ([]byte, error) {
		return o.CommitStoreReader.EncodeCommitReport(report)
	})
}

func (o *ObservedCommitStoreReader) DecodeCommitReport(report []byte) (ccipdata.CommitStoreReport, error) {
	return withObservedContract(o.metric, "DecodeCommitReport", func() (ccipdata.CommitStoreReport, error) {
		return o.CommitStoreReader.DecodeCommitReport(report)
	})
}

func (o *ObservedCommitStoreReader) VerifyExecutionReport(ctx context.Context, report ccipdata.ExecReport) (bool, error) {
	return withObservedContract(o.metric, "VerifyExecutionReport", func() (bool, error) {
		return o.CommitStoreReader.VerifyExecutionReport(ctx, report)
	})
}
