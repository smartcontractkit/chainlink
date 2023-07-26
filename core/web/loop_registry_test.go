package web_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type mockLoopImpl struct {
	t *testing.T
	*plugins.PromServer
}

// test prom var to avoid collision with real chainlink metrics
var (
	testRegistry   = prometheus.NewRegistry()
	testHandler    = promhttp.HandlerFor(testRegistry, promhttp.HandlerOpts{})
	testMetricName = "super_great_counter"
	testMetric     = prometheus.NewCounter(prometheus.CounterOpts{
		Name: testMetricName,
	})
)

func configurePromRegistry() {
	testRegistry.MustRegister(testMetric)
}

func newMockLoopImpl(t *testing.T, port int) *mockLoopImpl {
	return &mockLoopImpl{
		t:          t,
		PromServer: plugins.NewPromServer(port, logger.TestLogger(t).Named("mock-loop"), plugins.WithHandler(testHandler)),
	}
}

func (m *mockLoopImpl) start() {
	require.NoError(m.t, m.PromServer.Start())
}

func (m *mockLoopImpl) close() {
	require.NoError(m.t, m.PromServer.Close())
}

func (m *mockLoopImpl) run() {
	testMetric.Inc()
}

func TestLoopRegistry(t *testing.T) {

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V1.Enabled = ptr(true)
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
	})
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey)
	// shim a reference to the promserver that is running in our mock loop
	// this ensures the client.Get calls below have a reference to mock loop impl

	expectedEndPoint := "/plugins/mockLoopImpl/metrics"

	require.NoError(t, app.KeyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, app.Start(testutils.Context(t)))

	// register a mock loop
	loop, err := app.GetLoopRegistry().Register("mockLoopImpl")
	require.NoError(t, err)
	require.NotNil(t, loop)
	require.Len(t, app.GetLoopRegistry().List(), 1)

	// set up a test prometheus registry and test metric that is used by
	// our mock loop impl and isolated from the default prom register
	configurePromRegistry()
	mockLoop := newMockLoopImpl(t, loop.EnvCfg.PrometheusPort())
	mockLoop.start()
	defer mockLoop.close()
	mockLoop.run()

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	t.Run("discovery endpoint", func(t *testing.T) {
		// under the covers this is routing thru the app into loop registry
		resp, cleanup := client.Get("/discovery")
		t.Cleanup(cleanup)
		cltest.AssertServerResponse(t, resp, http.StatusOK)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		t.Logf("discovery response %s", b)
		require.Contains(t, string(b), expectedEndPoint)
	})

	t.Run("plugin metrics OK", func(t *testing.T) {
		// plugin name `mockLoopImpl` matches key in PluginConfigs
		resp, cleanup := client.Get(expectedEndPoint)
		t.Cleanup(cleanup)
		cltest.AssertServerResponse(t, resp, http.StatusOK)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		t.Logf("plugin metrics response %s", b)

		var (
			exceptedCount  = 1
			expectedMetric = fmt.Sprintf("%s %d", testMetricName, exceptedCount)
		)
		require.Contains(t, string(b), expectedMetric)
	})

	t.Run("no existent plugin metrics ", func(t *testing.T) {
		// request plugin that doesn't exist
		resp, cleanup := client.Get("/plugins/noexist/metrics")
		t.Cleanup(cleanup)
		cltest.AssertServerResponse(t, resp, http.StatusNotFound)
	})
}
