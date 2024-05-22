package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/smartcontractkit/chainlink-common/pkg/monitoring/config"
)

func TestKafkaExporter(t *testing.T) {
	t.Run("one call to export translates into two prometheus messages", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		producer := fakeProducer{make(chan producerMessage), make(chan struct{})}
		defer func() { assert.NoError(t, producer.Close()) }()
		cfg := config.Config{}
		cfg.Kafka.TransmissionTopic = "transmissions"
		cfg.Kafka.ConfigSetSimplifiedTopic = "config-set-simplified"
		transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}
		configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec, SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic)}
		factory, err := NewKafkaExporterFactory(
			log, producer,
			[]Pipeline{
				{cfg.Kafka.TransmissionTopic, MakeTransmissionMapping, transmissionSchema},
				{cfg.Kafka.ConfigSetSimplifiedTopic, MakeConfigSetSimplifiedMapping, configSetSimplifiedSchema},
			},
		)
		require.NoError(t, err)
		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		nodes := []NodeConfig{generateNodeConfig()}
		exporter, err := factory.NewExporter(ExporterParams{chainConfig, feedConfig, nodes})
		require.NoError(t, err)
		envelope, err := generateEnvelope(ctx)
		require.NoError(t, err)

		go exporter.Export(ctx, envelope)

		var receivedTransmission, receivedConfigSetSimplified producerMessage
		for i := 0; i < 2; i++ {
			select {
			case message := <-producer.sendCh:
				if message.topic == cfg.Kafka.TransmissionTopic {
					receivedTransmission = message
				} else if message.topic == cfg.Kafka.ConfigSetSimplifiedTopic {
					receivedConfigSetSimplified = message
				} else {
					t.Fatalf("received unexpected message with topic %s", message.topic)
				}
			case <-ctx.Done():
				break
			}
		}
		require.NotNil(t, receivedTransmission)
		require.Equal(t, receivedTransmission.topic, cfg.Kafka.TransmissionTopic)
		require.Equal(t, receivedTransmission.key, feedConfig.GetContractAddressBytes())
		require.NotNil(t, receivedConfigSetSimplified)
		require.Equal(t, receivedConfigSetSimplified.topic, cfg.Kafka.ConfigSetSimplifiedTopic)
		require.Equal(t, receivedConfigSetSimplified.key, feedConfig.GetContractAddressBytes())

		// Checking whether the right payload is written to corresponding topic.

		decodedTransmission, err := transmissionSchema.Decode(receivedTransmission.value)
		require.NoError(t, err)
		transmission, ok := decodedTransmission.(map[string]interface{})
		require.True(t, ok)
		answer, ok := transmission["answer"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, answer["data"], envelope.LatestAnswer.Bytes())

		decodedConfigSetSimplified, err := configSetSimplifiedSchema.Decode(receivedConfigSetSimplified.value)
		require.NoError(t, err)
		configSetSimplified, ok := decodedConfigSetSimplified.(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, configSetSimplified["block_number"], uint64ToBeBytes(envelope.BlockNumber))
	})
}
