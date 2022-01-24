package monitoring

import (
	"context"
	"sync"
	"testing"

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

func BenchmarkMultichainMonitor(b *testing.B) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Config{}
	chainCfg := generateChainConfig()
	feeds := []FeedConfig{generateFeedConfig()}

	transmissionSchema := fakeSchema{transmissionCodec}
	configSetSimplifiedSchema := fakeSchema{configSetSimplifiedCodec}

	producer := fakeProducer{make(chan producerMessage), ctx}
	factory := &fakeRandomDataSourceFactory{make(chan Envelope), ctx}

	monitor := NewMultiFeedMonitor(
		chainCfg,

		newNullLogger(),
		factory,
		producer,
		&devnullMetrics{},

		cfg.Kafka.TransmissionTopic,
		cfg.Kafka.ConfigSetSimplifiedTopic,

		transmissionSchema,
		configSetSimplifiedSchema,
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
	for i := 0; i < b.N; i++ {
		select {
		case factory.updates <- envelope:
		case <-ctx.Done():
			continue
		}
		// for each update from the chain, the system produces two kafka updates:
		// transmissions and config_set_simplified.
		select {
		case <-producer.sendCh:
		case <-ctx.Done():
			continue
		}
		select {
		case <-producer.sendCh:
		case <-ctx.Done():
			continue
		}
	}
}
