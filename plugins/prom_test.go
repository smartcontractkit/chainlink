package plugins

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

func TestPromServer(t *testing.T) {

	testReg := prometheus.NewRegistry()
	testHandler := promhttp.HandlerFor(testReg, promhttp.HandlerOpts{})
	testMetric := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_metric",
	})
	testReg.MustRegister(testMetric)
	testMetric.Inc()

	s := NewPromServer(0, logger.Test(t), WithHandler(testHandler))
	// check that port is not resolved yet
	require.Equal(t, -1, s.Port())
	require.NoError(t, s.Start())

	url := fmt.Sprintf("http://localhost:%d/metrics", s.Port())
	resp, err := http.Get(url) //nolint
	require.NoError(t, err)
	require.NoError(t, err, "endpoint %s", url)
	require.NotNil(t, resp.Body)
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(b), "test_metric")
	defer resp.Body.Close()

	require.NoError(t, s.Close())
}
