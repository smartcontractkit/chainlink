package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestCors_DefaultOrigins(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfig(t)
	config.Set("ALLOW_ORIGINS", "http://localhost:3000,http://localhost:6689")
	app, appCleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer appCleanup()
	client := app.NewHTTPClient()

	tests := []struct {
		origin     string
		statusCode int
	}{
		{"http://localhost:3000", 200},
		{"http://localhost:6689", 200},
		{"http://localhost:1234", 403},
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
		{"http://chainlink.com", "http://chainlink.com", 200},
		{"http://chainlink.com", "http://localhost:3000", 403},
		{"*", "http://chainlink.com", 200},
		{"*", "http://localhost:3000", 200},
	}

	for _, test := range tests {
		t.Run(test.origin, func(t *testing.T) {
			config, _ := cltest.NewConfig(t)
			config.Set("ALLOW_ORIGINS", test.allow)
			app, appCleanup := cltest.NewApplicationWithConfigAndKey(t, config)
			defer appCleanup()
			client := app.NewHTTPClient()

			headers := map[string]string{"Origin": test.origin}
			resp, cleanup := client.Get("/v2/config", headers)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, test.statusCode)
		})
	}
}
