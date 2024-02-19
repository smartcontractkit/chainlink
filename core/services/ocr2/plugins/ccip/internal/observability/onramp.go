package observability

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

type ObservedOnRampReader struct {
	ccipdata.OnRampReader
	metric metricDetails
}

func NewObservedOnRampReader(origin ccipdata.OnRampReader, chainID int64, pluginName string) *ObservedOnRampReader {
	return &ObservedOnRampReader{
		OnRampReader: origin,
		metric: metricDetails{
			interactionDuration: readerHistogram,
			resultSetSize:       readerDatasetSize,
			pluginName:          pluginName,
			readerName:          "OnRampReader",
			chainId:             chainID,
		},
	}
}

func (o ObservedOnRampReader) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, finalized bool) ([]cciptypes.EVM2EVMMessageWithTxMeta, error) {
	return withObservedInteractionAndResults(o.metric, "GetSendRequestsBetweenSeqNums", func() ([]cciptypes.EVM2EVMMessageWithTxMeta, error) {
		return o.OnRampReader.GetSendRequestsBetweenSeqNums(ctx, seqNumMin, seqNumMax, finalized)
	})
}

func (o ObservedOnRampReader) RouterAddress() (cciptypes.Address, error) {
	return withObservedInteraction(o.metric, "RouterAddress", func() (cciptypes.Address, error) {
		return o.OnRampReader.RouterAddress()
	})
}

func (o ObservedOnRampReader) GetDynamicConfig() (cciptypes.OnRampDynamicConfig, error) {
	return withObservedInteraction(o.metric, "GetDynamicConfig", func() (cciptypes.OnRampDynamicConfig, error) {
		return o.OnRampReader.GetDynamicConfig()
	})
}
