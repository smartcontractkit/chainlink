package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/smartcontractkit/chainlink-common/pkg/monitoring/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

func TestFeedMonitor(t *testing.T) {
	t.Run("processes updates from multiple pollers", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		cfg := config.Config{}
		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()
		nodes := []NodeConfig{generateNodeConfig()}

		sourceFactory1 := &fakeRandomDataSourceFactory{make(chan interface{})}
		source1, err := sourceFactory1.NewSource(chainConfig, feedConfig)
		require.NoError(t, err)

		sourceFactory2 := &fakeRandomDataSourceFactory{make(chan interface{})}
		source2, err := sourceFactory2.NewSource(chainConfig, feedConfig)
		require.NoError(t, err)

		var bufferCapacity uint32 // no buffering

		pollInterval := 100 * time.Millisecond
		readTimeout := 100 * time.Millisecond

		poller1 := NewSourcePoller(
			source1,
			newNullLogger(),
			pollInterval, readTimeout,
			bufferCapacity,
		)
		poller2 := NewSourcePoller(
			source2,
			newNullLogger(),
			pollInterval, readTimeout,
			bufferCapacity,
		)

		var subs utils.Subprocesses
		defer subs.Wait()
		subs.Go(func() {
			poller1.Run(ctx)
		})
		subs.Go(func() {
			poller2.Run(ctx)
		})

		producer := fakeProducer{make(chan producerMessage), make(chan struct{})}
		defer func() { assert.NoError(t, producer.Close()) }()

		transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}
		configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec, SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic)}

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
		prometheusExporter, err := prometheusExporterFactory.NewExporter(ExporterParams{
			chainConfig,
			feedConfig,
			nodes,
		})
		require.NoError(t, err)
		kafkaExporter, err := kafkaExporterFactory.NewExporter(ExporterParams{
			chainConfig,
			feedConfig,
			nodes,
		})
		require.NoError(t, err)

		exporters := []Exporter{prometheusExporter, kafkaExporter}

		monitor := NewFeedMonitor(
			newNullLogger(),
			[]Poller{poller1, poller2},
			exporters,
		)
		subs.Go(func() {
			monitor.Run(ctx)
		})

		envelope, err := generateEnvelope(ctx)
		require.NoError(t, err)

		var countEnvelopes int64
		var countMessages int64

	LOOP:
		for {
			select {
			case sourceFactory1.updates <- envelope:
				countEnvelopes++
			case sourceFactory2.updates <- envelope:
				countEnvelopes++
			case <-producer.sendCh:
				countMessages++
			case <-ctx.Done():
				break LOOP
			}
		}

		// There should be two prometheus metrics for each envelope + a little bit of wiggle room.
		require.GreaterOrEqual(t, countMessages, 2*countEnvelopes-1)
		require.LessOrEqual(t, countMessages, 2*countEnvelopes+1)
	})
	t.Run("cleanup is called once for each exporter", func(t *testing.T) {
		// put timers on exports of 100ms + keep a counter of all the running exporters.
		// check how many running exporters still execute when Cleanup happens.
		poller := &fakePoller{0, make(chan interface{})}
		exporter1 := new(ExporterMock)
		exporter2 := new(ExporterMock)

		monitor := NewFeedMonitor(
			newNullLogger(),
			[]Poller{poller},
			[]Exporter{exporter1, exporter2},
		)

		var subs utils.Subprocesses
		ctx, cancel := context.WithCancel(context.Background())

		subs.Go(func() {
			monitor.Run(ctx)
		})

		exporter1.On("Export", mock.Anything, mock.Anything).Once()
		exporter1.On("Cleanup", mock.Anything).Once()

		exporter2.On("Export", mock.Anything, mock.Anything).Once()
		exporter2.On("Cleanup", mock.Anything).Once()

		poller.ch <- "update"
		<-time.After(100 * time.Millisecond)
		cancel()
		subs.Wait()

		mock.AssertExpectationsForObjects(t, exporter1)
		mock.AssertExpectationsForObjects(t, exporter2)
	})
	t.Run("panics during Export() or Cleanup() get reported but don't crash the monitor", func(t *testing.T) {
		poller := &fakePoller{0, make(chan interface{})}
		exporter := new(ExporterMock)

		monitor := NewFeedMonitor(
			newNullLogger(),
			[]Poller{poller},
			[]Exporter{exporter},
		)

		var subs utils.Subprocesses
		ctx, cancel := context.WithCancel(context.Background())

		subs.Go(func() {
			monitor.Run(ctx)
		})

		exporter.On("Export", mock.Anything, mock.Anything).Once()
		exporter.On("Export", mock.Anything, mock.Anything).Panic("some error during Export()").Once()
		exporter.On("Export", mock.Anything, mock.Anything).Once()
		exporter.On("Cleanup", mock.Anything).Panic("some error during Cleanup()").Once()

		poller.ch <- "update-before-panic"
		<-time.After(100 * time.Millisecond)
		poller.ch <- "update-causes-panic"
		<-time.After(100 * time.Millisecond)
		poller.ch <- "update-after-panic"
		<-time.After(100 * time.Millisecond)
		cancel()
		subs.Wait()

		mock.AssertExpectationsForObjects(t, exporter)
	})
}
