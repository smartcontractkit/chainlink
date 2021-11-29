package web_test

import (
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/config"

	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/require"
)

func TestConfigController_Show(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/config")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	cp := config.ConfigPrinter{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &cp))

	assert.Contains(t, cp.RootDir, "/tmp/chainlink_test/")
	assert.Equal(t, uint16(6688), cp.Port)
	assert.Equal(t, uint16(6689), cp.TLSPort)
	assert.Equal(t, "", cp.TLSHost)
	assert.Len(t, cp.EthereumURL, 0)
	assert.Equal(t, big.NewInt(eth.NullClientChainID).String(), cp.DefaultChainID)
	assert.Contains(t, cp.ClientNodeURL, "http://127.0.0.1:")
	assert.Equal(t, cltest.NewTestGeneralConfig(t).BlockBackfillDepth(), cp.BlockBackfillDepth)
	assert.Equal(t, time.Second*5, cp.DatabaseTimeout.Duration())
}
