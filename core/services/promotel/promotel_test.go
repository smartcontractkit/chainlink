package promotel

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestForwarder(t *testing.T) {
	var (
		g              = prometheus.DefaultGatherer
		r              = prometheus.DefaultRegisterer
		lggr, observed = logger.TestObserved(t, zap.DebugLevel)
		testMetricName = t.Name() + "_test_counter_metric"
		interval       = 10 * time.Millisecond
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go reportTestMetrics(ctx, r, testMetricName)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for _, l := range observed.All() {
					metricName, ok := l.ContextMap()["name"].(string)
					if ok && strings.Contains(metricName, testMetricName) {
						doneCh <- struct{}{}
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	forwarder, err := NewForwarderService(g, r, lggr, Options{
		Endpoint:    "localhost:4317",
		TLSInsecure: true,
		Interval:    interval,
		Verbose:     true,
	})
	require.NoError(t, err)
	require.NoError(t, forwarder.Start(ctx))
	defer forwarder.Close()

	select {
	case <-ctx.Done():
		t.Fatal("Test timed out. Expected metric not found")
	case <-doneCh:
		t.Log("Found metric.")
	}
}

func reportTestMetrics(ctx context.Context, reg prometheus.Registerer, metricName string) {
	m := promauto.With(reg).NewCounter(prometheus.CounterOpts{Name: metricName})
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m.Inc()
			time.Sleep(1 * time.Second)
		}
	}
}
