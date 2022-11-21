package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	clhttptest "github.com/smartcontractkit/chainlink/core/internal/testutils/httptest"

	"github.com/stretchr/testify/require"
)

func TestBuildInfoController_Show_APICredentials(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

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
	require.NoError(t, app.Start(testutils.Context(t)))

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	url := app.Server.URL + "/v2/build_info"
	resp, err := client.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
