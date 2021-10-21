package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestCors_DefaultOrigins(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestEVMConfig(t)
	config.GeneralConfig.Overrides.AllowOrigins = null.StringFrom("http://localhost:3000,http://localhost:6689")
	config.GeneralConfig.Overrides.EthereumDisabled = null.BoolFrom(true)
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	tests := []struct {
		origin     string
		statusCode int
	}{
		{"http://localhost:3000", http.StatusOK},
		{"http://localhost:6689", http.StatusOK},
		{"http://localhost:1234", http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.origin, func(t *testing.T) {
			headers := map[string]string{"Origin": test.origin}
			resp, cleanup := client.Get("/v2/config", headers)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, test.statusCode)
		})
	}
}

func TestCors_OverrideOrigins(t *testing.T) {
	t.Parallel()

	tests := []struct {
		allow      string
		origin     string
		statusCode int
	}{
		{"http://chainlink.com", "http://chainlink.com", http.StatusOK},
		{"http://chainlink.com", "http://localhost:3000", http.StatusForbidden},
		{"*", "http://chainlink.com", http.StatusOK},
		{"*", "http://localhost:3000", http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.origin, func(t *testing.T) {
			config := cltest.NewTestEVMConfig(t)
			config.GeneralConfig.Overrides.AllowOrigins = null.StringFrom(test.allow)
			config.GeneralConfig.Overrides.EthereumDisabled = null.BoolFrom(true)
			app, cleanup := cltest.NewApplicationWithConfig(t, config)
			defer cleanup()

			client := app.NewHTTPClient()

			headers := map[string]string{"Origin": test.origin}
			resp, cleanup := client.Get("/v2/config", headers)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, test.statusCode)
		})
	}
}
