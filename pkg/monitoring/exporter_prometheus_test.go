package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPrometheusExporter(t *testing.T) {
	t.Run("should set correct labels and cleanup all labels associated with different transmitters", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := NewMetricsMock(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		nodes := []NodeConfig{generateNodeConfig(), generateNodeConfig()}
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
		).Once()
		exporter, err := factory.NewExporter(ExporterParams{chainConfig, feedConfig, nodes})
		require.NoError(t, err)

		envelope1, err := generateEnvelope()
		require.NoError(t, err)
		envelope1.Transmitter = nodes[0].GetAccount()
		envelope2, err := generateEnvelope()
		require.NoError(t, err)
		envelope2.Transmitter = nodes[1].GetAccount()

		humanizedAnswer1 := toFloat64(envelope1.LatestAnswer) / toFloat64(feedConfig.GetMultiply())
		humanizedAnswer2 := toFloat64(envelope2.LatestAnswer) / toFloat64(feedConfig.GetMultiply())
		humanizedJuelsPerFeeCoin1 := toFloat64(envelope1.JuelsPerFeeCoin) / toFloat64(feedConfig.GetMultiply())
		humanizedJuelsPerFeeCoin2 := toFloat64(envelope2.JuelsPerFeeCoin) / toFloat64(feedConfig.GetMultiply())

		metrics.On("SetFeedContractLinkBalance",
			toFloat64(envelope1.LinkBalance), // balance
			feedConfig.GetID(),               // contractAddress
			feedConfig.GetID(),               // feedID
			chainConfig.GetChainID(),         // chainID
			feedConfig.GetContractStatus(),   // contractStatus
			feedConfig.GetContractType(),     // contractType
			feedConfig.GetName(),             // feedName
			feedConfig.GetPath(),             // feedPath
			chainConfig.GetNetworkID(),       // networkID
			chainConfig.GetNetworkName(),     // networkName
		).Once()
		metrics.On("SetLinkAvailableForPayment",
			toFloat64(envelope1.LinkAvailableForPayment), //link available for payment
			feedConfig.GetID(),                           // feedID
			chainConfig.GetChainID(),                     // chainID
			feedConfig.GetContractStatus(),               // contractStatus
			feedConfig.GetContractType(),                 // contractType
			feedConfig.GetName(),                         // feedName
			feedConfig.GetPath(),                         // feedPath
			chainConfig.GetNetworkID(),                   // networkID
			chainConfig.GetNetworkName(),                 // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(nodes[0].GetName()),    // oracleName
			string(envelope1.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope1.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
		metrics.On("SetOffchainAggregatorAnswers",
			humanizedAnswer1,               // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorAnswersRaw",
			toFloat64(envelope1.LatestAnswer), // answer
			feedConfig.GetID(),                // contractAddress
			feedConfig.GetID(),                // feedID
			chainConfig.GetChainID(),          // chainID
			feedConfig.GetContractStatus(),    // contractStatus
			feedConfig.GetContractType(),      // contractType
			feedConfig.GetName(),              // feedName
			feedConfig.GetPath(),              // feedPath
			chainConfig.GetNetworkID(),        // networkID
			chainConfig.GetNetworkName(),      // networkName
		).Once()
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
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinRaw",
			toFloat64(envelope1.JuelsPerFeeCoin),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoin",
			humanizedJuelsPerFeeCoin1,
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
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
		).Once()
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			humanizedAnswer1,               // answer
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
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinReceivedValues",
			humanizedJuelsPerFeeCoin1,      //juels/feecoin
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
		).Once()
		metrics.On("SetOffchainAggregatorRoundID",
			float64(envelope1.AggregatorRoundID),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		exporter.Export(ctx, envelope1)

		metrics.On("SetFeedContractLinkBalance",
			toFloat64(envelope2.LinkBalance), // balance
			feedConfig.GetID(),               // contractAddress
			feedConfig.GetID(),               // feedID
			chainConfig.GetChainID(),         // chainID
			feedConfig.GetContractStatus(),   // contractStatus
			feedConfig.GetContractType(),     // contractType
			feedConfig.GetName(),             // feedName
			feedConfig.GetPath(),             // feedPath
			chainConfig.GetNetworkID(),       // networkID
			chainConfig.GetNetworkName(),     // networkName
		).Once()
		metrics.On("SetLinkAvailableForPayment",
			toFloat64(envelope2.LinkAvailableForPayment), //link available for payment
			feedConfig.GetID(),                           // feedID
			chainConfig.GetChainID(),                     // chainID
			feedConfig.GetContractStatus(),               // contractStatus
			feedConfig.GetContractType(),                 // contractType
			feedConfig.GetName(),                         // feedName
			feedConfig.GetPath(),                         // feedPath
			chainConfig.GetNetworkID(),                   // networkID
			chainConfig.GetNetworkName(),                 // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(nodes[1].GetName()),    // oracleName
			string(envelope2.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope2.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
		metrics.On("SetOffchainAggregatorAnswers",
			humanizedAnswer2,
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorAnswersRaw",
			toFloat64(envelope2.LatestAnswer), // answer
			feedConfig.GetID(),                // contractAddress
			feedConfig.GetID(),                // feedID
			chainConfig.GetChainID(),          // chainID
			feedConfig.GetContractStatus(),    // contractStatus
			feedConfig.GetContractType(),      // contractType
			feedConfig.GetName(),              // feedName
			feedConfig.GetPath(),              // feedPath
			chainConfig.GetNetworkID(),        // networkID
			chainConfig.GetNetworkName(),      // networkName
		).Once()
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
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinRaw",
			toFloat64(envelope2.JuelsPerFeeCoin),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoin",
			humanizedJuelsPerFeeCoin2,
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
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
		).Once()
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			humanizedAnswer2,
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
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinReceivedValues",
			humanizedJuelsPerFeeCoin2,
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
		).Once()
		metrics.On("SetOffchainAggregatorRoundID",
			float64(envelope2.AggregatorRoundID),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		exporter.Export(ctx, envelope2)

		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(nodes[0].GetName()),     // oracleName
			string(envelope1.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		).Once()
		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(nodes[1].GetName()),     // oracleName
			string(envelope2.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		).Once()
		exporter.Cleanup(ctx)
	})
	t.Run("should not emit metrics for stale transmissions", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := NewMetricsMock(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		nodes := []NodeConfig{generateNodeConfig()}
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
		).Once()
		exporter, err := factory.NewExporter(ExporterParams{chainConfig, feedConfig, nodes})
		require.NoError(t, err)

		envelope1, err := generateEnvelope()
		require.NoError(t, err)
		envelope2, err := generateEnvelope()
		require.NoError(t, err)
		envelope2.LatestAnswer = envelope1.LatestAnswer
		envelope2.LatestTimestamp = envelope1.LatestTimestamp
		envelope2.Transmitter = nodes[0].GetAccount()
		envelope1.Transmitter = nodes[0].GetAccount()

		humanizedAnswer := toFloat64(envelope1.LatestAnswer) / toFloat64(feedConfig.GetMultiply())
		humanizedJuelsPerFeeCoin := toFloat64(envelope1.JuelsPerFeeCoin) / toFloat64(feedConfig.GetMultiply())

		metrics.On("SetFeedContractLinkBalance",
			toFloat64(envelope1.LinkBalance), // balance
			feedConfig.GetID(),               // contractAddress
			feedConfig.GetID(),               // feedID
			chainConfig.GetChainID(),         // chainID
			feedConfig.GetContractStatus(),   // contractStatus
			feedConfig.GetContractType(),     // contractType
			feedConfig.GetName(),             // feedName
			feedConfig.GetPath(),             // feedPath
			chainConfig.GetNetworkID(),       // networkID
			chainConfig.GetNetworkName(),     // networkName
		).Once()
		metrics.On("SetLinkAvailableForPayment",
			toFloat64(envelope1.LinkAvailableForPayment), //link available for payment
			feedConfig.GetID(),                           // feedID
			chainConfig.GetChainID(),                     // chainID
			feedConfig.GetContractStatus(),               // contractStatus
			feedConfig.GetContractType(),                 // contractType
			feedConfig.GetName(),                         // feedName
			feedConfig.GetPath(),                         // feedPath
			chainConfig.GetNetworkID(),                   // networkID
			chainConfig.GetNetworkName(),                 // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(nodes[0].GetName()),    // oracleName
			string(envelope1.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope1.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
		metrics.On("SetOffchainAggregatorAnswers",
			humanizedAnswer,                // answer
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		metrics.On("SetOffchainAggregatorAnswersRaw",
			toFloat64(envelope1.LatestAnswer), // answer
			feedConfig.GetID(),                // contractAddress
			feedConfig.GetID(),                // feedID
			chainConfig.GetChainID(),          // chainID
			feedConfig.GetContractStatus(),    // contractStatus
			feedConfig.GetContractType(),      // contractType
			feedConfig.GetName(),              // feedName
			feedConfig.GetPath(),              // feedPath
			chainConfig.GetNetworkID(),        // networkID
			chainConfig.GetNetworkName(),      // networkName
		).Once()
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
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinRaw",
			toFloat64(envelope1.JuelsPerFeeCoin),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoin",
			humanizedJuelsPerFeeCoin,
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
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
		).Once()
		metrics.On("SetOffchainAggregatorSubmissionReceivedValues",
			humanizedAnswer,                // answer
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
		).Once()
		metrics.On("SetOffchainAggregatorJuelsPerFeeCoinReceivedValues",
			humanizedJuelsPerFeeCoin,       // answer
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
		).Once()
		metrics.On("SetOffchainAggregatorRoundID",
			float64(envelope1.AggregatorRoundID),
			feedConfig.GetID(),
			feedConfig.GetID(),
			chainConfig.GetChainID(),
			feedConfig.GetContractStatus(),
			feedConfig.GetContractType(),
			feedConfig.GetName(),
			feedConfig.GetPath(),
			chainConfig.GetNetworkID(),
			chainConfig.GetNetworkName(),
		).Once()
		exporter.Export(ctx, envelope1)

		metrics.On("SetFeedContractLinkBalance",
			toFloat64(envelope2.LinkBalance), // balance
			feedConfig.GetID(),               // contractAddress
			feedConfig.GetID(),               // feedID
			chainConfig.GetChainID(),         // chainID
			feedConfig.GetContractStatus(),   // contractStatus
			feedConfig.GetContractType(),     // contractType
			feedConfig.GetName(),             // feedName
			feedConfig.GetPath(),             // feedPath
			chainConfig.GetNetworkID(),       // networkID
			chainConfig.GetNetworkName(),     // networkName
		).Once()
		metrics.On("SetLinkAvailableForPayment",
			toFloat64(envelope2.LinkAvailableForPayment), //link available for payment
			feedConfig.GetID(),                           // feedID
			chainConfig.GetChainID(),                     // chainID
			feedConfig.GetContractStatus(),               // contractStatus
			feedConfig.GetContractType(),                 // contractType
			feedConfig.GetName(),                         // feedName
			feedConfig.GetPath(),                         // feedPath
			chainConfig.GetNetworkID(),                   // networkID
			chainConfig.GetNetworkName(),                 // networkName
		).Once()
		metrics.On("SetNodeMetadata",
			chainConfig.GetChainID(),      // chainID
			chainConfig.GetNetworkID(),    // networkID
			chainConfig.GetNetworkName(),  // networkName
			string(nodes[0].GetName()),    // oracleName
			string(envelope2.Transmitter), // sender
		).Once()
		metrics.On("SetHeadTrackerCurrentHead",
			float64(envelope2.BlockNumber), // blockNumber
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetChainID(),       // chainID
			chainConfig.GetNetworkID(),     // networkID
		).Once()
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
		).Once()
		exporter.Export(ctx, envelope2)

		metrics.On("Cleanup",
			chainConfig.GetNetworkName(),   // networkName
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetChainID(),       // chainID
			string(nodes[0].GetName()),     // oracleName
			string(envelope1.Transmitter),  // sender
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			feedConfig.GetSymbol(),         // symbol
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
		).Once()
		exporter.Cleanup(ctx)

		metrics.AssertNumberOfCalls(t, "SetOffchainAggregatorAnswers", 1)
		metrics.AssertNumberOfCalls(t, "SetOffchainAggregatorAnswersRaw", 1)
		metrics.AssertNumberOfCalls(t, "IncOffchainAggregatorAnswersTotal", 1)
		metrics.AssertNumberOfCalls(t, "SetOffchainAggregatorSubmissionReceivedValues", 1)
	})
	t.Run("should emit transaction results metrics", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		metrics := NewMetricsMock(t)
		factory := NewPrometheusExporterFactory(log, metrics)

		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		nodes := []NodeConfig{generateNodeConfig()}

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
		).Once()
		exporter, err := factory.NewExporter(ExporterParams{chainConfig, feedConfig, nodes})
		require.NoError(t, err)

		txResults := generateTxResults()
		metrics.On("SetFeedContractTransactionsSucceeded",
			float64(txResults.NumSucceeded), // succeeded
			feedConfig.GetID(),              // contractAddress
			feedConfig.GetID(),              // feedID
			chainConfig.GetChainID(),        // chainID
			feedConfig.GetContractStatus(),  // contractStatus
			feedConfig.GetContractType(),    // contractType
			feedConfig.GetName(),            // feedName
			feedConfig.GetPath(),            // feedPath
			chainConfig.GetNetworkID(),      // networkID
			chainConfig.GetNetworkName(),    // networkName
		).Once()
		metrics.On("SetFeedContractTransactionsFailed",
			float64(txResults.NumFailed),   // failed
			feedConfig.GetID(),             // contractAddress
			feedConfig.GetID(),             // feedID
			chainConfig.GetChainID(),       // chainID
			feedConfig.GetContractStatus(), // contractStatus
			feedConfig.GetContractType(),   // contractType
			feedConfig.GetName(),           // feedName
			feedConfig.GetPath(),           // feedPath
			chainConfig.GetNetworkID(),     // networkID
			chainConfig.GetNetworkName(),   // networkName
		).Once()
		exporter.Export(ctx, txResults)
	})
}
