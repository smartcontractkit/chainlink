package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/mocks"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPrometheusExporter(t *testing.T) {
	t.Run("should remove all labels associated with different transmitters on the same feed", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := new(mocks.Metrics)
		metrics.Test(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		metrics.On("SetFeedContractMetadata", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()
		exporter, err := factory.NewExporter(chainConfig, feedConfig)
		require.NoError(t, err)

		envelope1, err := generateEnvelope()
		require.NoError(t, err)
		envelope1.Transmitter = types.Account(hexutil.Encode([]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, uint8(1),
		}))
		envelope2, err := generateEnvelope()
		require.NoError(t, err)
		envelope2.Transmitter = types.Account(hexutil.Encode([]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, uint8(2),
		}))

		metrics.On("SetNodeMetadata", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetHeadTrackerCurrentHead", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetOffchainAggregatorAnswers", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("IncOffchainAggregatorAnswersTotal", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetOffchainAggregatorAnswerStalled", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		exporter.Export(ctx, envelope1)

		metrics.On("SetNodeMetadata", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetHeadTrackerCurrentHead", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetOffchainAggregatorAnswers", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("IncOffchainAggregatorAnswersTotal", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetOffchainAggregatorAnswerStalled", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		exporter.Export(ctx, envelope2)

		metrics.On("Cleanup", mock.Anything, mock.Anything, mock.Anything, string(envelope1.Transmitter), string(envelope1.Transmitter), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		metrics.On("Cleanup", mock.Anything, mock.Anything, mock.Anything, string(envelope2.Transmitter), string(envelope2.Transmitter), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		exporter.Cleanup(ctx)
	})
}
