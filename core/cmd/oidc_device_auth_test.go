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
	denyReason   string
	// verificationURI overrides the URI the start endpoint returns.
	verificationURI string
	// startStatus, when set, is returned by the start endpoint instead of 200.
	startStatus int
}

func (m *mockNode) handler(t *testing.T) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/oidc-enabled", func(w http.ResponseWriter, r *http.Request) {
		if !m.oidcEnabled {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"enabled": true})
	})
	mux.HandleFunc("/oidc-device/start", func(w http.ResponseWriter, r *http.Request) {
		// assert (not require) inside handler goroutines: require calls Goexit
		// on the wrong goroutine.
		assert.Equal(t, http.MethodPost, r.Method)
		if m.startStatus != 0 {
			w.WriteHeader(m.startStatus)
			_, _ = w.Write([]byte(`{"error":"node text"}`))
			return
		}
		uri := "https://sso.example.com/activate"
		if m.verificationURI != "" {
			uri = m.verificationURI
		}
		// An older node may still send verification_uri_complete. The CLI
		// must ignore it.
		_ = json.NewEncoder(w).Encode(map[string]any{
			"device_handle":             "test-handle",
			"user_code":                 "WDJB-MJHT",
			"verification_uri":          uri,
			"verification_uri_complete": uri + "?user_code=WDJB-MJHT",
			"expires_in":                300,
			"interval":                  1,
		})
	})
	mux.HandleFunc("/oidc-device/poll", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]string
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "test-handle", req["device_handle"], "CLI must echo back the opaque handle")

		if atomic.AddInt32(&m.pollsPending, -1) >= 0 {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
			return
		}
		if m.denyAfter {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "denied", "reason": m.denyReason, "message": "node text"})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "clsession",
			Value:    "abc123",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
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

	// HTTP 200 without {"enabled":true} must not trigger the device flow.
	bare200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer bare200.Close()
	assert.False(t, cmd.NodeHasOIDCEnabled(context.Background(), newClientOpts(t, bare200), lggr))
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

	// The operator must be shown the verification URI, the code, and a warning
	// to confirm the host before approving.
	assert.Contains(t, out.String(), "https://sso.example.com/activate")
	assert.Contains(t, out.String(), "WDJB-MJHT")
	assert.Contains(t, out.String(), "An approval you did not start is an attack.")
	// The pre-filled one-click link must never be shown.
	assert.NotContains(t, out.String(), "user_code=")
	assert.Contains(t, out.String(), "When prompted, enter the code: WDJB-MJHT")
}

func TestOIDCDeviceLogin_RefusesPlaintextNode(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{"http://node.internal:6688", "http://10.0.0.5:6688", "ftp://node.internal", "https://"} {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()
			u, err := url.Parse(raw)
			require.NoError(t, err)
			auth := cmd.NewOIDCDeviceCookieAuthenticator(cmd.ClientOpts{RemoteNodeURL: *u}, &cmd.MemoryCookieStore{}, &bytes.Buffer{}, logger.TestLogger(t))
			err = auth.Login(context.Background())
			require.ErrorContains(t, err, "requires an https remote-node-url")
		})
	}
}

func TestOIDCDeviceLogin_AllowsLoopbackHTTP(t *testing.T) {
	t.Parallel()
	// httptest serves http://127.0.0.1, so the success test covers 127.0.0.1.
	// localhost and [::1] are covered by rejecting at a later step instead of
	// the scheme check.
	for _, raw := range []string{"http://localhost:1", "http://[::1]:1"} {
		u, err := url.Parse(raw)
		require.NoError(t, err)
		auth := cmd.NewOIDCDeviceCookieAuthenticator(cmd.ClientOpts{RemoteNodeURL: *u}, &cmd.MemoryCookieStore{}, &bytes.Buffer{}, logger.TestLogger(t))
		err = auth.Login(context.Background())
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "requires an https remote-node-url", raw)
	}
}

func TestOIDCDeviceLogin_RefusesPlaintextVerificationURI(t *testing.T) {
	t.Parallel()
	for _, uri := range []string{"http://sso.example.com/activate", "javascript:alert(1)", "sso.example.com/activate"} {
		t.Run(uri, func(t *testing.T) {
			t.Parallel()
			node := &mockNode{oidcEnabled: true, verificationURI: uri}
			srv := httptest.NewServer(node.handler(t))
			defer srv.Close()
			var out bytes.Buffer
			auth := cmd.NewOIDCDeviceCookieAuthenticator(newClientOpts(t, srv), &cmd.MemoryCookieStore{}, &out, logger.TestLogger(t))
			err := auth.Login(context.Background())
			require.ErrorContains(t, err, "verification URI that is not https")
			assert.NotContains(t, out.String(), "sso.example.com", "a rejected URI must not be printed")
		})
	}
}

func TestOIDCDeviceLogin_NetworkRefused(t *testing.T) {
	t.Parallel()
	node := &mockNode{oidcEnabled: true, startStatus: http.StatusForbidden}
	srv := httptest.NewServer(node.handler(t))
	defer srv.Close()
	auth := cmd.NewOIDCDeviceCookieAuthenticator(newClientOpts(t, srv), &cmd.MemoryCookieStore{}, &bytes.Buffer{}, logger.TestLogger(t))
	err := auth.Login(context.Background())
	require.ErrorContains(t, err, "does not allow device login from your network")
	assert.NotContains(t, err.Error(), "node text")
}

func TestOIDCDeviceLogin_RefusesInsecureSkipVerify(t *testing.T) {
	t.Parallel()
	opts := cmd.ClientOpts{InsecureSkipVerify: true}
	auth := cmd.NewOIDCDeviceCookieAuthenticator(opts, &cmd.MemoryCookieStore{}, &bytes.Buffer{}, logger.TestLogger(t))
	err := auth.Login(context.Background())
	require.EqualError(t, err, "device login requires TLS verification; do not set insecure-skip-verify")
}

func TestOIDCDeviceLogin_Denied(t *testing.T) {
	t.Parallel()
	cases := []struct {
		reason string
		want   string
	}{
		{"denied", "denied at the identity provider"},
		{"no_role", "none of your identity provider groups map to a node role"},
		{"expired", "expired before it was approved"},
		{"", "was not completed"},
		{"something-new", "was not completed"},
	}
	for _, tc := range cases {
		t.Run(tc.reason, func(t *testing.T) {
			t.Parallel()
			node := &mockNode{oidcEnabled: true, pollsPending: 1, denyAfter: true, denyReason: tc.reason}
			srv := httptest.NewServer(node.handler(t))
			defer srv.Close()

			store := &cmd.MemoryCookieStore{}
			var out bytes.Buffer
			auth := cmd.NewOIDCDeviceCookieAuthenticator(newClientOpts(t, srv), store, &out, logger.TestLogger(t))

			err := auth.Login(context.Background())
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
			assert.NotContains(t, err.Error(), "node text", "CLI must not print the node's message")

			// On denial nothing must be written to the cookie jar.
			cookie, _ := store.Retrieve()
			assert.Nil(t, cookie)
		})
	}
}
