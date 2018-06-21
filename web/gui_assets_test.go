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

	resp := cltest.BasicAuthGet(app.Server.URL + "/")
	cltest.AssertServerResponse(t, resp, 200)

	resp = cltest.BasicAuthGet(app.Server.URL + "/not_found")
	cltest.AssertServerResponse(t, resp, 404)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123")
	cltest.AssertServerResponse(t, resp, 200)

	resp = cltest.BasicAuthGet(app.Server.URL + "/jjob_specs/abc123")
	cltest.AssertServerResponse(t, resp, 404)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs")
	cltest.AssertServerResponse(t, resp, 200)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/rruns")
	cltest.AssertServerResponse(t, resp, 404)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs/abc123")
	cltest.AssertServerResponse(t, resp, 200)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/rruns/abc123")
	cltest.AssertServerResponse(t, resp, 404)
}

func TestGuiAssets_WildcardRouteInfo(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/routeInfo.json")
	cltest.AssertServerResponse(t, resp, 200)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/rrouteInfo.json")
	cltest.AssertServerResponse(t, resp, 404)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs/routeInfo.json")
	cltest.AssertServerResponse(t, resp, 200)

	resp = cltest.BasicAuthGet(app.Server.URL + "/job_specs/abc123/runs/rrouteInfo.json")
	cltest.AssertServerResponse(t, resp, 404)
}

func TestGuiAssets_Exact(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	resp := cltest.BasicAuthGet(app.Server.URL + "/main.js")
	cltest.AssertServerResponse(t, resp, 200)

	resp = cltest.BasicAuthGet(app.Server.URL + "/mmain.js")
	cltest.AssertServerResponse(t, resp, 404)
}
