package oidcauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"golang.org/x/oauth2"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
	webauth "github.com/smartcontractkit/chainlink/v2/core/web/auth"
)

// Device authorization grant (RFC 8628) for headless clients (the CLI).
//
// The node brokers the flow as a confidential or public OAuth client: it calls
// the provider's device authorization endpoint, hands the user-facing code and
// verification URI back to the CLI, then polls the provider's token endpoint on
// the CLI's behalf. The raw provider device_code never leaves the node; the CLI
// only ever holds an opaque, single-use handle minted here.
//
// Polling the provider runs in one background goroutine per flow, driven by
// oauth2's DeviceAccessToken which honours the server-mandated interval and
// slow_down backoff. The CLI's poll to the node is therefore decoupled from the
// node's poll to the provider: no matter how fast the CLI asks, the node never
// exceeds the provider's cadence.
//
// Multi-replica note: deviceFlowStore is process-local memory. Deployments with
// more than one replica behind a load balancer MUST use session affinity
// (sticky sessions) for /oidc-device/* (and preferably the whole node API), or
// replace this store with a shared short-TTL backend. Without affinity, start
// and poll can land on different replicas and the login fails closed.

const (
	// maxConcurrentDeviceFlows bounds the in-memory store so unauthenticated
	// callers cannot exhaust memory or fan out provider load without limit.
	maxConcurrentDeviceFlows = 100
	// maxDeviceFlowsPerIP stops a single source from consuming the entire
	// global budget (fairness / DoS).
	maxDeviceFlowsPerIP = 5
	// deviceHandleBytes is the entropy of the opaque handle returned to the CLI.
	deviceHandleBytes = 32 // 256 bits

	// Dedicated pre-auth rate limits for device endpoints. These sit on top of
	// the outer authenticated-tier limiter and match the unauth tier defaults
	// for start (expensive: outbound IdP call + goroutine). Poll is looser so a
	// legitimate CLI can tick at the provider interval.
	deviceStartRateLimit  = 5
	deviceStartRatePeriod = 20 * time.Second
	devicePollRateLimit   = 60
	devicePollRatePeriod  = time.Minute
)

var (
	errTooManyDeviceFlows      = errors.New("too many concurrent device authorization flows in progress")
	errTooManyDeviceFlowsPerIP = errors.New("too many concurrent device authorization flows from this address")
	errUnknownDeviceFlow       = errors.New("unknown or expired device authorization flow")
	errDeviceFlowAbandoned     = errors.New("device authorization flow expired or abandoned before completion")
)

// deviceFlowState is the per-flow state held server-side while a device
// authorization is pending. Exactly one background goroutine writes the
// terminal fields (sessionID/email/err/done); readers take the mutex.
type deviceFlowState struct {
	mu        sync.Mutex
	expiresAt time.Time
	clientIP  string
	done      bool
	sessionID string
	email     string
	role      clsessions.UserRole
	err       error
}

// deviceFlowStore maps opaque handles to pending flows. Handles are single-use:
// a successful poll consumes the flow. Expired flows are swept lazily.
type deviceFlowStore struct {
	mu    sync.Mutex
	flows map[string]*deviceFlowState
}

func newDeviceFlowStore() *deviceFlowStore {
	return &deviceFlowStore{flows: make(map[string]*deviceFlowState)}
}

func (s *deviceFlowStore) sweepLocked(now time.Time) {
	for h, f := range s.flows {
		f.mu.Lock()
		expired := now.After(f.expiresAt)
		f.mu.Unlock()
		if expired {
			delete(s.flows, h)
		}
	}
}

// add stores a flow under a fresh handle, enforcing the global and per-IP
// concurrency caps. The caller must not have logged or persisted the provider
// device_code. Prefer reserving the slot *before* calling the IdP so a refused
// flow never produces outbound provider load.
func (s *deviceFlowStore) add(handle string, f *deviceFlowState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sweepLocked(time.Now())
	if len(s.flows) >= maxConcurrentDeviceFlows {
		return errTooManyDeviceFlows
	}
	if f.clientIP != "" {
		n := 0
		for _, existing := range s.flows {
			if existing.clientIP == f.clientIP {
				n++
			}
		}
		if n >= maxDeviceFlowsPerIP {
			return errTooManyDeviceFlowsPerIP
		}
	}
	s.flows[handle] = f
	return nil
}

func (s *deviceFlowStore) get(handle string) (*deviceFlowState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, ok := s.flows[handle]
	return f, ok
}

// contains reports whether handle is still present (not swept/consumed).
func (s *deviceFlowStore) contains(handle string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.flows[handle]
	return ok
}

func (s *deviceFlowStore) remove(handle string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.flows, handle)
}

// consumeIfDone atomically reports whether the flow for handle has reached a
// terminal state and, if so, removes it in the same critical section. This
// closes the double-poll window: two concurrent polls on a completed handle
// cannot both observe done and both be issued the session cookie; exactly one
// sees terminal=true, the rest see (nil, false) once it is gone.
func (s *deviceFlowStore) consumeIfDone(handle string) (sessionID, email string, role clsessions.UserRole, terminal bool, known bool, flowErr error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, ok := s.flows[handle]
	if !ok {
		return "", "", "", false, false, nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.done {
		return "", "", "", false, true, nil
	}
	delete(s.flows, handle)
	return f.sessionID, f.email, f.role, true, true, f.err
}

// generateDeviceHandle returns a URL-safe 256-bit random handle. It is the only
// device-flow identifier the CLI ever sees; the provider device_code stays
// server-side.
func (oi *oidcAuthenticator) generateDeviceHandle() (string, error) {
	b := make([]byte, deviceHandleBytes)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// DeviceStartResponse is the CLI-facing payload that begins a device flow.
// It deliberately omits the provider device_code.
type DeviceStartResponse struct {
	DeviceHandle            string `json:"device_handle"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int64  `json:"expires_in"`
	Interval                int64  `json:"interval"`
}

// DevicePollRequest is the CLI's request to check on a pending flow.
type DevicePollRequest struct {
	DeviceHandle string `json:"device_handle" binding:"required"`
}

// DevicePollResponse reports flow status to the CLI. Status is one of
// "pending", "complete", or "denied". On "complete" the session cookie is set
// on the response; the body carries no token material.
type DevicePollResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// deviceEndpointRateLimit builds a dedicated per-IP limiter for pre-auth device
// routes. Nested under the outer api group so both limits apply.
func deviceEndpointRateLimit(period time.Duration, limit int64) gin.HandlerFunc {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}
	return mgin.NewMiddleware(limiter.New(store, rate))
}

// handleDeviceStart begins a device authorization flow. Unauthenticated, like
// the browser /oidc-login route: it only initiates the provider handshake.
func (oi *oidcAuthenticator) handleDeviceStart(c *gin.Context) {
	handle, err := oi.generateDeviceHandle()
	if err != nil {
		oi.lggr.Errorf("failed to generate device handle: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start device authorization"})
		return
	}

	// Reserve a concurrency slot *before* calling the IdP so a refused flow
	// never produces outbound provider load, and so a flood cannot open more
	// IdP device codes than the store will accept.
	state := &deviceFlowState{
		expiresAt: time.Now().Add(5 * time.Minute), // provisional; refined after IdP response
		clientIP:  c.ClientIP(),
	}
	if err := oi.deviceFlows.add(handle, state); err != nil {
		oi.lggr.Warnf("refusing new device flow: %v", err)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	da, err := oi.oauth2Config.DeviceAuth(ctx)
	if err != nil {
		oi.deviceFlows.remove(handle)
		oi.lggr.Errorf("device authorization request failed: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to start device authorization"})
		return
	}

	expiry := da.Expiry
	if expiry.IsZero() {
		// Per RFC 8628 the provider should send expires_in; fall back to a
		// conservative default so a malformed response cannot create an
		// unbounded flow.
		expiry = time.Now().Add(5 * time.Minute)
	}
	state.mu.Lock()
	state.expiresAt = expiry
	state.mu.Unlock()

	// Poll the provider in the background. DeviceAccessToken blocks until the
	// user approves, denies, or the code expires, internally respecting the
	// provider's interval and slow_down. Bounded by the device-code expiry.
	go oi.pollDeviceToken(handle, state, da)

	interval := da.Interval
	if interval == 0 {
		interval = 5 // RFC 8628 default
	}
	c.JSON(http.StatusOK, DeviceStartResponse{
		DeviceHandle:            handle,
		UserCode:                da.UserCode,
		VerificationURI:         da.VerificationURI,
		VerificationURIComplete: da.VerificationURIComplete,
		ExpiresIn:               int64(time.Until(expiry).Seconds()),
		Interval:                interval,
	})
}

// pollDeviceToken runs the blocking provider poll for one flow and records the
// terminal result. On success it verifies the id_token and creates a session
// through the same issueSessionFromIDToken path the browser flow uses.
func (oi *oidcAuthenticator) pollDeviceToken(handle string, state *deviceFlowState, da *oauth2.DeviceAuthResponse) {
	ctx, cancel := context.WithDeadline(context.Background(), state.expiresAt)
	defer cancel()

	token, err := oi.oauth2Config.DeviceAccessToken(ctx, da)
	if err != nil {
		oi.finishDeviceFlow(state, "", "", "", err)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		oi.finishDeviceFlow(state, "", "", "", errors.New("missing id_token in device token response"))
		return
	}

	// If the handle was swept/expired while we waited on the IdP, do not mint a
	// session nobody can claim (orphan oidc_sessions row).
	if !oi.deviceFlows.contains(handle) {
		oi.lggr.Warnf("device flow handle gone before token completion; abandoning without session")
		oi.finishDeviceFlow(state, "", "", "", errDeviceFlowAbandoned)
		return
	}

	// Device flow does not use a browser nonce; pass empty expectedNonce.
	sessionID, email, role, err := oi.issueSessionFromIDToken(ctx, rawIDToken, "")
	if err != nil {
		oi.finishDeviceFlow(state, "", "", "", err)
		return
	}

	// Race: handle may have been swept between contains and issue. Drop the
	// orphan session rather than leave a live unclaimable row.
	if !oi.deviceFlows.contains(handle) {
		oi.lggr.Warnf("device flow handle gone after session insert; deleting orphan session")
		oi.deleteOIDCSession(ctx, sessionID)
		oi.finishDeviceFlow(state, "", "", "", errDeviceFlowAbandoned)
		return
	}

	oi.finishDeviceFlow(state, sessionID, email, role, nil)
}

func (oi *oidcAuthenticator) finishDeviceFlow(state *deviceFlowState, sessionID, email string, role clsessions.UserRole, err error) {
	state.mu.Lock()
	state.done = true
	state.sessionID = sessionID
	state.email = email
	state.role = role
	state.err = err
	state.mu.Unlock()
	if err != nil {
		oi.lggr.Errorf("device authorization flow failed: %v", err)
	}
}

// deleteOIDCSession removes a session row that was minted but never bound to a
// client cookie (Save failure or abandoned device flow).
func (oi *oidcAuthenticator) deleteOIDCSession(ctx context.Context, sessionID string) {
	if sessionID == "" {
		return
	}
	if _, err := oi.ds.ExecContext(ctx, "DELETE FROM oidc_sessions WHERE id = $1", sessionID); err != nil {
		oi.lggr.Errorf("failed to delete orphan OIDC session %s: %v", sessionID, err)
	}
}

// handleDevicePoll reports the status of a pending flow to the CLI and, once
// complete, sets the session cookie. The handle is single-use: a completed flow
// is removed after its cookie is issued or its failure reported.
func (oi *oidcAuthenticator) handleDevicePoll(c *gin.Context) {
	var req DevicePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, DevicePollResponse{Status: "denied", Message: "Invalid request"})
		return
	}

	// Read the terminal result and consume the handle in one critical section
	// so two concurrent polls cannot both be handed the session cookie.
	sessionID, email, role, terminal, known, flowErr := oi.deviceFlows.consumeIfDone(req.DeviceHandle)
	if !known {
		c.JSON(http.StatusNotFound, DevicePollResponse{Status: "denied", Message: errUnknownDeviceFlow.Error()})
		return
	}

	if !terminal {
		c.JSON(http.StatusOK, DevicePollResponse{Status: "pending"})
		return
	}

	if flowErr != nil {
		// Terminal failure already consumed the handle; no session to clean up
		// unless a prior success path left one (should not happen).
		c.JSON(http.StatusOK, DevicePollResponse{Status: "denied", Message: "Authorization was not completed"})
		return
	}

	// Bind the session ID to the CLI's cookie jar, exactly as the browser flow
	// does on a successful exchange.
	ginSession := sessions.Default(c)
	ginSession.Set(webauth.SessionIDKey, sessionID)
	if err := ginSession.Save(); err != nil {
		// Handle is already consumed (single-winner). Delete the session row so
		// we do not leave a live unclaimable credential; the CLI must restart.
		oi.lggr.Errorf("failed to save device flow session: %v", err)
		oi.deleteOIDCSession(context.Background(), sessionID)
		c.JSON(http.StatusInternalServerError, DevicePollResponse{Status: "denied", Message: "Failed to establish session"})
		return
	}

	oi.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]any{"email": email, "role": role, "method": "device_flow"})
	c.JSON(http.StatusOK, DevicePollResponse{Status: "complete"})
}
