package web_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/params"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/mocks"
)

func TestMain(m *testing.M) {
	params.InitCosmosSdk(
		/* bech32Prefix= */ "wasm",
		/* token= */ "cosm",
	)

	os.Exit(m.Run())
}

func TestHealthController_Readyz(t *testing.T) {
	var tt = []struct {
		name   string
		ready  bool
		status int
	}{
		{
			name:   "not ready",
			ready:  false,
			status: http.StatusServiceUnavailable,
		},
		{
			name:   "ready",
			ready:  true,
			status: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app := cltest.NewApplicationWithKey(t)
			healthChecker := new(mocks.Checker)
			healthChecker.On("Start").Return(nil).Once()
			healthChecker.On("IsReady").Return(tc.ready, nil).Once()
			healthChecker.On("Close").Return(nil).Once()

			app.HealthChecker = healthChecker
			require.NoError(t, app.Start(testutils.Context(t)))

			client := app.NewHTTPClient(nil)
			resp, cleanup := client.Get("/readyz")
			t.Cleanup(cleanup)
			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}

func TestHealthController_Health_status(t *testing.T) {
	var tt = []struct {
		name   string
		ready  bool
		status int
	}{
		{
			name:   "not ready",
			ready:  false,
			status: http.StatusMultiStatus,
		},
		{
			name:   "ready",
			ready:  true,
			status: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app := cltest.NewApplicationWithKey(t)
			healthChecker := new(mocks.Checker)
			healthChecker.On("Start").Return(nil).Once()
			healthChecker.On("IsHealthy").Return(tc.ready, nil).Once()
			healthChecker.On("Close").Return(nil).Once()

			app.HealthChecker = healthChecker
			require.NoError(t, app.Start(testutils.Context(t)))

			client := app.NewHTTPClient(nil)
			resp, cleanup := client.Get("/health")
			t.Cleanup(cleanup)
			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}

var (
	//go:embed testdata/body/health.json
	bodyJSON string
	//go:embed testdata/body/health.html
	bodyHTML string
	//go:embed testdata/body/health.txt
	bodyTXT string
	//go:embed testdata/body/health-failing.json
	bodyJSONFailing string
	//go:embed testdata/body/health-failing.html
	bodyHTMLFailing string
	//go:embed testdata/body/health-failing.txt
	bodyTXTFailing string
)

func TestHealthController_Health_body(t *testing.T) {
	for _, tc := range []struct {
		name    string
		path    string
		headers map[string]string
		expBody string
	}{
		{"default", "/health", nil, bodyJSON},
		{"json", "/health", map[string]string{"Accept": gin.MIMEJSON}, bodyJSON},
		{"html", "/health", map[string]string{"Accept": gin.MIMEHTML}, bodyHTML},
		{"text", "/health", map[string]string{"Accept": gin.MIMEPlain}, bodyTXT},
		{".txt", "/health.txt", nil, bodyTXT},

		{"default-failing", "/health?failing", nil, bodyJSONFailing},
		{"json-failing", "/health?failing", map[string]string{"Accept": gin.MIMEJSON}, bodyJSONFailing},
		{"html-failing", "/health?failing", map[string]string{"Accept": gin.MIMEHTML}, bodyHTMLFailing},
		{"text-failing", "/health?failing", map[string]string{"Accept": gin.MIMEPlain}, bodyTXTFailing},
		{".txt-failing", "/health.txt?failing", nil, bodyTXTFailing},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := configtest.NewGeneralConfig(t, func(cfg *chainlink.Config, secrets *chainlink.Secrets) {
				cfg.Cosmos = append(cfg.Cosmos, &coscfg.TOMLConfig{
					ChainID: ptr("Foo"),
					Nodes: coscfg.Nodes{
						{Name: ptr("primary"), TendermintURL: config.MustParseURL("http://tender.mint")},
					},
				})
				cfg.Cosmos[0].SetDefaults()
				cfg.Solana = append(cfg.Solana, &solcfg.TOMLConfig{
					ChainID: ptr("Bar"),
					Nodes: solcfg.Nodes{
						{Name: ptr("primary"), URL: config.MustParseURL("http://solana.web")},
					},
				})
				cfg.Solana[0].SetDefaults()
				cfg.Starknet = append(cfg.Starknet, &stkcfg.TOMLConfig{
					ChainID: ptr("Baz"),
					Nodes: stkcfg.Nodes{
						{Name: ptr("primary"), URL: config.MustParseURL("http://stark.node")},
					},
				})
				cfg.Starknet[0].SetDefaults()
			})
			app := cltest.NewApplicationWithConfigAndKey(t, cfg)
			require.NoError(t, app.Start(testutils.Context(t)))

			client := app.NewHTTPClient(nil)
			resp, cleanup := client.Get(tc.path, tc.headers)
			t.Cleanup(cleanup)
			assert.Equal(t, http.StatusMultiStatus, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tc.expBody == bodyJSON {
				// pretty print for comparison
				var b bytes.Buffer
				require.NoError(t, json.Indent(&b, body, "", "  "))
				body = b.Bytes()
			}
			assert.Equal(t, strings.TrimSpace(tc.expBody), strings.TrimSpace(string(body)))
		})
	}
}
