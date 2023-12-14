package web_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/mocks"
)

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
	healthJSON string
)

func TestHealthController_Health_body(t *testing.T) {
	app := cltest.NewApplicationWithKey(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/health")
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// pretty print for comparison
	var b bytes.Buffer
	require.NoError(t, json.Indent(&b, body, "", "  "))
	body = b.Bytes()

	assert.Equal(t, healthJSON, string(body))
}
