package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestKafkaExporter(t *testing.T) {
	t.Run("one call to export translates into two prometheus messages", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()
		log := newNullLogger()
		producer := fakeProducer{make(chan producerMessage), ctx}
		cfg := config.Config{}
		cfg.Kafka.TransmissionTopic = "transmissions"
		cfg.Kafka.ConfigSetSimplifiedTopic = "config-set-simplified"
		transmissionSchema := fakeSchema{transmissionCodec}
		configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec}
		factory := NewKafkaExporterFactory(
			log, producer,
			transmissionSchema, configSetSimplifiedSchema,
			cfg.Kafka.TransmissionTopic, cfg.Kafka.ConfigSetSimplifiedTopic,
		)
		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		exporter, err := factory.NewExporter(chainConfig, feedConfig)
		require.NoError(t, err)
		envelope, err := generateEnvelope()
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
	})
}
