package oidcauth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
)

// These tests drive the real routes from ExtendRouter against the mock IdP
// and a real database.

type auditEvent struct {
	id   audit.EventID
	data audit.Data
}

// captureAudit records audit events. It embeds the no-op service for the
// services.Service half of audit.Logger.
type captureAudit struct {
	audit.AuditLoggerService
	mu     sync.Mutex
	events []auditEvent
}

func (c *captureAudit) Audit(id audit.EventID, data audit.Data) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, auditEvent{id: id, data: data})
}

func (c *captureAudit) find(id audit.EventID) (audit.Data, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, e := range c.events {
		if e.id == id {
			return e.data, true
		}
	}
	return nil, false
}

// idpConfig points ProviderURL at the mock IdP so the node's verification URI
// host check passes.
type idpConfig struct {
	*TestConfig
	provider string
}

func (c idpConfig) ProviderURL() string { return c.provider }

type routerHarness struct {
	t      *testing.T
	engine *gin.Engine
	oi     *oidcAuthenticator
	idp    *mockIDP
	audit  *captureAudit
}

func newRouterHarness(t *testing.T) *routerHarness {
	t.Helper()
	db := pgtest.NewSqlxDB(t)
	idp := newMockIDP(t, testClientID)
	idp.email = "operator@example.com"
	idp.groups = []string{EditorClaim}
	oi := newAuthenticatorForIDP(t, idp, db)
	capture := &captureAudit{}
	oi.auditLogger = capture
	oi.config = idpConfig{TestConfig: &TestConfig{}, provider: idp.server.URL}
	// Operators on the VPN (10.8.0.0/16) behind one load balancer (10.0.0.0/24).
	oi.deviceAllowed = mustCIDRs(t, "10.8.0.0/16")
	oi.clients = clientResolver{trustedProxies: mustCIDRs(t, "10.0.0.0/24")}

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := engine.Group("/", sessions.Sessions("clsession", cookie.NewStore([]byte("0123456789abcdef0123456789abcdef"))))
	require.NoError(t, oi.ExtendRouter(api))
	return &routerHarness{t: t, engine: engine, oi: oi, idp: idp, audit: capture}
}

// do sends a request as if it came through the load balancer for client.
func (h *routerHarness) do(method, path, client string, body any) *httptest.ResponseRecorder {
	h.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(h.t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequestWithContext(h.t.Context(), method, path, &buf)
	req.RemoteAddr = "10.0.0.2:40000"
	req.Header.Set("X-Forwarded-For", client)
	req.Header.Set("User-Agent", "chainlink-cli-test")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	h.engine.ServeHTTP(w, req)
	return w
}

func (h *routerHarness) start(client string) DeviceStartResponse {
	h.t.Helper()
	w := h.do(http.MethodPost, "/oidc-device/start", client, nil)
	require.Equal(h.t, http.StatusOK, w.Code, w.Body.String())
	assert.NotContains(h.t, w.Body.String(), "verification_uri_complete")
	assert.NotContains(h.t, w.Body.String(), "user_code=")
	var sr DeviceStartResponse
	require.NoError(h.t, json.Unmarshal(w.Body.Bytes(), &sr))
	return sr
}

func TestDeviceRoutes_OutsideAllowListRefused(t *testing.T) {
	t.Parallel()
	h := newRouterHarness(t)

	for _, path := range []string{"/oidc-device/start", "/oidc-device/poll"} {
		w := h.do(http.MethodPost, path, "203.0.113.5", map[string]string{"device_handle": "x"})
		assert.Equal(t, http.StatusForbidden, w.Code, path)
	}
	// Nothing reached the identity provider or the database.
	var n int
	require.NoError(t, h.oi.ds.GetContext(t.Context(), &n, `SELECT count(*) FROM oidc_device_flows`))
	assert.Equal(t, 0, n)
}

func TestDeviceRoutes_DisabledWhenAllowListEmpty(t *testing.T) {
	t.Parallel()
	h := newRouterHarness(t)
	h.oi.deviceAllowed = nil
	w := h.do(http.MethodPost, "/oidc-device/start", "10.8.1.1", nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "disabled")
}

// TestDeviceRoutes_EndToEndAudit runs a full device login through the router
// and checks the audit trail can join start and redemption.
func TestDeviceRoutes_EndToEndAudit(t *testing.T) {
	t.Parallel()
	h := newRouterHarness(t)

	sr := h.start("10.8.1.1")
	assert.Equal(t, "WDJB-MJHT", sr.UserCode)
	assert.Equal(t, h.idp.server.URL+"/activate", sr.VerificationURI)

	started, ok := h.audit.find(audit.AuthDeviceFlowStarted)
	require.True(t, ok, "start must be audited")
	assert.Equal(t, "10.8.1.1", started["client_ip"])
	assert.Equal(t, "chainlink-cli-test", started["user_agent"])
	assert.Equal(t, handleFingerprint(sr.DeviceHandle), started["handle_fp"])
	assert.NotContains(t, started, "device_handle")

	// Redeem from a different VPN address, as a phished approval would.
	var w *httptest.ResponseRecorder
	require.Eventually(t, func() bool {
		w = h.do(http.MethodPost, "/oidc-device/poll", "10.8.9.9", map[string]string{"device_handle": sr.DeviceHandle})
		return w.Code == http.StatusOK && bytes.Contains(w.Body.Bytes(), []byte(`"complete"`))
	}, 10*time.Second, 100*time.Millisecond)
	assert.NotEmpty(t, w.Result().Cookies(), "the session cookie is set on completion")

	success, ok := h.audit.find(audit.AuthLoginSuccessNo2FA)
	require.True(t, ok)
	assert.Equal(t, "operator@example.com", success["email"])
	assert.Equal(t, started["handle_fp"], success["handle_fp"], "start and success join on the fingerprint")
	assert.Equal(t, "10.8.1.1", success["start_client_ip"])
	assert.Equal(t, "10.8.9.9", success["client_ip"])
	assert.Equal(t, true, success["client_ip_mismatch"])

	// Single use.
	w = h.do(http.MethodPost, "/oidc-device/poll", "10.8.9.9", map[string]string{"device_handle": sr.DeviceHandle})
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestDeviceRoutes_AbandonedFlowStopsProviderPolling: a flow nobody polls
// gives up its slot and stops calling the identity provider.
func TestDeviceRoutes_AbandonedFlowStopsProviderPolling(t *testing.T) {
	t.Parallel()
	h := newRouterHarness(t)
	h.oi.liveCheckEvery = 100 * time.Millisecond
	h.idp.pending.Store(true)

	sr := h.start("10.8.1.1")
	require.True(t, h.oi.deviceFlows.contains(sr.DeviceHandle))

	// Simulate the poll grace running out.
	_, err := h.oi.ds.ExecContext(t.Context(), `UPDATE oidc_device_flows SET poll_deadline = now() - interval '1 second' WHERE handle = $1`, sr.DeviceHandle)
	require.NoError(t, err)

	require.Eventually(t, func() bool { return !h.oi.deviceFlows.contains(sr.DeviceHandle) }, 5*time.Second, 50*time.Millisecond,
		"the poller must exit and free the replica slot")
	failed, ok := h.audit.find(audit.AuthLoginFailedDevice)
	require.True(t, ok)
	assert.Equal(t, "abandoned", failed["reason"])
	assert.Equal(t, "10.8.1.1", failed["start_client_ip"])
	assert.Equal(t, handleFingerprint(sr.DeviceHandle), failed["handle_fp"])

	calls := h.idp.tokenCalls.Load()
	time.Sleep(1500 * time.Millisecond)
	assert.Equal(t, calls, h.idp.tokenCalls.Load(), "no provider calls after abandonment")
}

func TestSignInRoute_RateLimitedPerClient(t *testing.T) {
	t.Parallel()
	h := newRouterHarness(t)
	for i := range signInRateLimit {
		w := h.do(http.MethodGet, "/oidc-login", "198.51.100.7", nil)
		require.Equal(t, http.StatusFound, w.Code, "request %d", i)
	}
	w := h.do(http.MethodGet, "/oidc-login", "198.51.100.7", nil)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	// Another client is unaffected.
	w = h.do(http.MethodGet, "/oidc-login", "198.51.100.8", nil)
	assert.Equal(t, http.StatusFound, w.Code)

	var n int
	require.NoError(t, h.oi.ds.GetContext(t.Context(), &n, `SELECT count(*) FROM oidc_pending_auth WHERE client_key = '198.51.100.7'`))
	assert.Equal(t, signInRateLimit, n, "refused sign-ins write nothing")
}
