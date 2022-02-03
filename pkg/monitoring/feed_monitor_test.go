package monitoring

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestFeedMonitor(t *testing.T) {
	t.Run("processes updates from multiple pollers", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		wg := &sync.WaitGroup{}
		defer wg.Wait()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		cfg := config.Config{}
		chainConfig := generateChainConfig()
		feedConfig := generateFeedConfig()

		sourceFactory1 := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}
		source1, err := sourceFactory1.NewSource(chainConfig, feedConfig)
		require.NoError(t, err)

		sourceFactory2 := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}
		source2, err := sourceFactory2.NewSource(chainConfig, feedConfig)
		require.NoError(t, err)

		var bufferCapacity uint32 = 0 // no buffering

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

		wg.Add(1)
		go func() {
			defer wg.Done()
			poller1.Run(ctx)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			poller2.Run(ctx)
		}()

		producer := fakeProducer{make(chan producerMessage), ctx}

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
		prometheusExporter, err := prometheusExporterFactory.NewExporter(
			chainConfig,
			feedConfig,
		)
		require.NoError(t, err)
		kafkaExporter, err := kafkaExporterFactory.NewExporter(
			chainConfig,
			feedConfig,
		)
		require.NoError(t, err)

		exporters := []Exporter{prometheusExporter, kafkaExporter}

		monitor := NewFeedMonitor(
			newNullLogger(),
			[]Poller{poller1, poller2},
			exporters,
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			monitor.Run(ctx)
		}()

		envelope, err := generateEnvelope()
		require.NoError(t, err)

		var countEnvelopes int64 = 0
		var countMessages int64 = 0

	LOOP:
		for {
			select {
			case sourceFactory1.updates <- envelope:
				countEnvelopes += 1
			case sourceFactory2.updates <- envelope:
				countEnvelopes += 1
			case _ = <-producer.sendCh:
				countMessages += 1
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
		exporter1 := new(mocks.Exporter)
		exporter2 := new(mocks.Exporter)

		monitor := NewFeedMonitor(
			newNullLogger(),
			[]Poller{poller},
			[]Exporter{exporter1, exporter2},
		)

		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())

		wg.Add(1)
		go func() {
			defer wg.Done()
			monitor.Run(ctx)
		}()

		exporter1.On("Export", mock.Anything, mock.Anything).Once()
		exporter1.On("Cleanup", mock.Anything).Once()

		exporter2.On("Export", mock.Anything, mock.Anything).Once()
		exporter2.On("Cleanup", mock.Anything).Once()

		poller.ch <- "update"
		<-time.After(100 * time.Millisecond)
		cancel()
		wg.Wait()

		mock.AssertExpectationsForObjects(t, exporter1)
		mock.AssertExpectationsForObjects(t, exporter2)
	})
}
