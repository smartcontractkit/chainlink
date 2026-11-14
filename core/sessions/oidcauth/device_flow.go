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
// Multi-replica: the opaque handle and its terminal result live in
// oidc_device_flows, so start and poll may land on different replicas. The
// provider device code stays in the poller process and is not written to the
// database. The in-memory map below is only the poller's local copy.

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

	// Browser sign-in writes one oidc_pending_auth row per call.
	signInRateLimit  = 20
	signInRatePeriod = time.Minute

	// minDevicePollGrace is how long a flow may go unpolled before it stops
	// holding a slot. The grace is the larger of this and three intervals.
	minDevicePollGrace = 30 * time.Second
	// deviceLiveCheckEvery is how often the provider poller checks that the
	// CLI is still polling, so an abandoned flow stops calling the provider.
	deviceLiveCheckEvery = 10 * time.Second
)

var (
	errTooManyDeviceFlows      = errors.New("too many concurrent device authorization flows in progress")
	errTooManyDeviceFlowsPerIP = errors.New("too many concurrent device authorization flows from this address")
	errUnknownDeviceFlow       = errors.New("unknown or expired device authorization flow")
	errDeviceLoginDisabled     = errors.New("device login is disabled on this node")
	errDeviceLoginNetwork      = errors.New("device login is not allowed from this network")
	errDeviceFlowAbandoned     = errors.New("device authorization flow expired or abandoned before completion")
)

// devicePollGrace is how long a flow may go unpolled before it is dropped.
func devicePollGrace(interval int64) time.Duration {
	return max(minDevicePollGrace, 3*time.Duration(interval)*time.Second)
}

// deviceFlowState is the per-flow state held server-side while a device
// authorization is pending. Exactly one background goroutine writes the
// terminal fields (sessionID/email/err/done); readers take the mutex.
type deviceFlowState struct {
	mu        sync.Mutex
	expiresAt time.Time
	// clientIP is the address that started the flow, for audit.
	clientIP string
	// bucket is the rate-limit key the per-client cap counts.
	bucket    string
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
	n := 0
	for _, existing := range s.flows {
		if existing.bucket == f.bucket {
			n++
		}
	}
	if n >= maxDeviceFlowsPerIP {
		return errTooManyDeviceFlowsPerIP
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
// It deliberately omits the provider device_code. It also omits the
// provider's verification_uri_complete: a link with the code already filled
// in turns a forwarded login into a one-click approval (RFC 8628 section
// 5.4), so the operator has to type or paste the user code.
type DeviceStartResponse struct {
	DeviceHandle    string `json:"device_handle"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}

// DevicePollRequest is the CLI's request to check on a pending flow.
type DevicePollRequest struct {
	DeviceHandle string `json:"device_handle" binding:"required"`
}

// DevicePollResponse reports flow status to the CLI. Status is one of
// "pending", "complete", or "denied". On "denied" after a finished flow,
// Reason is one of the deviceReasonMessages keys. On "complete" the session
// cookie is set on the response; the body carries no token material.
type DevicePollResponse struct {
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// preAuthRateLimit builds a dedicated limiter for a pre-auth OIDC route.
// The key is the client bucket (IPv4 address or IPv6 /64) resolved through
// the trusted proxies, so a load balancer does not collapse every caller into
// one bucket and a forged X-Forwarded-For from an untrusted hop is ignored.
// Nested under the outer api group so both limits apply.
func (oi *oidcAuthenticator) preAuthRateLimit(period time.Duration, limit int64) gin.HandlerFunc {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}
	key := func(c *gin.Context) string { return clientBucket(oi.clients.clientAddr(c)) }
	return mgin.NewMiddleware(limiter.New(store, rate),
		mgin.WithKeyGetter(key),
		mgin.WithLimitReachedHandler(func(c *gin.Context) {
			oi.lggr.Warnf("OIDC pre-auth rate limit reached path=%s client=%s", c.FullPath(), key(c))
			c.String(http.StatusTooManyRequests, "Limit exceeded")
		}),
	)
}

// deviceNetworkGuard refuses device-flow requests from outside
// DeviceFlowAllowedCIDRs. An empty list turns the device login off. Only CLI
// operators need these routes, and keeping them off the public internet is
// what stops a stranger from starting a flow to phish an approval.
func (oi *oidcAuthenticator) deviceNetworkGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(oi.deviceAllowed) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errDeviceLoginDisabled.Error()})
			return
		}
		addr := oi.clients.clientAddr(c)
		if !prefixesContain(oi.deviceAllowed, addr) {
			// Debug, not Warn: anyone on the internet can trigger this.
			oi.lggr.Debugf("refusing device login request from outside DeviceFlowAllowedCIDRs path=%s client=%s", c.FullPath(), addrString(addr))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errDeviceLoginNetwork.Error()})
			return
		}
		c.Next()
	}
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
	addr := oi.clients.clientAddr(c)
	clientIP := addrString(addr)
	bucket := clientBucket(addr)
	provisional := time.Now().Add(maxDeviceFlowLifetime)
	state := &deviceFlowState{
		expiresAt: provisional,
		clientIP:  clientIP,
		bucket:    bucket,
	}
	if addErr := reserveDeviceFlow(c.Request.Context(), oi.ds, deviceFlowReservation{
		handle:    handle,
		clientIP:  clientIP,
		ipBucket:  bucket,
		expiresAt: provisional,
		pollGrace: minDevicePollGrace,
	}); addErr != nil {
		oi.lggr.Warnf("refusing new device flow client=%s: %v", bucket, addErr)
		if errors.Is(addErr, errTooManyDeviceFlows) || errors.Is(addErr, errTooManyDeviceFlowsPerIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": addErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start device authorization"})
		return
	}
	if addErr := oi.deviceFlows.add(handle, state); addErr != nil {
		_ = removeDeviceFlow(c.Request.Context(), oi.ds, handle)
		oi.lggr.Warnf("refusing new device flow client=%s: %v", bucket, addErr)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": addErr.Error()})
		return
	}

	ctx := context.Background()
	da, err := oi.oauth2Config.DeviceAuth(ctx)
	if err != nil {
		_ = removeDeviceFlow(c.Request.Context(), oi.ds, handle)
		oi.deviceFlows.remove(handle)
		oi.lggr.Errorf("device authorization request failed: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to start device authorization"})
		return
	}

	expiry := capDeviceExpiry(da.Expiry, time.Now())
	if !expiry.After(time.Now()) {
		_ = removeDeviceFlow(c.Request.Context(), oi.ds, handle)
		oi.deviceFlows.remove(handle)
		oi.lggr.Errorf("device authorization expiry is already past")
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to start device authorization"})
		return
	}
	if !verificationURIAllowed(oi.config.ProviderURL(), da.VerificationURI) {
		_ = removeDeviceFlow(c.Request.Context(), oi.ds, handle)
		oi.deviceFlows.remove(handle)
		oi.lggr.Errorf("device authorization verification URI host does not match the configured provider")
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to start device authorization"})
		return
	}
	interval := da.Interval
	if interval == 0 {
		interval = 5 // RFC 8628 default
	}
	state.mu.Lock()
	state.expiresAt = expiry
	state.mu.Unlock()
	if err := setDeviceFlowTiming(c.Request.Context(), oi.ds, handle, expiry, devicePollGrace(interval)); err != nil {
		oi.lggr.Errorf("failed to record device flow expiry: %v", err)
	}

	oi.audit(audit.AuthDeviceFlowStarted, map[string]any{
		"method":      "device_flow",
		"handle_fp":   handleFingerprint(handle),
		"client_ip":   clientIP,
		"user_agent":  truncateUserAgent(c.Request.UserAgent()),
		"expires_in":  int64(time.Until(expiry).Seconds()),
		"remote_addr": c.Request.RemoteAddr,
	})

	// Poll the provider in the background. DeviceAccessToken blocks until the
	// user approves, denies, or the code expires, internally respecting the
	// provider's interval and slow_down. Bounded by the device-code expiry,
	// and cut short if the CLI stops polling.
	go oi.pollDeviceToken(handle, state, da)

	c.JSON(http.StatusOK, DeviceStartResponse{
		DeviceHandle:    handle,
		UserCode:        da.UserCode,
		VerificationURI: da.VerificationURI,
		ExpiresIn:       int64(time.Until(expiry).Seconds()),
		Interval:        interval,
	})
}

// audit forwards to the audit logger when one is configured.
func (oi *oidcAuthenticator) audit(event audit.EventID, data map[string]any) {
	if oi.auditLogger != nil {
		oi.auditLogger.Audit(event, data)
	}
}

// pollDeviceToken runs the blocking provider poll for one flow and records the
// terminal result. On success it verifies the id_token and creates a session
// through the same issueSessionFromIDToken path the browser flow uses.
func (oi *oidcAuthenticator) pollDeviceToken(handle string, state *deviceFlowState, da *oauth2.DeviceAuthResponse) {
	state.mu.Lock()
	deadline := state.expiresAt
	state.mu.Unlock()
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// Free this replica's slot however the flow ends.
	defer oi.deviceFlows.remove(handle)

	abandoned := make(chan struct{})
	go oi.watchDeviceFlow(ctx, handle, cancel, abandoned)

	token, err := oi.oauth2Config.DeviceAccessToken(ctx, da)
	if err != nil {
		select {
		case <-abandoned:
			err = errDeviceFlowAbandoned
		default:
		}
		oi.finishDeviceFlow(state, handle, "", "", "", err)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		oi.finishDeviceFlow(state, handle, "", "", "", errors.New("missing id_token in device token response"))
		return
	}

	// If the handle was swept/expired while we waited on the IdP, do not mint a
	// session nobody can claim (orphan oidc_sessions row).
	if !oi.deviceHandleLive(ctx, handle) {
		oi.lggr.Warnf("device flow handle gone before token completion; abandoning without session")
		oi.finishDeviceFlow(state, handle, "", "", "", errDeviceFlowAbandoned)
		return
	}

	// Device flow does not use a browser nonce; pass empty expectedNonce.
	sessionID, email, role, err := oi.issueSessionFromIDToken(ctx, rawIDToken, "")
	if err != nil {
		oi.finishDeviceFlow(state, handle, "", "", "", err)
		return
	}

	// Race: handle may have been swept between the live check and issue. Drop the
	// orphan session rather than leave a live unclaimable row.
	if !oi.deviceHandleLive(ctx, handle) {
		oi.lggr.Warnf("device flow handle gone after session insert; deleting orphan session")
		oi.deleteOIDCSession(ctx, sessionID)
		oi.finishDeviceFlow(state, handle, "", "", "", errDeviceFlowAbandoned)
		return
	}

	oi.finishDeviceFlow(state, handle, sessionID, email, role, nil)
}

// watchDeviceFlow stops the provider poll once the flow is no longer live,
// usually because the CLI stopped polling. Without it an abandoned flow keeps
// calling the provider until the device code expires. A database error is
// not taken as abandonment, so a blip does not fail a real login.
func (oi *oidcAuthenticator) watchDeviceFlow(ctx context.Context, handle string, cancel context.CancelFunc, abandoned chan<- struct{}) {
	every := oi.liveCheckEvery
	if every <= 0 {
		every = deviceLiveCheckEvery
	}
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			live, err := oi.deviceHandleLiveErr(ctx, handle)
			if err != nil {
				if ctx.Err() == nil {
					oi.lggr.Warnf("device flow live check failed: %v", err)
				}
				continue
			}
			if !live {
				close(abandoned)
				cancel()
				return
			}
		}
	}
}

func (oi *oidcAuthenticator) deviceHandleLiveErr(ctx context.Context, handle string) (bool, error) {
	if oi.ds != nil {
		return deviceFlowLive(ctx, oi.ds, handle)
	}
	return oi.deviceFlows.contains(handle), nil
}

func (oi *oidcAuthenticator) deviceHandleLive(ctx context.Context, handle string) bool {
	if oi.ds != nil {
		live, err := deviceFlowLive(ctx, oi.ds, handle)
		if err != nil {
			oi.lggr.Errorf("device flow live check failed: %v", err)
			return false
		}
		return live
	}
	return oi.deviceFlows.contains(handle)
}

func (oi *oidcAuthenticator) finishDeviceFlow(state *deviceFlowState, handle, sessionID, email string, role clsessions.UserRole, err error) {
	errText := ""
	if err != nil {
		errText = deviceFailureReason(err)
	}
	state.mu.Lock()
	state.done = true
	state.sessionID = sessionID
	state.email = email
	state.role = role
	state.err = err
	startIP := state.clientIP
	state.mu.Unlock()
	if err != nil {
		oi.lggr.Errorf("device authorization flow failed: %v", err)
		oi.audit(audit.AuthLoginFailedDevice, map[string]any{
			"method":          "device_flow",
			"reason":          errText,
			"handle_fp":       handleFingerprint(handle),
			"start_client_ip": startIP,
		})
	}
	if oi.ds != nil && handle != "" {
		if werr := recordDeviceFlowResult(context.Background(), oi.ds, handle, sessionID, email, string(role), errText); werr != nil {
			oi.lggr.Errorf("failed to record device flow result: %v", werr)
		}
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

	// Read the terminal result and consume the handle in one statement so two
	// polls, including polls on different replicas, cannot both be handed the
	// session cookie.
	res, err := consumeDeviceFlow(c.Request.Context(), oi.ds, req.DeviceHandle)
	if err != nil {
		oi.lggr.Errorf("device flow poll failed: %v", err)
		c.JSON(http.StatusInternalServerError, DevicePollResponse{Status: "denied", Message: "Failed to read device authorization"})
		return
	}
	if !res.known {
		c.JSON(http.StatusNotFound, DevicePollResponse{Status: "denied", Message: errUnknownDeviceFlow.Error()})
		return
	}

	if !res.terminal {
		c.JSON(http.StatusOK, DevicePollResponse{Status: "pending"})
		return
	}

	if res.flowErr != nil {
		// Terminal failure already consumed the handle; no session to clean up
		// unless a prior success path left one (should not happen).
		reason := res.flowErr.Error()
		msg, ok := deviceReasonMessages[reason]
		if !ok {
			reason, msg = "failed", deviceReasonMessages["failed"]
		}
		c.JSON(http.StatusOK, DevicePollResponse{Status: "denied", Reason: reason, Message: msg})
		return
	}

	// Bind the session ID to the CLI's cookie jar, exactly as the browser flow
	// does on a successful exchange.
	ginSession := sessions.Default(c)
	ginSession.Set(webauth.SessionIDKey, res.sessionID)
	if err := ginSession.Save(); err != nil {
		// Handle is already consumed (single-winner). Delete the session row so
		// we do not leave a live unclaimable credential; the CLI must restart.
		oi.lggr.Errorf("failed to save device flow session: %v", err)
		oi.deleteOIDCSession(context.Background(), res.sessionID)
		c.JSON(http.StatusInternalServerError, DevicePollResponse{Status: "denied", Message: "Failed to establish session"})
		return
	}

	// A different address redeeming the flow than the one that started it is
	// the signature of a forwarded or phished approval. It can also be an
	// operator whose egress address changed, so flag it rather than refuse.
	redeemIP := addrString(oi.clients.clientAddr(c))
	mismatch := redeemIP != res.startIP
	if mismatch {
		oi.lggr.Warnf("device flow redeemed from a different address than it was started from email=%s start=%s redeem=%s", res.email, res.startIP, redeemIP)
	}
	oi.audit(audit.AuthLoginSuccessNo2FA, map[string]any{
		"email":              res.email,
		"role":               res.role,
		"method":             "device_flow",
		"handle_fp":          handleFingerprint(req.DeviceHandle),
		"client_ip":          redeemIP,
		"start_client_ip":    res.startIP,
		"client_ip_mismatch": mismatch,
		"user_agent":         truncateUserAgent(c.Request.UserAgent()),
	})
	c.JSON(http.StatusOK, DevicePollResponse{Status: "complete"})
}
