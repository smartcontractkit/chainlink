package monitoring

import (
	"context"
	"fmt"
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
	chainCfg.ReadTimeout = 1 * time.Second
	chainCfg.PollInterval = 5 * time.Second
	feeds := make([]FeedConfig, numFeeds)
	for i := 0; i < numFeeds; i++ {
		feeds[i] = generateFeedConfig()
	}

	transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec, SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic)}

	producer := fakeProducer{make(chan producerMessage), ctx}
	factory := &fakeRandomDataSourceFactory{make(chan interface{})}

	prometheusExporterFactory := NewPrometheusExporterFactory(
		newNullLogger(),
		&devnullMetrics{},
	)
	kafkaExporterFactory, err := NewKafkaExporterFactory(
		newNullLogger(),
		producer,
		[]Pipeline{
			{cfg.Kafka.TransmissionTopic, MakeTransmissionMapping, transmissionSchema},
			{cfg.Kafka.ConfigSetSimplifiedTopic, MakeConfigSetSimplifiedMapping, configSetSimplifiedSchema},
		},
	)
	require.NoError(t, err)

	monitor := NewMultiFeedMonitor(
		chainCfg,
		newNullLogger(),
		[]SourceFactory{factory},
		[]ExporterFactory{prometheusExporterFactory, kafkaExporterFactory},
		100, // bufferCapacity for source pollers
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
	chainCfg.ReadTimeout = 1 * time.Second
	chainCfg.PollInterval = 5 * time.Second
	feeds := []FeedConfig{}
	for i := 0; i < numFeeds; i++ {
		feeds = append(feeds, generateFeedConfig())
	}

	transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec, SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic)}

	producer := fakeProducer{make(chan producerMessage), ctx}
	factory := &fakeRandomDataSourceFactory{make(chan interface{})}

	prometheusExporterFactory := NewPrometheusExporterFactory(
		newNullLogger(),
		&devnullMetrics{},
	)
	kafkaExporterFactory, err := NewKafkaExporterFactory(
		newNullLogger(),
		producer,
		[]Pipeline{
			{cfg.Kafka.TransmissionTopic, MakeTransmissionMapping, transmissionSchema},
			{cfg.Kafka.ConfigSetSimplifiedTopic, MakeConfigSetSimplifiedMapping, configSetSimplifiedSchema},
		},
	)
	require.NoError(t, err)

	monitor := NewMultiFeedMonitor(
		chainCfg,
		newNullLogger(),
		[]SourceFactory{factory},
		[]ExporterFactory{prometheusExporterFactory, kafkaExporterFactory},
		100, // bufferCapacity for source pollers
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
	t.Run("all sources fail for one feed and all exporters fail for the other", func(t *testing.T) {
		sourceFactory1 := new(SourceFactoryMock)
		sourceFactory2 := new(SourceFactoryMock)
		source1 := new(SourceMock)
		source2 := new(SourceMock)

		exporterFactory1 := new(ExporterFactoryMock)
		exporterFactory2 := new(ExporterFactoryMock)
		exporter1 := new(ExporterMock)
		exporter2 := new(ExporterMock)

		chainConfig := generateChainConfig()
		feeds := []FeedConfig{
			generateFeedConfig(),
			generateFeedConfig(),
		}

		monitor := NewMultiFeedMonitor(
			chainConfig,
			newNullLogger(),
			[]SourceFactory{sourceFactory1, sourceFactory2},
			[]ExporterFactory{exporterFactory1, exporterFactory2},
			10, // bufferCapacity for source pollers
		)

		sourceFactory1.On("NewSource", chainConfig, feeds[0]).Return(nil, fmt.Errorf("source_factory1/feed1 failed"))
		sourceFactory2.On("NewSource", chainConfig, feeds[0]).Return(nil, fmt.Errorf("source_factory2/feed1 failed"))
		sourceFactory1.On("NewSource", chainConfig, feeds[1]).Return(source1, nil)
		sourceFactory2.On("NewSource", chainConfig, feeds[1]).Return(source2, nil)

		sourceFactory1.On("GetType").Return("fake")
		sourceFactory2.On("GetType").Return("fake")

		exporterFactory1.On("NewExporter", chainConfig, feeds[0]).Return(exporter1, nil)
		exporterFactory2.On("NewExporter", chainConfig, feeds[0]).Return(exporter2, nil)
		exporterFactory1.On("NewExporter", chainConfig, feeds[1]).Return(nil, fmt.Errorf("exporter_factory1/feed2 failed"))
		exporterFactory2.On("NewExporter", chainConfig, feeds[1]).Return(nil, fmt.Errorf("exporter_factory2/feed2 failed"))

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		monitor.Run(ctx, feeds)
	})
	t.Run("one SourceFactory and an ExporterFactory fail for one feed", func(t *testing.T) {
		feeds := []FeedConfig{generateFeedConfig()}

		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

		sourceFactory1 := &fakeRandomDataSourceFactory{make(chan interface{})}
		sourceFactory2 := &fakeSourceFactoryWithError{make(chan interface{}), make(chan error), true}
		sourceFactory3 := &fakeRandomDataSourceFactory{make(chan interface{})}

		exporterFactory1 := &fakeExporterFactory{make(chan interface{}), false}
		exporterFactory2 := &fakeExporterFactory{make(chan interface{}), true} // factory errors out on NewExporter.
		exporterFactory3 := &fakeExporterFactory{make(chan interface{}), false}

		chainCfg := fakeChainConfig{}
		monitor := NewMultiFeedMonitor(
			chainCfg,
			newNullLogger(),
			[]SourceFactory{sourceFactory1, sourceFactory2, sourceFactory3},
			[]ExporterFactory{exporterFactory1, exporterFactory2, exporterFactory3},
			100, // bufferCapacity for source pollers
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
