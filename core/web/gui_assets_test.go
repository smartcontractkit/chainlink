package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/require"
)

func TestGuiAssets_DefaultIndexHtml(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := &http.Client{}

	// Valid app routes should return OK
	resp, err := client.Get(app.Server.URL + "/")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/runs")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	// Potentially valid app routes should also return OK
	resp, err = client.Get(app.Server.URL + "/valid/route")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/another/valid/route")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	// Bad routes that point to files should return 404
	resp, err = client.Get(app.Server.URL + "/invalidFile.json")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/another/invalidFile.css")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	// Bad API routes should return 404
	resp, err = client.Get(app.Server.URL + "/v2/bad/route")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/v2/another/bad/route")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/v3/new/api/version")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/v123/newer/api/version")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}

func TestGuiAssets_Exact(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := &http.Client{}

	resp, err := client.Get(app.Server.URL + "/main.js")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/mmain.js")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}
