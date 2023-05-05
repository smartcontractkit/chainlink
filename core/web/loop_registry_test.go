package web_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
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
	counter prometheus.Counter
}

// test prom var to avoid collision with real chainlink metrics
var (
	testRegistry   = prometheus.NewRegistry()
	testMetricName = "super_great_counter"
	testMetric     = prometheus.NewCounter(prometheus.CounterOpts{
		Name: testMetricName,
	})
)

func configurePromRegistry(t *testing.T) {
	testRegistry.MustRegister(testMetric)
}

func newMockLoopImpl(t *testing.T) *mockLoopImpl {
	return &mockLoopImpl{
		t:          t,
		PromServer: plugins.NewPromServer(0, logger.TestLogger(t).Named("mock-loop"), plugins.WithRegistry(testRegistry)),
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

	// set up a test register and test metric that is used by
	// our mock loop impl
	configurePromRegistry(t)

	mockLoop := newMockLoopImpl(t)
	mockLoop.start()
	defer mockLoop.close()
	mockLoop.run()

	// shim a reference to the promserver that is running in our mock loop
	// this ensures the client.Get calls below have a reference to mock loop impl

	app.LOOPConfigs = map[string]plugins.EnvConfigurer{
		"foo": plugins.NewEnvConfig(app.Config.LogLevel(), app.Config.JSONConsole(), app.Config.LogUnixTimestamps(), mockLoop.PromServer.Port()),
	}

	require.NoError(t, app.KeyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, app.Start(testutils.Context(t)))

	require.Len(t, app.GetLoopEnvConfig(), 1)

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	// under the covers this is routing thru the app into loop registry
	resp, cleanup := client.Get("/discovery")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("disco response %s", b)

	// plugin name `foo` matches key in PluginConfigs
	resp, cleanup2 := client.Get("/plugins/foo/metrics")
	t.Cleanup(cleanup2)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	b, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("plugin metrics response %s", b)

	var (
		exceptedCount  = 1
		expectedMetric = fmt.Sprintf("%s %d", testMetricName, exceptedCount)
	)
	require.Contains(t, string(b), expectedMetric)

}
