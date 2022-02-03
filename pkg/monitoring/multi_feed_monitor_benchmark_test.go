package monitoring

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
)

// Results:
// goos: darwin
// goarch: amd64
// pkg: github.com/smartcontractkit/chainlink-relay/pkg/monitoring
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// (11 Dec 2021)
//    48993	     35111 ns/op	   44373 B/op	     251 allocs/op
// (13 Dec 2021)
//    47331	     34285 ns/op	   41074 B/op	     235 allocs/op
// (3 Jan 2022)
//    6985	    162187 ns/op	  114802 B/op	    1506 allocs/op
// (4 Jan 2022)
//    9332	    166275 ns/op	  157078 B/op	    1590 allocs/op
// (17 Jan 2022)
//    7374	    202079 ns/op	  164301 B/op	    1712 allocs/op
// (30 Jan 2022)
//    19083	     61491 ns/op	   75157 B/op	     723 allocs/op

func BenchmarkMultiFeedMonitor(b *testing.B) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Config{}
	chainCfg := generateChainConfig().(fakeChainConfig)
	chainCfg.ReadTimeout = 0 * time.Second
	chainCfg.PollInterval = 0 * time.Second
	feeds := []FeedConfig{generateFeedConfig()}

	transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec, SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic)}

	producer := fakeProducer{make(chan producerMessage), ctx}
	factory := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}

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
	if err != nil {
		b.Fatalf("failed to build kafka exporter: %v", err)
	}

	monitor := NewMultiFeedMonitor(
		chainCfg,
		newNullLogger(),
		[]SourceFactory{factory},
		[]ExporterFactory{
			prometheusExporterFactory,
			kafkaExporterFactory,
		},
		100, // bufferCapacity for source pollers
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitor.Run(ctx, feeds)
	}()

	envelope, err := generateEnvelope()
	if err != nil {
		b.Fatalf("failed to generate config: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

BENCH_LOOP:
	for i := 0; i < b.N; i++ {
		//for {
		select {
		case factory.updates <- envelope:
		case <-ctx.Done():
			continue BENCH_LOOP
		}
		// for each update from the chain, the system produces two kafka updates:
		// transmissions and config_set_simplified.
		for {
			select {
			case <-producer.sendCh:
			case <-ctx.Done():
				continue BENCH_LOOP
			default:
				continue BENCH_LOOP
			}
		}
	}
}
