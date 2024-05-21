package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
)

func TestBuildInfoController_Show_APICredentials(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := app.NewHTTPClient(nil)

	resp, cleanup := client.Get("/v2/build_info")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	body := string(cltest.ParseResponseBody(t, resp))

	require.Contains(t, strings.TrimSpace(body), "commitSHA")
	require.Contains(t, strings.TrimSpace(body), "version")
}

func TestBuildInfoController_Show_NoCredentials(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	url := app.Server.URL + "/v2/build_info"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
