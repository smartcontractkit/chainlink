package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenAuthRequired_NoCredentials(t *testing.T) {
	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", ts.URL+"/v2/jobs/", bytes.NewBufferString("{}"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", web.MediaType)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTokenAuthRequired_SessionCredentials(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Post("/v2/bridge_types/", nil)
	defer cleanup()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestTokenAuthRequired_TokenCredentials(t *testing.T) {
	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	eia := auth.NewToken()
	url := cltest.WebURL(t, "http://localhost:8888")
	eir := &bridges.ExternalInitiatorRequest{
		Name: uuid.New().String(),
		URL:  &url,
	}
	ea, err := bridges.NewExternalInitiator(eia, eir)
	require.NoError(t, err)
	err = app.BridgeORM().CreateExternalInitiator(ctx, ea)
	require.NoError(t, err)

	request, err := http.NewRequestWithContext(ctx, "GET", ts.URL+"/v2/ping/", bytes.NewBufferString("{}"))
	require.NoError(t, err)
	request.Header.Set("Content-Type", web.MediaType)
	request.Header.Set("X-Chainlink-EA-AccessKey", eia.AccessKey)
	request.Header.Set("X-Chainlink-EA-Secret", eia.Secret)

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	resp, err := client.Do(request)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTokenAuthRequired_BadTokenCredentials(t *testing.T) {
	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	eia := auth.NewToken()
	url := cltest.WebURL(t, "http://localhost:8888")
	eir := &bridges.ExternalInitiatorRequest{
		Name: uuid.New().String(),
		URL:  &url,
	}
	ea, err := bridges.NewExternalInitiator(eia, eir)
	require.NoError(t, err)
	err = app.BridgeORM().CreateExternalInitiator(ctx, ea)
	require.NoError(t, err)

	request, err := http.NewRequestWithContext(ctx, "GET", ts.URL+"/v2/ping/", bytes.NewBufferString("{}"))
	require.NoError(t, err)
	request.Header.Set("Content-Type", web.MediaType)
	request.Header.Set("X-Chainlink-EA-AccessKey", eia.AccessKey)
	request.Header.Set("X-Chainlink-EA-Secret", "every unpleasant commercial color from aquamarine to beige")

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	resp, err := client.Do(request)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSessions_RateLimited(t *testing.T) {
	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := clhttptest.NewTestLocalOnlyHTTPClient()
	input := `{"email":"brute@force.com", "password": "wrongpassword"}`

	for i := 0; i < 5; i++ {
		request, err := http.NewRequestWithContext(ctx, "POST", ts.URL+"/sessions", bytes.NewBufferString(input))
		require.NoError(t, err)

		resp, err := client.Do(request)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", ts.URL+"/sessions", bytes.NewBufferString(input))
	require.NoError(t, err)

	resp, err := client.Do(request)
	require.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode)
}

func TestRouter_LargePOSTBody(t *testing.T) {
	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := clhttptest.NewTestLocalOnlyHTTPClient()

	body := string(make([]byte, 70000))
	request, err := http.NewRequestWithContext(ctx, "POST", ts.URL+"/sessions", bytes.NewBufferString(body))
	require.NoError(t, err)

	resp, err := client.Do(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
}

func TestRouter_GinHelmetHeaders(t *testing.T) {
	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	require.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	for _, tt := range []struct {
		HelmetName  string
		HeaderKey   string
		HeaderValue string
	}{
		{"NoSniff", "X-Content-Type-Options", "nosniff"},
		{"DNSPrefetchControl", "X-DNS-Prefetch-Control", "off"},
		{"FrameGuard", "X-Frame-Options", "DENY"},
		{"SetHSTS", "Strict-Transport-Security", "max-age=5184000; includeSubDomains"},
		{"IENoOpen", "X-Download-Options", "noopen"},
		{"XSSFilter", "X-Xss-Protection", "1; mode=block"},
	} {
		assert.Equal(t, res.Header.Get(tt.HeaderKey), tt.HeaderValue,
			"wrong header for helmet's %s handler", tt.HelmetName)
	}
}
