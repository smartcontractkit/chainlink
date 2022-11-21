package web_test

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	clhttptest "github.com/smartcontractkit/chainlink/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/operator_ui/assets
var testFs embed.FS

func TestGuiAssets_DefaultIndexHtml_OK(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

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

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))

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

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.DevMode = false
	})
	app := cltest.NewApplicationWithConfig(t, config)
	require.NoError(t, app.Start(testutils.Context(t)))

	client := clhttptest.NewTestLocalOnlyHTTPClient()

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

func TestGuiAssets_AssetsFS(t *testing.T) {
	t.Parallel()

	efs := web.NewEmbedFileSystem(testFs, "fixtures/operator_ui")
	handler := web.ServeGzippedAssets("/fixtures/operator_ui/", efs, logger.TestLogger(t))

	t.Run("it get exact assets if Accept-Encoding is not specified", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		var err error
		c.Request, err = http.NewRequest("GET", "http://localhost:6688/fixtures/operator_ui/assets/main.js", nil)
		require.NoError(t, err)
		handler(c)

		require.Equal(t, http.StatusOK, recorder.Result().StatusCode)

		recorder = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(recorder)
		c.Request, err = http.NewRequest("GET", "http://localhost:6688/fixtures/operator_ui/assets/kinda_main.js", nil)
		require.NoError(t, err)
		handler(c)

		require.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
	})

	t.Run("it respects Accept-Encoding header", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		var err error
		c.Request, err = http.NewRequest("GET", "http://localhost:6688/fixtures/operator_ui/assets/main.js", nil)
		require.NoError(t, err)
		c.Request.Header.Set("Accept-Encoding", "gzip")
		handler(c)

		require.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		require.Equal(t, "gzip", recorder.Result().Header.Get("Content-Encoding"))

		recorder = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(recorder)
		c.Request, err = http.NewRequest("GET", "http://localhost:6688/fixtures/operator_ui/assets/kinda_main.js", nil)
		require.NoError(t, err)
		c.Request.Header.Set("Accept-Encoding", "gzip")
		handler(c)

		require.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
	})
}
