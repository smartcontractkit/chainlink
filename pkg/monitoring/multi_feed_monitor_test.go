package monitoring

import (
	"context"
	"sync"
	"sync/atomic"
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

func TestMultiFeedMonitorErroringFactories(t *testing.T) {
	t.Run("a SourceFactory and an ExporterFactory fail", func(t *testing.T) {
		feeds := []FeedConfig{generateFeedConfig()}

		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

		sourceFactory1 := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}
		sourceFactory2 := &fakeSourceFactoryWithError{make(chan interface{}), make(chan error), true}
		sourceFactory3 := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}

		exporterFactory1 := &fakeExporterFactory{make(chan interface{}), false}
		exporterFactory2 := &fakeExporterFactory{make(chan interface{}), true} // factory errors out on NewExporter.
		exporterFactory3 := &fakeExporterFactory{make(chan interface{}), false}

		chainCfg := fakeChainConfig{}
		monitor := NewMultiFeedMonitor(
			chainCfg,
			newNullLogger(),
			[]SourceFactory{sourceFactory1, sourceFactory2, sourceFactory3},
			[]ExporterFactory{exporterFactory1, exporterFactory2, exporterFactory3},
		)

		envelope, err := generateEnvelope()
		require.NoError(t, err)

		wg.Add(1)
		go func() {
			defer wg.Done()
			monitor.Run(ctx, feeds)
		}()

		wg.Add(2)
		for _, factory := range []*fakeRandomDataSourceFactory{
			sourceFactory1, sourceFactory3,
		} {
			go func(factory *fakeRandomDataSourceFactory) {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					select {
					case factory.updates <- envelope:
					case <-ctx.Done():
						return
					}
				}
			}(factory)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				select {
				case sourceFactory2.updates <- envelope:
				case <-ctx.Done():
					return
				}
			}
		}()

		var countMessages int64 = 0
		wg.Add(1)
		go func() {
			defer wg.Done()
		LOOP:
			for {
				select {
				case _ = <-exporterFactory1.data:
					atomic.AddInt64(&countMessages, 1)
				case _ = <-exporterFactory2.data:
					atomic.AddInt64(&countMessages, 1)
				case _ = <-exporterFactory3.data:
					atomic.AddInt64(&countMessages, 1)
				case <-ctx.Done():
					break LOOP
				}
			}
		}()

		<-time.After(100 * time.Millisecond)
		cancel()
		wg.Wait()

		// Two sources produce 10 messages each (the third source is broken) and two exporters ingest each message.
		require.GreaterOrEqual(t, countMessages, int64(10*2*2))
	})
}
