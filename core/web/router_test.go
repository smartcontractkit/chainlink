package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenAuthRequired_NoCredentials(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/v2/specs/", web.MediaType, bytes.NewBufferString("{}"))
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTokenAuthRequired_SessionCredentials(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/specs/", nil)
	defer cleanup()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTokenAuthRequired_TokenCredentials(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	eia := auth.NewToken()
	url := cltest.WebURL(t, "http://localhost:8888")
	eir := &models.ExternalInitiatorRequest{
		Name: "bitcoin",
		URL:  &url,
	}
	ea, err := models.NewExternalInitiator(eia, eir)
	require.NoError(t, err)
	err = app.GetStore().CreateExternalInitiator(ea)
	require.NoError(t, err)

	request, err := http.NewRequest("GET", ts.URL+"/v2/ping/", bytes.NewBufferString("{}"))
	require.NoError(t, err)
	request.Header.Set("Content-Type", web.MediaType)
	request.Header.Set("X-Chainlink-EA-AccessKey", eia.AccessKey)
	request.Header.Set("X-Chainlink-EA-Secret", eia.Secret)

	client := http.Client{}
	resp, err := client.Do(request)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTokenAuthRequired_BadTokenCredentials(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	eia := auth.NewToken()
	url := cltest.WebURL(t, "http://localhost:8888")
	eir := &models.ExternalInitiatorRequest{
		Name: "bitcoin",
		URL:  &url,
	}
	ea, err := models.NewExternalInitiator(eia, eir)
	require.NoError(t, err)
	err = app.GetStore().CreateExternalInitiator(ea)
	require.NoError(t, err)

	request, err := http.NewRequest("GET", ts.URL+"/v2/ping/", bytes.NewBufferString("{}"))
	require.NoError(t, err)
	request.Header.Set("Content-Type", web.MediaType)
	request.Header.Set("X-Chainlink-EA-AccessKey", eia.AccessKey)
	request.Header.Set("X-Chainlink-EA-Secret", "every unpleasant commercial color from aquamarine to beige")

	client := http.Client{}
	resp, err := client.Do(request)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSessions_RateLimited(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := http.Client{}
	input := `{"email":"brute@force.com", "password": "wrongpassword"}`

	for i := 0; i < 5; i++ {
		request, err := http.NewRequest("POST", ts.URL+"/sessions", bytes.NewBufferString(input))
		require.NoError(t, err)

		resp, err := client.Do(request)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	}

	request, err := http.NewRequest("POST", ts.URL+"/sessions", bytes.NewBufferString(input))
	require.NoError(t, err)

	resp, err := client.Do(request)
	require.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode)
}

func TestRouter_LargePOSTBody(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := http.Client{}

	body := string(make([]byte, 70000))
	request, err := http.NewRequest("POST", ts.URL+"/sessions", bytes.NewBufferString(body))
	require.NoError(t, err)

	resp, err := client.Do(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
}

func TestRouter_GinHelmetHeaders(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	router := web.Router(app)
	ts := httptest.NewServer(router)
	defer ts.Close()
	res, err := http.Get(ts.URL)
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
