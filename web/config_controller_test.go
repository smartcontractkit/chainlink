package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigController_Show(t *testing.T) {
	t.Parallel()

	config, _ := cltest.NewConfigWithPrivateKey()
	app, cleanup := cltest.NewApplicationWithConfig(config)
	defer cleanup()
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/config")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	cfg := make(map[string]string)
	require.NoError(t, cltest.ParseJSONAPIResponse(resp, &cfg))

	assert.Equal(t, store.LogLevel{Level: -1}, cfg["logLevel"])
}
