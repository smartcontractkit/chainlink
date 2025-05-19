package web_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type mockLoopImpl struct {
	t *testing.T
	*loop.PromServer
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
		PromServer: loop.PromServerOpts{Handler: testHandler}.New(port, logger.TestLogger(t).Named("mock-loop")),
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
	ctx := testutils.Context(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.OCR.Enabled = ptr(true)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
	})
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey)
	// shim a reference to the promserver that is running in our mock loop
	// this ensures the client.Get calls below have a reference to mock loop impl

	expectedLooppEndPoint, expectedCoreEndPoint := "/plugins/mockLoopImpl/metrics", "/metrics"

	// note we expect this to be an ordered result
	expectedLabels := []model.LabelSet{
		model.LabelSet{"__metrics_path__": model.LabelValue(expectedCoreEndPoint)},
		model.LabelSet{"__metrics_path__": model.LabelValue(expectedLooppEndPoint)},
	}

	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, app.Start(testutils.Context(t)))

	// register a mock loop
	loop, err := app.GetLoopRegistry().Register("mockLoopImpl")
	require.NoError(t, err)
	require.NotNil(t, loop)
	require.Len(t, app.GetLoopRegistry().List(), 1)

	// set up a test prometheus registry and test metric that is used by
	// our mock loop impl and isolated from the default prom register
	configurePromRegistry()
	mockLoop := newMockLoopImpl(t, loop.EnvCfg.PrometheusPort)
	mockLoop.start()
	defer mockLoop.close()
	mockLoop.run()

	client := app.NewHTTPClient(nil)

	t.Run("discovery endpoint", func(t *testing.T) {
		// under the covers this is routing thru the app into loop registry
		resp, cleanup := client.Get("/discovery")
		t.Cleanup(cleanup)
		cltest.AssertServerResponse(t, resp, http.StatusOK)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		t.Logf("discovery response %s", b)
		var got []*targetgroup.Group
		require.NoError(t, json.Unmarshal(b, &got))

		gotLabels := make([]model.LabelSet, 0)
		for _, ls := range got {
			gotLabels = append(gotLabels, ls.Labels)
		}
		assert.Equal(t, len(expectedLabels), len(gotLabels))
		for i := range expectedLabels {
			assert.EqualValues(t, expectedLabels[i], gotLabels[i])
		}
	})

	t.Run("plugin metrics OK", func(t *testing.T) {
		// plugin name `mockLoopImpl` matches key in PluginConfigs
		resp, cleanup := client.Get(expectedLooppEndPoint)
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

	t.Run("core metrics OK", func(t *testing.T) {
		// core node metrics endpoint
		resp, cleanup := client.Get(expectedCoreEndPoint)
		t.Cleanup(cleanup)
		cltest.AssertServerResponse(t, resp, http.StatusOK)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		t.Logf("core metrics response %s", b)
	})

	t.Run("no existent plugin metrics ", func(t *testing.T) {
		// request plugin that doesn't exist
		resp, cleanup := client.Get("/plugins/noexist/metrics")
		t.Cleanup(cleanup)
		cltest.AssertServerResponse(t, resp, http.StatusNotFound)
	})
}
