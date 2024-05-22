package monitoring

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/smartcontractkit/chainlink-common/pkg/monitoring/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

const numFeeds = 10

func TestMultiFeedMonitorSynchronousMode(t *testing.T) {
	// Synchronous mode means that the a source update is produced and the
	// corresponding exporter message is consumed in the same goroutine.
	defer goleak.VerifyNone(t)

	var subs utils.Subprocesses
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := config.Config{}
	chainCfg := fakeChainConfig{}
	chainCfg.ReadTimeout = 1 * time.Second
	chainCfg.PollInterval = 5 * time.Second
	feeds := make([]FeedConfig, numFeeds)
	for i := 0; i < numFeeds; i++ {
		feeds[i] = generateFeedConfig()
	}
	nodes := []NodeConfig{generateNodeConfig()}

	transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec, SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic)}

	producer := fakeProducer{make(chan producerMessage), make(chan struct{})}
	defer func() { assert.NoError(t, producer.Close()) }()
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
	subs.Go(func() {
		monitor.Run(ctx, RDDData{feeds, nodes})
	})

	count := 0
	messages := []producerMessage{}

	envelope, err := generateEnvelope(ctx)
	require.NoError(t, err)

LOOP:
	for {
		select {
		case factory.updates <- envelope:
			count++
			envelope, err = generateEnvelope(ctx)
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

	subs.Wait()
	require.Equal(t, 10, count, "should only be able to do initial read of the chain")
	require.Equal(t, 20, len(messages))
}

func TestMultiFeedMonitorForPerformance(t *testing.T) {
	defer goleak.VerifyNone(t)

	var subs utils.Subprocesses
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := config.Config{}
	chainCfg := fakeChainConfig{}
	chainCfg.ReadTimeout = 1 * time.Second
	chainCfg.PollInterval = 5 * time.Second
	feeds := []FeedConfig{}
	for i := 0; i < numFeeds; i++ {
		feeds = append(feeds, generateFeedConfig())
	}
	nodes := []NodeConfig{generateNodeConfig()}

	transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec, SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic)}

	producer := fakeProducer{make(chan producerMessage), make(chan struct{})}
	defer func() { assert.NoError(t, producer.Close()) }()
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
	subs.Go(func() {
		monitor.Run(ctx, RDDData{feeds, nodes})
	})

	var count int64
	messages := []producerMessage{}

	envelope, err := generateEnvelope(ctx)
	require.NoError(t, err)

	subs.Go(func() {
	LOOP:
		for {
			select {
			case factory.updates <- envelope:
				count++
				envelope, err = generateEnvelope(ctx)
				require.NoError(t, err)
			case <-ctx.Done():
				break LOOP
			}
		}
	})
	subs.Go(func() {
	LOOP:
		for {
			select {
			case message := <-producer.sendCh:
				messages = append(messages, message)
			case <-ctx.Done():
				break LOOP
			}
		}
	})

	subs.Wait()
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
		nodes := []NodeConfig{
			generateNodeConfig(),
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

		exporterFactory1.On("NewExporter", ExporterParams{chainConfig, feeds[0], nodes}).Return(exporter1, nil)
		exporterFactory2.On("NewExporter", ExporterParams{chainConfig, feeds[0], nodes}).Return(exporter2, nil)
		exporterFactory1.On("NewExporter", ExporterParams{chainConfig, feeds[1], nodes}).Return(nil, fmt.Errorf("exporter_factory1/feed2 failed"))
		exporterFactory2.On("NewExporter", ExporterParams{chainConfig, feeds[1], nodes}).Return(nil, fmt.Errorf("exporter_factory2/feed2 failed"))

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		monitor.Run(ctx, RDDData{feeds, nodes})
	})
	t.Run("one SourceFactory and an ExporterFactory fail for one feed", func(t *testing.T) {
		chainCfg := fakeChainConfig{}
		feeds := []FeedConfig{generateFeedConfig()}
		nodes := []NodeConfig{generateNodeConfig()}

		var subs utils.Subprocesses
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

		sourceFactory1 := &fakeRandomDataSourceFactory{make(chan interface{})}
		sourceFactory2 := &fakeSourceFactoryWithError{make(chan interface{}), make(chan error), true}
		sourceFactory3 := &fakeRandomDataSourceFactory{make(chan interface{})}

		exporterFactory1 := &fakeExporterFactory{make(chan interface{}), false}
		exporterFactory2 := &fakeExporterFactory{make(chan interface{}), true} // factory errors out on NewExporter.
		exporterFactory3 := &fakeExporterFactory{make(chan interface{}), false}

		monitor := NewMultiFeedMonitor(
			chainCfg,
			newNullLogger(),
			[]SourceFactory{sourceFactory1, sourceFactory2, sourceFactory3},
			[]ExporterFactory{exporterFactory1, exporterFactory2, exporterFactory3},
			100, // bufferCapacity for source pollers
		)

		envelope, err := generateEnvelope(ctx)
		require.NoError(t, err)

		subs.Go(func() {
			monitor.Run(ctx, RDDData{feeds, nodes})
		})

		for _, factory := range []*fakeRandomDataSourceFactory{
			sourceFactory1, sourceFactory3,
		} {
			factory := factory
			subs.Go(func() {
				for i := 0; i < 10; i++ {
					select {
					case factory.updates <- envelope:
					case <-ctx.Done():
						return
					}
				}
			})
		}

		subs.Go(func() {
			for i := 0; i < 10; i++ {
				select {
				case sourceFactory2.updates <- envelope:
				case <-ctx.Done():
					return
				}
			}
		})

		var countMessages int64
		subs.Go(func() {
		LOOP:
			for {
				select {
				case <-exporterFactory1.data:
					atomic.AddInt64(&countMessages, 1)
				case <-exporterFactory2.data:
					atomic.AddInt64(&countMessages, 1)
				case <-exporterFactory3.data:
					atomic.AddInt64(&countMessages, 1)
				case <-ctx.Done():
					break LOOP
				}
			}
		})

		<-time.After(100 * time.Millisecond)
		cancel()
		subs.Wait()

		// Two sources produce 10 messages each (the third source is broken) and two exporters ingest each message.
		require.GreaterOrEqual(t, countMessages, int64(10*2*2))
	})
}
