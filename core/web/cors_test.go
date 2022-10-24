package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

func TestCors_DefaultOrigins(t *testing.T) {
	t.Parallel()

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.WebServer.AllowOrigins = ptr("http://localhost:3000,http://localhost:6689")
	})

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
			app := cltest.NewApplicationWithConfig(t, config)

			client := app.NewHTTPClient(cltest.APIEmailAdmin)

			headers := map[string]string{"Origin": test.origin}
			resp, cleanup := client.Get("/v2/chains/evm", headers)
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

			config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.WebServer.AllowOrigins = ptr(test.allow)
			})
			app := cltest.NewApplicationWithConfig(t, config)

			client := app.NewHTTPClient(cltest.APIEmailAdmin)

			headers := map[string]string{"Origin": test.origin}
			resp, cleanup := client.Get("/v2/chains/evm", headers)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, test.statusCode)
		})
	}
}
