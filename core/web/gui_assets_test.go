package web_test

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

//go:embed fixtures/operator_ui/assets
var testFs embed.FS

func TestGuiAssets_DefaultIndexHtml_OK(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(t.Context()))

	client := clhttptest.NewTestLocalOnlyHTTPClient()

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
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(), "GET", app.Server.URL+tc.path, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			cltest.AssertServerResponse(t, resp, http.StatusOK)
		})
	}
}

func TestGuiAssets_DefaultIndexHtml_NotFound(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(t.Context()))

	client := clhttptest.NewTestLocalOnlyHTTPClient()

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
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(t.Context(), "GET", app.Server.URL+tc.path, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			cltest.AssertServerResponse(t, resp, http.StatusNotFound)
		})
	}
}

func TestGuiAssets_DefaultIndexHtml_RateLimited(t *testing.T) {
	t.Parallel()

	config := configtest.NewGeneralConfig(t, nil)
	app := cltest.NewApplicationWithConfig(t, config)
	require.NoError(t, app.Start(t.Context()))

	client := clhttptest.NewTestLocalOnlyHTTPClient()

	// Make calls equal to the rate limit
	rateLimit := 20
	for range rateLimit {
		req, err := http.NewRequestWithContext(t.Context(), "GET", app.Server.URL+"/", nil)
		require.NoError(t, err)
		resp, err := client.Do(req)
		require.NoError(t, err)
		cltest.AssertServerResponse(t, resp, http.StatusOK)
	}

	// Last request fails
	req, err := http.NewRequestWithContext(t.Context(), "GET", app.Server.URL+"/", nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}

func TestGuiAssets_AssetsFS(t *testing.T) {
	t.Parallel()

	efs := web.NewEmbedFileSystem(testFs, "fixtures/operator_ui")
	handler := web.ServeGzippedAssets("/fixtures/operator_ui/", efs, logger.TestLogger(t))

	t.Run("it get exact assets if Accept-Encoding is not specified", func(t *testing.T) {
		t.Parallel()
		engine := gin.New()
		engine.GET("/fixtures/operator_ui/*filepath", handler)

		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(t.Context(), "GET", "/fixtures/operator_ui/assets/main.js", nil)
		require.NoError(t, err)
		engine.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Result().StatusCode)

		recorder = httptest.NewRecorder()
		req, err = http.NewRequestWithContext(t.Context(), "GET", "/fixtures/operator_ui/assets/kinda_main.js", nil)
		require.NoError(t, err)
		engine.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
	})

	t.Run("it respects Accept-Encoding header", func(t *testing.T) {
		t.Parallel()
		engine := gin.New()
		engine.GET("/fixtures/operator_ui/*filepath", handler)

		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(t.Context(), "GET", "/fixtures/operator_ui/assets/main.js", nil)
		require.NoError(t, err)
		req.Header.Set("Accept-Encoding", "gzip")
		engine.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		require.Equal(t, "gzip", recorder.Result().Header.Get("Content-Encoding"))

		recorder = httptest.NewRecorder()
		req, err = http.NewRequestWithContext(t.Context(), "GET", "/fixtures/operator_ui/assets/kinda_main.js", nil)
		require.NoError(t, err)
		req.Header.Set("Accept-Encoding", "gzip")
		engine.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
	})
}
