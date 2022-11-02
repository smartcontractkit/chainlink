package web_test

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
)

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestConfigController_Show(t *testing.T) {
	t.Parallel()

	app := cltest.NewLegacyApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)

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
	assert.Equal(t, big.NewInt(evmclient.NullClientChainID).String(), cp.DefaultChainID)
	assert.Equal(t, app.Config.BlockBackfillDepth(), cp.BlockBackfillDepth)
}
