package vault

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

var (
	ErrMissingToken                    = errors.New("missing JWT token")
	ErrInvalidToken                    = errors.New("invalid JWT token")
	ErrMissingOrgID                    = errors.New("missing org_id claim")
	ErrMissingTenantID                 = errors.New("missing urn:chainlink:tenant_id or tenant_id claim for JWT-derived vault workflow owner")
	ErrMissingWorkflowOwner            = errors.New("missing workflow_owner in authorization_details")
	ErrMissingRequestDigest            = errors.New("missing request_digest in authorization_details")
	ErrVaultSecretManagementNotEnabled = errors.New("claim_vault_secret_management_enabled claim must be true")
	ErrJWTTenantIDJobSpecMismatch      = errors.New("JWT tenant id does not match auth0 tenantID from job specification")
	ErrJWKSFetchFailed                 = errors.New("failed to fetch JWKS")
	ErrJWKSKeyNotFound                 = errors.New("signing key not found in JWKS")
)

const ClaimVaultSecretManagementEnabled = "urn:chainlink:claim_vault_secret_management_enabled"

// ClaimChainlinkTenantID is the CRE tenant numeric id (typically a decimal string) as emitted on Auth0 JWTs.
// See cre-platform-graphql/internal/auth/jwt_auth0.go (resolveStringClaim(..., "tenant_id")).
const ClaimChainlinkTenantID = "urn:chainlink:tenant_id"

const (
	defaultJWTAuthJobSpecTenantID uint64 = 1
	defaultJWKSRefreshInterval           = 15 * time.Minute
	defaultHTTPTimeout                   = 5 * time.Second
)

// Auth0Config captures the Vault JWT issuer settings shared by gateway and node handlers.
type Auth0Config struct {
	IssuerURL string `json:"issuerURL" toml:"issuerURL" yaml:"issuerURL"`
	Audience  string `json:"audience" toml:"audience" yaml:"audience"`
	TenantID  uint64 `json:"tenantID" toml:"tenantID" yaml:"tenantID"`
}

// JWTBasedAuthConfig holds the configuration for JWTBasedAuth validation.
type JWTBasedAuthConfig struct {
	IssuerURL           string
	Audience            string
	TenantID            uint64        // omit or zero defaults to defaultJWTAuthJobSpecTenantID; must match JWT urn:chainlink:tenant_id / tenant_id claim
	JWKSRefreshInterval time.Duration // minimum interval between JWKS fetches; 0 uses default (30s)
	HTTPClient          *http.Client  // nil uses a default client with 5s timeout
}

// JWTClaims contains the validated claims extracted from an Auth0 JWT
// relevant to Vault request authorization.
type JWTClaims struct {
	OrgID         string
	TenantID      uint64 // from urn:chainlink:tenant_id or tenant_id
	WorkflowOwner string // from authorization_details
	RequestDigest string // from authorization_details
	ExpiresAt     time.Time
	OAuthScopes   []string // from scope / permissions claims
}

type jsonWebKey struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jsonWebKeySet struct {
	Keys []jsonWebKey `json:"keys"`
}

// JWTBasedAuth verifies Auth0-issued RS256 JWTs using the provider's
// public JWKS endpoint and extracts Vault-specific claims (org_id,
// tenant id, optional workflow_owner, request_digest). It is safe for concurrent use.
//
// JWKS keys are fetched lazily on the first token validation and refreshed
// on key-ID misses, rate-limited to at most once per JWKSRefreshInterval.
//
// Reference: cre-platform-graphql/internal/auth/jwt_auth0.go
type jwtBasedAuth struct {
	services.Service
	eng *services.Engine

	issuerURL       string
	audience        string
	jwksURL         string
	refreshInterval time.Duration
	authEnabledGate limits.GateLimiter

	expectedTenantID uint64

	mu            sync.RWMutex
	keySet        *jsonWebKeySet
	lastRefreshed time.Time

	refreshMu sync.Mutex // serializes JWKS refresh attempts

	httpClient *http.Client
	lggr       logger.Logger

	limitsFactory limits.Factory
}

type jwtBasedAuthOptions struct {
	authEnabledGate limits.GateLimiter
}

// JWTBasedAuthOption customizes JWTBasedAuth construction without multiplying constructors.
type JWTBasedAuthOption func(*jwtBasedAuthOptions)

// WithJWTBasedAuthGateLimiter overrides the gate limiter that decides whether JWT-based auth is enabled.
func WithJWTBasedAuthGateLimiter(gateLimiter limits.GateLimiter) JWTBasedAuthOption {
	return func(opts *jwtBasedAuthOptions) {
		opts.authEnabledGate = gateLimiter
	}
}

// NewJWTBasedAuth creates a JWTBasedAuth authorizer that verifies Auth0-issued JWTs
// against the provider's JWKS endpoint. The JWKS is fetched lazily on first
// use and refreshed on key-ID cache misses (rate-limited).
func NewJWTBasedAuth(cfg JWTBasedAuthConfig, limitsFactory limits.Factory, lggr logger.Logger, opts ...JWTBasedAuthOption) (*jwtBasedAuth, error) {
	options := jwtBasedAuthOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	if options.authEnabledGate == nil {
		options.authEnabledGate = newVaultJWTAuthEnabledGateLimiter(limitsFactory, lggr)
	}
	if cfg.IssuerURL == "" {
		return nil, errors.New("issuer URL is required")
	}
	if cfg.Audience == "" {
		return nil, errors.New("audience is required")
	}

	expectedTenantID := cfg.TenantID
	if expectedTenantID == 0 {
		expectedTenantID = defaultJWTAuthJobSpecTenantID
	}

	trimmedIssuer := strings.TrimSuffix(cfg.IssuerURL, "/")
	jwksURL := trimmedIssuer + "/.well-known/jwks.json"

	refreshInterval := cfg.JWKSRefreshInterval
	if refreshInterval == 0 {
		refreshInterval = defaultJWKSRefreshInterval
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	v := &jwtBasedAuth{
		issuerURL:        cfg.IssuerURL,
		audience:         cfg.Audience,
		jwksURL:          jwksURL,
		refreshInterval:  refreshInterval,
		authEnabledGate:  options.authEnabledGate,
		expectedTenantID: expectedTenantID,
		httpClient:       httpClient,
		lggr:             logger.Named(lggr, "VaultJWTBasedAuth"),
		limitsFactory:    limitsFactory,
	}
	v.Service, v.eng = services.Config{
		Name:  "VaultJWTBasedAuth",
		Start: v.start,
		Close: v.close,
	}.NewServiceEngine(v.lggr)

	return v, nil
}

func newVaultJWTAuthEnabledGateLimiter(limitsFactory limits.Factory, lggr logger.Logger) limits.GateLimiter {
	limiter, err := limits.MakeGateLimiter(limitsFactory, cresettings.Default.VaultJWTAuthEnabled)
	if err != nil {
		logger.Named(lggr, "VaultJWTBasedAuth").Errorw("failed to create VaultJWTAuthEnabled limiter", "error", err)
		return limits.NewGateLimiter(false)
	}

	return limiter
}

func (v *jwtBasedAuth) start(context.Context) error {
	v.eng.GoTick(services.NewTicker(v.refreshInterval), func(ctx context.Context) {
		if err := v.refreshJWKS(ctx); err != nil {
			v.lggr.Warnw("periodic JWKS refresh failed", "error", err)
		}
	})
	return nil
}

func (v *jwtBasedAuth) close() error {
	return v.authEnabledGate.Close()
}

// AuthorizeRequest verifies JWTBasedAuth state and token claims, and returns a common AuthResult.
func (v *jwtBasedAuth) AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	isEnabled, err := v.authEnabledGate.Limit(ctx)
	if err != nil {
		v.lggr.Errorw("failed to resolve JWTBasedAuth gate", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, fmt.Errorf("failed to resolve JWTBasedAuth gate: %w", err)
	}
	if !isEnabled {
		v.lggr.Debugw("JWTBasedAuth rejected request because it is disabled", "method", req.Method, "requestID", req.ID)
		return nil, errors.New("JWTBasedAuth is disabled")
	}

	claims, err := v.validateToken(ctx, req.Auth)
	if err != nil {
		v.lggr.Debugw("JWTBasedAuth token validation failed", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, fmt.Errorf("invalid JWT auth token: %w", err)
	}

	if scopeErr := enforceVaultJWTOAuthScopes(req.Method, claims.OAuthScopes); scopeErr != nil {
		v.lggr.Debugw("JWTBasedAuth OAuth scope rejected request", "method", req.Method, "requestID", req.ID, "orgID", claims.OrgID, "scopes", claims.OAuthScopes, "error", scopeErr)
		return nil, fmt.Errorf("invalid JWT auth token: %w", scopeErr)
	}

	if claims.TenantID == 0 {
		return nil, ErrMissingTenantID
	}
	if claims.TenantID != v.expectedTenantID {
		v.lggr.Debugw("JWT tenant id does not match job spec auth0 tenantID", "method", req.Method, "requestID", req.ID, "orgID", claims.OrgID, "claimsTenantID", claims.TenantID, "expectedTenantID", v.expectedTenantID)
		return nil, fmt.Errorf("%w: jwt tenant id %d expected tenant id %d", ErrJWTTenantIDJobSpecMismatch, claims.TenantID, v.expectedTenantID)
	}

	requestDigest, err := req.Digest()
	if err != nil {
		v.lggr.Debugw("JWTBasedAuth failed to compute request digest", "method", req.Method, "requestID", req.ID, "orgID", claims.OrgID, "workflowOwner", claims.WorkflowOwner, "error", err)
		return nil, fmt.Errorf("failed to compute request digest: %w", err)
	}

	if !strings.EqualFold(requestDigest, claims.RequestDigest) {
		v.lggr.Debugw("JWTBasedAuth request digest mismatch", "method", req.Method, "requestID", req.ID, "orgID", claims.OrgID, "workflowOwner", claims.WorkflowOwner, "computedDigest", requestDigest, "claimedDigest", claims.RequestDigest)
		return nil, fmt.Errorf("request digest mismatch: computed=%s claimed=%s", requestDigest, claims.RequestDigest)
	}

	derivedWorkflowOwner, err := DeriveJWTAuthorizedVaultWorkflowOwner(claims.OrgID, claims.TenantID, claims.WorkflowOwner)
	if err != nil {
		v.lggr.Debugw("JWTBasedAuth failed to derive authorized workflow owner", "method", req.Method, "requestID", req.ID, "orgID", claims.OrgID, "error", err)
		return nil, fmt.Errorf("invalid JWT auth token: %w", err)
	}

	v.lggr.Debugw("JWTBasedAuth authorization succeeded", "method", req.Method, "requestID", req.ID, "orgID", claims.OrgID, "workflowOwner", derivedWorkflowOwner, "digest", requestDigest, "expiresAt", claims.ExpiresAt.UTC().Unix())
	return &AuthResult{
		orgID:         claims.OrgID,
		workflowOwner: derivedWorkflowOwner,
		digest:        requestDigest,
		expiresAt:     claims.ExpiresAt.UTC().Unix(),
	}, nil
}

// validateToken verifies the JWT signature via Auth0 JWKS, validates
// standard claims (iss, aud, exp), and extracts Vault-specific claims
// (org_id, tenant id, optional workflow_owner, request_digest).
func (v *jwtBasedAuth) validateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, ErrMissingToken
	}

	unverified, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	kid, ok := unverified.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fmt.Errorf("%w: missing kid header", ErrInvalidToken)
	}

	rsaKey, err := v.resolveSigningKey(ctx, kid)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, methodOK := token.Method.(*jwt.SigningMethodRSA); !methodOK {
			return nil, fmt.Errorf("%w: unsupported alg %v", ErrInvalidToken, token.Header["alg"])
		}
		return rsaKey, nil
	},
		jwt.WithIssuer(v.issuerURL),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w. Expected Issuer: %s, Actual Issuer: %s", ErrInvalidToken, err, v.issuerURL, unverified.Claims.(jwt.MapClaims)["iss"])
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return extractVaultClaims(claims)
}

func extractVaultClaims(claims jwt.MapClaims) (*JWTClaims, error) {
	orgID, _ := claims["org_id"].(string)
	if orgID == "" {
		return nil, ErrMissingOrgID
	}

	if v, ok := claims[ClaimVaultSecretManagementEnabled].(string); !ok || v != "true" {
		return nil, ErrVaultSecretManagementNotEnabled
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("%w: invalid exp claim", ErrInvalidToken)
	}

	tenantID, err := extractTenantNumericIDFromClaims(claims)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid tenant id claim", ErrInvalidToken)
	}

	workflowOwner, requestDigest, err := extractAuthorizationDetails(claims)
	if err != nil {
		return nil, err
	}

	oauthScopes := extractOAuthScopesFromClaims(claims)

	return &JWTClaims{
		OrgID:         orgID,
		TenantID:      tenantID,
		WorkflowOwner: workflowOwner,
		RequestDigest: requestDigest,
		ExpiresAt:     exp.Time,
		OAuthScopes:   oauthScopes,
	}, nil
}

func extractTenantNumericIDFromClaims(claims jwt.MapClaims) (uint64, error) {
	var raw interface{}
	ok := false
	if v, exists := claims[ClaimChainlinkTenantID]; exists && v != nil {
		raw, ok = v, true
	} else if v, exists := claims["tenant_id"]; exists && v != nil {
		raw, ok = v, true
	}
	if !ok {
		return 0, nil
	}
	id, err := parseJWTUnsignedIntegerClaim(raw)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func parseJWTUnsignedIntegerClaim(raw interface{}) (uint64, error) {
	switch v := raw.(type) {
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return 0, nil
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return n, nil
	case json.Number:
		n, err := strconv.ParseUint(v.String(), 10, 64)
		if err != nil {
			return 0, err
		}
		return n, nil
	case float64:
		if v <= 0 || v != math.Trunc(v) || v > math.MaxUint64 {
			return 0, strconv.ErrSyntax
		}
		return uint64(v), nil
	default:
		return 0, strconv.ErrSyntax
	}
}

func extractAuthorizationDetails(claims jwt.MapClaims) (workflowOwner, requestDigest string, err error) {
	rawDetails, ok := claims["authorization_details"]
	if !ok {
		return "", "", ErrMissingRequestDigest
	}

	details, ok := rawDetails.([]interface{})
	if !ok {
		return "", "", fmt.Errorf("%w: authorization_details must be an array", ErrInvalidToken)
	}

	for _, rawDetail := range details {
		detail, ok := rawDetail.(map[string]interface{})
		if !ok {
			continue
		}
		authDetailType, _ := detail["type"].(string)
		authDetailValue, _ := detail["value"].(string)
		switch authDetailType {
		case "request_digest":
			requestDigest = authDetailValue
		case "workflow_owner":
			workflowOwner = authDetailValue
		}
	}

	if requestDigest == "" {
		return "", "", ErrMissingRequestDigest
	}

	return workflowOwner, requestDigest, nil
}

// resolveSigningKey looks up the RSA public key for the given kid from the
// JWKS cache, refreshing the cache if necessary.
func (v *jwtBasedAuth) resolveSigningKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	key, err := v.findCachedKey(kid)
	if err != nil {
		return nil, err
	}
	if key != nil {
		return key, nil
	}

	if refreshErr := v.refreshJWKS(ctx); refreshErr != nil {
		v.lggr.Warnw("JWKS refresh failed", "error", refreshErr, "kid", kid)
		return nil, fmt.Errorf("%w: kid=%s", ErrJWKSKeyNotFound, kid)
	}

	key, err = v.findCachedKey(kid)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, fmt.Errorf("%w: kid=%s", ErrJWKSKeyNotFound, kid)
	}

	return key, nil
}

func (v *jwtBasedAuth) findCachedKey(kid string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.keySet == nil {
		return nil, nil
	}

	for _, key := range v.keySet.Keys {
		if key.Kid == kid {
			return parseRSAPublicKey(key)
		}
	}

	return nil, nil
}

// refreshJWKS fetches the JWKS from Auth0. Concurrent callers are serialized
// via refreshMu; if a recent fetch already happened within refreshInterval
// the call is a no-op.
func (v *jwtBasedAuth) refreshJWKS(ctx context.Context) error {
	v.refreshMu.Lock()
	defer v.refreshMu.Unlock()

	v.mu.RLock()
	if v.keySet != nil && time.Since(v.lastRefreshed) < v.refreshInterval {
		v.mu.RUnlock()
		return nil
	}
	v.mu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrJWKSFetchFailed, err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrJWKSFetchFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: HTTP %d", ErrJWKSFetchFailed, resp.StatusCode)
	}

	const maxJWKSBodySize = 1 << 20 // 1 MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxJWKSBodySize))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrJWKSFetchFailed, err)
	}

	var keySet jsonWebKeySet
	if err := json.Unmarshal(body, &keySet); err != nil {
		return fmt.Errorf("%w: invalid JWKS: %w", ErrJWKSFetchFailed, err)
	}

	v.mu.Lock()
	v.keySet = &keySet
	v.lastRefreshed = time.Now()
	v.mu.Unlock()

	v.lggr.Infow("Refreshed JWKS", "numKeys", len(keySet.Keys), "url", v.jwksURL)
	return nil
}

func parseRSAPublicKey(key jsonWebKey) (*rsa.PublicKey, error) {
	if key.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported key type: %s", key.Kty)
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RSA modulus: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RSA exponent: %w", err)
	}

	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}
