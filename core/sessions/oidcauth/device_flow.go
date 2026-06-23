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

const (
	// maxConcurrentDeviceFlows bounds the in-memory store so unauthenticated
	// callers cannot exhaust memory or fan out provider load without limit.
	maxConcurrentDeviceFlows = 100
	// deviceHandleBytes is the entropy of the opaque handle returned to the CLI.
	deviceHandleBytes = 32 // 256 bits
)

var (
	errTooManyDeviceFlows = errors.New("too many concurrent device authorization flows in progress")
	errUnknownDeviceFlow  = errors.New("unknown or expired device authorization flow")
)

// deviceFlowState is the per-flow state held server-side while a device
// authorization is pending. Exactly one background goroutine writes the
// terminal fields (sessionID/email/err/done); readers take the mutex.
type deviceFlowState struct {
	mu        sync.Mutex
	expiresAt time.Time
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

// add stores a flow under a fresh handle, enforcing the concurrency cap. The
// caller must not have logged or persisted the provider device_code.
func (s *deviceFlowStore) add(handle string, f *deviceFlowState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sweepLocked(time.Now())
	if len(s.flows) >= maxConcurrentDeviceFlows {
		return errTooManyDeviceFlows
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
func (s *deviceFlowStore) consumeIfDone(handle string) (sessionID, email string, role clsessions.UserRole, flowErr error, terminal bool, known bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, ok := s.flows[handle]
	if !ok {
		return "", "", "", nil, false, false
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.done {
		return "", "", "", nil, false, true
	}
	delete(s.flows, handle)
	return f.sessionID, f.email, f.role, f.err, true, true
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

// handleDeviceStart begins a device authorization flow. Unauthenticated, like
// the browser /oidc-login route: it only initiates the provider handshake.
func (oi *oidcAuthenticator) handleDeviceStart(c *gin.Context) {
	ctx := context.Background()
	da, err := oi.oauth2Config.DeviceAuth(ctx)
	if err != nil {
		oi.lggr.Errorf("device authorization request failed: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to start device authorization"})
		return
	}

	handle, err := oi.generateDeviceHandle()
	if err != nil {
		oi.lggr.Errorf("failed to generate device handle: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start device authorization"})
		return
	}

	expiry := da.Expiry
	if expiry.IsZero() {
		// Per RFC 8628 the provider should send expires_in; fall back to a
		// conservative default so a malformed response cannot create an
		// unbounded flow.
		expiry = time.Now().Add(5 * time.Minute)
	}
	state := &deviceFlowState{expiresAt: expiry}

	if err := oi.deviceFlows.add(handle, state); err != nil {
		oi.lggr.Warnf("refusing new device flow: %v", err)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

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

	sessionID, email, role, err := oi.issueSessionFromIDToken(ctx, rawIDToken)
	oi.finishDeviceFlow(state, sessionID, email, role, err)
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
	sessionID, email, role, flowErr, terminal, known := oi.deviceFlows.consumeIfDone(req.DeviceHandle)
	if !known {
		c.JSON(http.StatusNotFound, DevicePollResponse{Status: "denied", Message: errUnknownDeviceFlow.Error()})
		return
	}

	if !terminal {
		c.JSON(http.StatusOK, DevicePollResponse{Status: "pending"})
		return
	}

	if flowErr != nil {
		c.JSON(http.StatusOK, DevicePollResponse{Status: "denied", Message: "Authorization was not completed"})
		return
	}

	// Bind the session ID to the CLI's cookie jar, exactly as the browser flow
	// does on a successful exchange.
	ginSession := sessions.Default(c)
	ginSession.Set(webauth.SessionIDKey, sessionID)
	if err := ginSession.Save(); err != nil {
		oi.lggr.Errorf("failed to save device flow session: %v", err)
		c.JSON(http.StatusInternalServerError, DevicePollResponse{Status: "denied", Message: "Failed to establish session"})
		return
	}

	oi.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]any{"email": email, "role": role, "method": "device_flow"})
	c.JSON(http.StatusOK, DevicePollResponse{Status: "complete"})
}
