package observability

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
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

func (o *ObservedOffRampReader) EncodeExecutionReport(report ccipdata.ExecReport) ([]byte, error) {
	return withObservedInteraction(o.metric, "EncodeExecutionReport", func() ([]byte, error) {
		return o.OffRampReader.EncodeExecutionReport(report)
	})
}

func (o *ObservedOffRampReader) DecodeExecutionReport(report []byte) (ccipdata.ExecReport, error) {
	return withObservedInteraction(o.metric, "DecodeExecutionReport", func() (ccipdata.ExecReport, error) {
		return o.OffRampReader.DecodeExecutionReport(report)
	})
}

func (o *ObservedOffRampReader) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, confs int) ([]ccipdata.Event[ccipdata.ExecutionStateChanged], error) {
	return withObservedInteraction(o.metric, "GetExecutionStateChangesBetweenSeqNums", func() ([]ccipdata.Event[ccipdata.ExecutionStateChanged], error) {
		return o.OffRampReader.GetExecutionStateChangesBetweenSeqNums(ctx, seqNumMin, seqNumMax, confs)
	})
}

func (o *ObservedOffRampReader) GetDestinationTokens(ctx context.Context) ([]common.Address, error) {
	return withObservedInteractionAndResults(o.metric, "GetDestinationTokens", func() ([]common.Address, error) {
		return o.OffRampReader.GetDestinationTokens(ctx)
	})
}

func (o *ObservedOffRampReader) GetSenderNonce(ctx context.Context, sender common.Address) (uint64, error) {
	return withObservedInteraction(o.metric, "GetSenderNonce", func() (uint64, error) {
		return o.OffRampReader.GetSenderNonce(ctx, sender)
	})
}

func (o *ObservedOffRampReader) CurrentRateLimiterState(ctx context.Context) (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
	return withObservedInteraction(o.metric, "CurrentRateLimiterState", func() (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
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

func (o *ObservedOffRampReader) GetDestinationTokenPools(ctx context.Context) (map[common.Address]common.Address, error) {
	return withObservedInteraction(o.metric, "GetDestinationTokenPools", func() (map[common.Address]common.Address, error) {
		return o.OffRampReader.GetDestinationTokenPools(ctx)
	})
}
