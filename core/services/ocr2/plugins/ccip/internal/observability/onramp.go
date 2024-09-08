package observability

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

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

func (o ObservedOnRampReader) RouterAddress(ctx context.Context) (cciptypes.Address, error) {
	return withObservedInteraction(o.metric, "RouterAddress", func() (cciptypes.Address, error) {
		return o.OnRampReader.RouterAddress(ctx)
	})
}

func (o ObservedOnRampReader) GetDynamicConfig(ctx context.Context) (cciptypes.OnRampDynamicConfig, error) {
	return withObservedInteraction(o.metric, "GetDynamicConfig", func() (cciptypes.OnRampDynamicConfig, error) {
		return o.OnRampReader.GetDynamicConfig(ctx)
	})
}

func (o ObservedOnRampReader) IsSourceCursed(ctx context.Context) (bool, error) {
	return withObservedInteraction(o.metric, "IsSourceCursed", func() (bool, error) {
		return o.OnRampReader.IsSourceCursed(ctx)
	})
}

func (o ObservedOnRampReader) IsSourceChainHealthy(ctx context.Context) (bool, error) {
	return withObservedInteraction(o.metric, "IsSourceChainHealthy", func() (bool, error) {
		return o.OnRampReader.IsSourceChainHealthy(ctx)
	})
}

func (o ObservedOnRampReader) SourcePriceRegistryAddress(ctx context.Context) (cciptypes.Address, error) {
	return withObservedInteraction(o.metric, "SourcePriceRegistryAddress", func() (cciptypes.Address, error) {
		return o.OnRampReader.SourcePriceRegistryAddress(ctx)
	})
}
