package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuiAssets_DefaultIndexHtml_OK(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := &http.Client{}

	// Make sure the test cases don't exceed the rate limit
	testCases := []struct {
		name string
		path string
	}{
		{name: "root path", path: "/"},
		{name: "nested path", path: "/job_specs/abc123"},
		{name: "potentially valid path", path: "/valid/route"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := client.Get(app.Server.URL + tc.path)
			require.NoError(t, err)
			cltest.AssertServerResponse(t, resp, http.StatusOK)
		})
	}
}

func TestGuiAssets_DefaultIndexHtml_NotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := &http.Client{}

	// Make sure the test cases don't exceed the rate limit
	testCases := []struct {
		name string
		path string
	}{
		{name: "with extension", path: "/invalidFile.json"},
		{name: "nested path with extension", path: "/another/invalidFile.css"},
		{name: "bad api route", path: "/v2/bad/route"},
		{name: "non existent api version", path: "/v3/new/api/version"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := client.Get(app.Server.URL + tc.path)
			require.NoError(t, err)
			cltest.AssertServerResponse(t, resp, http.StatusNotFound)
		})
	}
}

func TestGuiAssets_DefaultIndexHtml_RateLimited(t *testing.T) {
	t.Parallel()

	config, cfgCleanup := cltest.NewConfig(t)
	config.Set("CHAINLINK_DEV", false)
	t.Cleanup(cfgCleanup)
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := &http.Client{}

	// Make calls equal to the rate limit
	rateLimit := 20
	for i := 0; i < rateLimit; i++ {
		resp, err := client.Get(app.Server.URL + "/")
		require.NoError(t, err)
		cltest.AssertServerResponse(t, resp, http.StatusOK)
	}

	// Last request fails
	resp, err := client.Get(app.Server.URL + "/")
	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}

func TestGuiAssets_AssetsExact(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := &http.Client{}

	resp, err := client.Get(app.Server.URL + "/assets/main.js")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	resp, err = client.Get(app.Server.URL + "/assets/mmain.js")
	require.NoError(t, err)
	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}

func TestGuiAssets_AssetsExactCompressed(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := &http.Client{}
	req, err := http.NewRequest("GET", app.Server.URL+"/assets/main.js", nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := client.Do(req)
	require.NoError(t, err)

	cltest.AssertServerResponse(t, resp, http.StatusOK)
	assert.Equal(t, "gzip", resp.Header["Content-Encoding"][0])
	assert.Equal(t, "Accept-Encoding", resp.Header["Vary"][0])

	req, err = http.NewRequest("GET", app.Server.URL+"/assets/doesnotexist.js", nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err = client.Do(req)
	require.NoError(t, err)

	cltest.AssertServerResponse(t, resp, http.StatusNotFound)
}
