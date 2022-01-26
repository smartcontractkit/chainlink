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
	t.Run("should set correct labels and cleanup all labels associated with different transmitters", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := new(mocks.Metrics)
		metrics.Test(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		metrics.On("SetFeedContractMetadata",
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
			feedConfig.GetSymbol(),         // symbol
		)
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

		metrics.On("SetFeedContractLinkBalance",
			envelope1.LinkBalance,          // balance
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(envelope1.Transmitter), // oracleName
			string(envelope1.Transmitter), // sender
		)
		metrics.On("SetHeadTrackerCurrentHead",
			envelope1.BlockNumber,        // blockNumber
			chainConfig.GetNetworkName(), // networkName
			chainConfig.GetChainID(),     // chainID
			chainConfig.GetNetworkID(),   // networkID
		)
		metrics.On("SetOffchainAggregatorAnswers",
			envelope1.LatestAnswer,         // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("IncOffchainAggregatorAnswersTotal",
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("SetOffchainAggregatorAnswerStalled",
			mock.Anything,                  // isSet
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			envelope1.LatestAnswer,         // value
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			string(envelope1.Transmitter),  // sender
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		exporter.Export(ctx, envelope1)

		metrics.On("SetFeedContractLinkBalance",
			envelope2.LinkBalance,          // balance
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(envelope2.Transmitter), // oracleName
			string(envelope2.Transmitter), // sender
		)
		metrics.On("SetHeadTrackerCurrentHead",
			envelope2.BlockNumber,        // blockNumber
			chainConfig.GetNetworkName(), // networkName
			chainConfig.GetChainID(),     // chainID
			chainConfig.GetNetworkID(),   // networkID
		)
		metrics.On("SetOffchainAggregatorAnswers",
			envelope2.LatestAnswer,         // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("IncOffchainAggregatorAnswersTotal",
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("SetOffchainAggregatorAnswerStalled",
			mock.Anything,                  // isSet
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			envelope2.LatestAnswer,         // value
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			string(envelope2.Transmitter),  // sender
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		)
		exporter.Export(ctx, envelope2)

		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(envelope1.Transmitter),  // oracleName
			string(envelope1.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		)
		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(envelope2.Transmitter),  // oracleName
			string(envelope2.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		)
		exporter.Cleanup(ctx)
	})
}
