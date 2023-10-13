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
			histogram:  offRampHistogram,
			pluginName: pluginName,
			chainId:    chainID,
		},
	}
}

func (o *ObservedOffRampReader) EncodeExecutionReport(report ccipdata.ExecReport) ([]byte, error) {
	return withObservedContract(o.metric, "EncodeExecutionReport", func() ([]byte, error) {
		return o.OffRampReader.EncodeExecutionReport(report)
	})
}

func (o *ObservedOffRampReader) DecodeExecutionReport(report []byte) (ccipdata.ExecReport, error) {
	return withObservedContract(o.metric, "DecodeExecutionReport", func() (ccipdata.ExecReport, error) {
		return o.OffRampReader.DecodeExecutionReport(report)
	})
}

func (o *ObservedOffRampReader) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]ccipdata.Event[ccipdata.ExecutionStateChanged], error) {
	return withObservedContract(o.metric, "GetExecutionStateChangesBetweenSeqNums", func() ([]ccipdata.Event[ccipdata.ExecutionStateChanged], error) {
		return o.OffRampReader.GetExecutionStateChangesBetweenSeqNums(ctx, seqNumMin, seqNumMax, confs)
	})
}

func (o *ObservedOffRampReader) GetDestinationTokens(ctx context.Context) ([]common.Address, error) {
	return withObservedContract(o.metric, "GetDestinationTokens", func() ([]common.Address, error) {
		return o.OffRampReader.GetDestinationTokens(ctx)
	})
}

func (o *ObservedOffRampReader) GetPoolByDestToken(ctx context.Context, address common.Address) (common.Address, error) {
	return withObservedContract(o.metric, "GetPoolByDestToken", func() (common.Address, error) {
		return o.OffRampReader.GetPoolByDestToken(ctx, address)
	})
}

func (o *ObservedOffRampReader) GetDestinationToken(ctx context.Context, address common.Address) (common.Address, error) {
	return withObservedContract(o.metric, "GetDestinationToken", func() (common.Address, error) {
		return o.OffRampReader.GetDestinationToken(ctx, address)
	})
}

func (o *ObservedOffRampReader) GetSupportedTokens(ctx context.Context) ([]common.Address, error) {
	return withObservedContract(o.metric, "GetSupportedTokens", func() ([]common.Address, error) {
		return o.OffRampReader.GetSupportedTokens(ctx)
	})
}
