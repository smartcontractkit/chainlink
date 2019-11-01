package web_test

import (
	"net/http"
	"testing"

	"chainlink/core/internal/cltest"
	"github.com/stretchr/testify/require"
)

func TestGuiAssets_WildcardIndexHtml(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	client := &http.Client{}

	resp, err := client.Get(app.Server.URL + "/")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/not_found")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/jjob_specs/abc123")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/runs")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/rruns")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/runs/abc123")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/rruns/abc123")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}

func TestGuiAssets_WildcardRouteInfo(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	client := &http.Client{}

	resp, err := client.Get(app.Server.URL + "/job_specs/abc123/routeInfo.json")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/rrouteInfo.json")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/runs/routeInfo.json")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/job_specs/abc123/runs/rrouteInfo.json")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}

func TestGuiAssets_Exact(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	client := &http.Client{}

	resp, err := client.Get(app.Server.URL + "/main.js")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/mmain.js")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}
