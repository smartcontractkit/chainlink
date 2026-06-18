package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// mockNode stands in for an OIDC-enabled chainlink node, serving the three
// endpoints the CLI device flow touches.
type mockNode struct {
	oidcEnabled  bool
	pollsPending int32 // number of "pending" polls to return before "complete"
	denyAfter    bool  // if true, return "denied" instead of completing
}

func (m *mockNode) handler(t *testing.T) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/oidc-enabled", func(w http.ResponseWriter, r *http.Request) {
		if !m.oidcEnabled {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/oidc-device/start", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"device_handle":             "test-handle",
			"user_code":                 "WDJB-MJHT",
			"verification_uri":          "https://sso.example.com/activate",
			"verification_uri_complete": "https://sso.example.com/activate?user_code=WDJB-MJHT",
			"expires_in":                300,
			"interval":                  1,
		})
	})
	mux.HandleFunc("/oidc-device/poll", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "test-handle", req["device_handle"], "CLI must echo back the opaque handle")

		if atomic.AddInt32(&m.pollsPending, -1) >= 0 {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
			return
		}
		if m.denyAfter {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "denied", "message": "expired"})
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "clsession", Value: "abc123", Path: "/"})
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "complete"})
	})
	return mux
}

func newClientOpts(t *testing.T, srv *httptest.Server) cmd.ClientOpts {
	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	return cmd.ClientOpts{RemoteNodeURL: *u}
}

func TestNodeHasOIDCEnabled(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	enabled := &mockNode{oidcEnabled: true}
	srv := httptest.NewServer(enabled.handler(t))
	defer srv.Close()
	assert.True(t, cmd.NodeHasOIDCEnabled(context.Background(), newClientOpts(t, srv), lggr))

	disabled := &mockNode{oidcEnabled: false}
	srv2 := httptest.NewServer(disabled.handler(t))
	defer srv2.Close()
	assert.False(t, cmd.NodeHasOIDCEnabled(context.Background(), newClientOpts(t, srv2), lggr))
}

func TestOIDCDeviceLogin_Success(t *testing.T) {
	t.Parallel()
	node := &mockNode{oidcEnabled: true, pollsPending: 2}
	srv := httptest.NewServer(node.handler(t))
	defer srv.Close()

	store := &cmd.MemoryCookieStore{}
	var out bytes.Buffer
	auth := cmd.NewOIDCDeviceCookieAuthenticator(newClientOpts(t, srv), store, &out, logger.TestLogger(t))

	require.NoError(t, auth.Login(context.Background()))

	// The session cookie produced by the node must land in the cookie jar.
	cookie, err := store.Retrieve()
	require.NoError(t, err)
	require.NotNil(t, cookie)
	assert.Equal(t, "clsession", cookie.Name)
	assert.Equal(t, "abc123", cookie.Value)

	// The operator must be shown the verification URI and code.
	assert.Contains(t, out.String(), "https://sso.example.com/activate")
	assert.Contains(t, out.String(), "WDJB-MJHT")
}

func TestOIDCDeviceLogin_Denied(t *testing.T) {
	t.Parallel()
	node := &mockNode{oidcEnabled: true, pollsPending: 1, denyAfter: true}
	srv := httptest.NewServer(node.handler(t))
	defer srv.Close()

	store := &cmd.MemoryCookieStore{}
	var out bytes.Buffer
	auth := cmd.NewOIDCDeviceCookieAuthenticator(newClientOpts(t, srv), store, &out, logger.TestLogger(t))

	err := auth.Login(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "denied")

	// On denial nothing must be written to the cookie jar.
	cookie, _ := store.Retrieve()
	assert.Nil(t, cookie)
}
