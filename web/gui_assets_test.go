// +build test

package web_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
)

func TestGuiAssets_WildcardIndexHtml(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp, cleanup := cltest.BasicAuthGet(app.Server.URL + "/")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/not_found")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/jjob_specs/abc123")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/rruns")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs/abc123")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/rruns/abc123")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)
}

func TestGuiAssets_WildcardRouteInfo(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp, cleanup := cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/routeInfo.json")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/rrouteInfo.json")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs/routeInfo.json")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs/rrouteInfo.json")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)
}

func TestGuiAssets_Exact(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp, cleanup := cltest.BasicAuthGet(app.Server.URL + "/main.js")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 200)

	resp, cleanup = cltest.BasicAuthGet(app.Server.URL + "/mmain.js")
	defer cleanup()
	cltest.AssertServerResponse(t, resp, 404)
}
