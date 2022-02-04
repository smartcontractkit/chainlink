package monitoring

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
)

// This benchmark measures how many messages end up in the kafka client given
// that the chain readers respond immediately with random data and the rdd poller
// will generate a new set of 5 random feeds every second.

//goos: darwin
//goarch: amd64
//pkg: github.com/smartcontractkit/chainlink-relay/pkg/monitoring
//cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// (10 jan 2022)
//    5719	    184623 ns/op	   91745 B/op	    1482 allocs/op
// (17 jan 2022)
//    6679	    180862 ns/op	   92230 B/op	    1493 allocs/op
// (18 jan 2022)
//   16504	     71478 ns/op	   77515 B/op	     963 allocs/op
// (3 feb 2022
//   59468	     23180 ns/op	    5921 B/op	      61 allocs/op
func BenchmarkManager(b *testing.B) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Config{}
	cfg.Feeds.URL = "http://some-fake-url-just-to-trigger-rdd-polling.com"

	chainCfg := generateChainConfig().(fakeChainConfig)
	chainCfg.ReadTimeout = 0 * time.Second
	chainCfg.PollInterval = 0 * time.Second

	transmissionSchema := fakeSchema{transmissionCodec, SubjectFromTopic(cfg.Kafka.TransmissionTopic)}

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
		0, // bufferCapacity for source pollers
	)

	rddPoller := NewSourcePoller(
		NewFakeRDDSource(5, 6), // Always produce 1 random feed.
		newNullLogger(),
		2*time.Second, // cfg.Feeds.RDDPollInterval,
		1*time.Second, // cfg.Feeds.RDDReadTimeout,
		0,             // no buffering!
	)

	manager := NewManager(
		newNullLogger(),
		rddPoller,
	)

	envelope, err := generateEnvelope()
	if err != nil {
		b.Fatalf("failed to generate config: %v", err)
	}
	_ = envelope

	wg.Add(1)
	go func() {
		defer wg.Done()
		rddPoller.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		manager.Run(ctx, monitor.Run)
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		factory.updates <- envelope
		<-producer.sendCh
	}
}
