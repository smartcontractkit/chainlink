package monitoring

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const numFeeds = 10

func TestMultiFeedMonitorSynchronousMode(t *testing.T) {
	// Synchronous mode means that the a source update is produced and the
	// corresponding exporter message is consumed in the same goroutine.
	defer goleak.VerifyNone(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}

	cfg := config.Config{}
	chainCfg := fakeChainConfig{}
	chainCfg.PollInterval = 5 * time.Second
	feeds := make([]FeedConfig, numFeeds)
	for i := 0; i < numFeeds; i++ {
		feeds[i] = generateFeedConfig()
	}

	transmissionSchema := fakeSchema{transmissionCodec}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec}

	producer := fakeProducer{make(chan producerMessage), ctx}
	factory := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}

	prometheusExporterFactory := NewPrometheusExporterFactory(
		newNullLogger(),
		&devnullMetrics{},
	)
	kafkaExporterFactory := NewKafkaExporterFactory(
		newNullLogger(),
		producer,

		transmissionSchema,
		configSetSimplifiedSchema,

		cfg.Kafka.ConfigSetSimplifiedTopic,
		cfg.Kafka.TransmissionTopic,
	)

	monitor := NewMultiFeedMonitor(
		chainCfg,
		newNullLogger(),

		[]SourceFactory{factory},
		[]ExporterFactory{prometheusExporterFactory, kafkaExporterFactory},
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitor.Run(ctx, feeds)
	}()

	count := 0
	messages := []producerMessage{}

	envelope, err := generateEnvelope()
	require.NoError(t, err)

LOOP:
	for {
		select {
		case factory.updates <- envelope:
			count += 1
			envelope, err = generateEnvelope()
			require.NoError(t, err)
		case <-ctx.Done():
			break LOOP
		}
		select {
		case message := <-producer.sendCh:
			messages = append(messages, message)
		case <-ctx.Done():
			break LOOP
		}
		select {
		case message := <-producer.sendCh:
			messages = append(messages, message)
		case <-ctx.Done():
			break LOOP
		}
	}

	wg.Wait()
	require.Equal(t, 10, count, "should only be able to do initial read of the chain")
	require.Equal(t, 20, len(messages))
}

func TestMultiFeedMonitorForPerformance(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}

	cfg := config.Config{}
	chainCfg := fakeChainConfig{}
	chainCfg.PollInterval = 5 * time.Second
	feeds := []FeedConfig{}
	for i := 0; i < numFeeds; i++ {
		feeds = append(feeds, generateFeedConfig())
	}

	transmissionSchema := fakeSchema{transmissionCodec}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec}

	producer := fakeProducer{make(chan producerMessage), ctx}
	factory := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}

	prometheusExporterFactory := NewPrometheusExporterFactory(
		newNullLogger(),
		&devnullMetrics{},
	)
	kafkaExporterFactory := NewKafkaExporterFactory(
		newNullLogger(),
		producer,

		transmissionSchema,
		configSetSimplifiedSchema,

		cfg.Kafka.ConfigSetSimplifiedTopic,
		cfg.Kafka.TransmissionTopic,
	)

	monitor := NewMultiFeedMonitor(
		chainCfg,
		newNullLogger(),

		[]SourceFactory{factory},
		[]ExporterFactory{prometheusExporterFactory, kafkaExporterFactory},
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitor.Run(ctx, feeds)
	}()

	var count int64 = 0
	messages := []producerMessage{}

	envelope, err := generateEnvelope()
	require.NoError(t, err)

	wg.Add(1)
	go func() {
		defer wg.Done()
	LOOP:
		for {
			select {
			case factory.updates <- envelope:
				count += 1
				envelope, err = generateEnvelope()
				require.NoError(t, err)
			case <-ctx.Done():
				break LOOP
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
	LOOP:
		for {
			select {
			case message := <-producer.sendCh:
				messages = append(messages, message)
			case <-ctx.Done():
				break LOOP
			}
		}
	}()

	wg.Wait()
	require.Equal(t, int64(10), count, "should only be able to do initial reads of the chain")
	require.Equal(t, 20, len(messages))
}
