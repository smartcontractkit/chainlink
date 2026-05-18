package web_test

import (
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
)

func TestLCAController_FindLCA(t *testing.T) {
	cfg := configtest.NewTestGeneralConfig(t)
	ec := setupEthClientForControllerTests(t)
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey, ec)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Get("/v2/find_lca?evmChainID=1")
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(b), "chain id does not match any local chains")
}
