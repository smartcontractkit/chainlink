package vault

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// --- test helpers ---

type testRSAKey struct {
	kid        string
	privateKey *rsa.PrivateKey
}

type testJWKSServer struct {
	server *httptest.Server
	mu     sync.Mutex
	keys   []testRSAKey
	hits   chan struct{}
}

func newTestJWKSServer(t *testing.T, keys ...testRSAKey) *testJWKSServer {
	t.Helper()
	s := &testJWKSServer{keys: keys, hits: make(chan struct{}, 32)}
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		currentKeys := s.keys
		s.mu.Unlock()
		select {
		case s.hits <- struct{}{}:
		default:
		}

		ks := jsonWebKeySet{}
		for _, k := range currentKeys {
			ks.Keys = append(ks.Keys, testRSAKeyToJWK(k.kid, &k.privateKey.PublicKey))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ks)
	})
	s.server = httptest.NewServer(mux)
	t.Cleanup(s.server.Close)
	return s
}

func (s *testJWKSServer) URL() string { return s.server.URL }

func (s *testJWKSServer) waitForHits(t *testing.T, count int, timeout time.Duration) {
	t.Helper()

	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	for range count {
		select {
		case <-s.hits:
		case <-deadline.C:
			t.Fatalf("timed out waiting for %d JWKS hits", count)
		}
	}
}

func (s *testJWKSServer) setKeys(keys ...testRSAKey) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys = keys
}

func testRSAKeyToJWK(kid string, pub *rsa.PublicKey) jsonWebKey {
	return jsonWebKey{
		Kid: kid,
		Alg: "RS256",
		Kty: "RSA",
		Use: "sig",
		N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}
}

func generateTestRSAKey(t *testing.T, kid string) testRSAKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return testRSAKey{kid: kid, privateKey: key}
}

func testJWTExpectedWorkflowOwner(tb testing.TB, tenantID uint64, orgID string) string {
	tb.Helper()
	prefix := strconv.FormatUint(tenantID, 10)
	addr, err := workflows.GenerateWorkflowOwnerAddress(prefix, orgID)
	require.NoError(tb, err)
	return common.BytesToAddress(addr).Hex()
}

func createTestJWT(t *testing.T, key testRSAKey, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = key.kid
	tokenString, err := token.SignedString(key.privateKey)
	require.NoError(t, err)
	return tokenString
}

func validTestClaims(tb testing.TB, issuer, audience string) jwt.MapClaims {
	return jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org_test123",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": "abc123def456",
			},
			map[string]interface{}{
				"type":  "workflow_owner",
				"value": testJWTExpectedWorkflowOwner(tb, 1, "org_test123"),
			},
		},
	}
}

func newTestValidator(t *testing.T, issuer, audience string) *jwtBasedAuth {
	t.Helper()
	v, err := NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL:           issuer,
		Audience:            audience,
		TenantID:            1,
		JWKSRefreshInterval: time.Millisecond,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, logger.TestLogger(t), WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.NoError(t, err)
	return v
}

// --- tests ---

func TestJWTBasedAuth_ValidToken(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	tokenString := createTestJWT(t, rsaKey, validTestClaims(t, issuer, audience))

	result, err := v.validateToken(context.Background(), tokenString)
	require.NoError(t, err)
	assert.Equal(t, "org_test123", result.OrgID)
	expectedOwner := testJWTExpectedWorkflowOwner(t, 1, "org_test123")
	assert.Equal(t, expectedOwner, result.WorkflowOwner)
	assert.Equal(t, uint64(1), result.TenantID)
	assert.Equal(t, "abc123def456", result.RequestDigest)
	assert.False(t, result.ExpiresAt.IsZero())
}

func TestJWTBasedAuth_ValidTokenWithoutWorkflowOwner(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org_no_wfowner",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": "digest456",
			},
		},
	}
	tokenString := createTestJWT(t, rsaKey, claims)

	result, err := v.validateToken(context.Background(), tokenString)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "org_no_wfowner", result.OrgID)
	require.Empty(t, result.WorkflowOwner)
	require.Equal(t, uint64(1), result.TenantID)
	require.Equal(t, "digest456", result.RequestDigest)
}

func TestJWTBasedAuth_ExpiredToken(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidToken)
	assert.Contains(t, err.Error(), "expired")
}

func TestJWTBasedAuth_WrongIssuer(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, "https://wrong-issuer.auth0.com/", audience)
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTBasedAuth_WrongAudience(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, "https://wrong-audience.com")
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTBasedAuth_MissingOrgID(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	delete(claims, "org_id")
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingOrgID)
}

func TestJWTBasedAuth_MissingVaultSecretManagementClaim(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	delete(claims, ClaimVaultSecretManagementEnabled)
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrVaultSecretManagementNotEnabled)
}

func TestJWTBasedAuth_VaultSecretManagementClaimNotTrue(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	claims[ClaimVaultSecretManagementEnabled] = "false"
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrVaultSecretManagementNotEnabled)
}

func TestJWTBasedAuth_MissingRequestDigest(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	claims["authorization_details"] = []interface{}{
		map[string]interface{}{
			"type":  "workflow_owner",
			"value": "0xAbCd",
		},
	}
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingRequestDigest)
}

func TestJWTBasedAuth_MissingAuthorizationDetails(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	delete(claims, "authorization_details")
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingRequestDigest)
}

func TestJWTBasedAuth_InvalidSignature(t *testing.T) {
	goodKey := generateTestRSAKey(t, "key-1")
	badKey := generateTestRSAKey(t, "key-1") // same kid, different key material
	jwksServer := newTestJWKSServer(t, goodKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	tokenString := createTestJWT(t, badKey, claims) // signed with wrong private key

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTBasedAuth_EmptyToken(t *testing.T) {
	v, err := NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL: "https://example.auth0.com/",
		Audience:  "https://api.test.chain.link",
		TenantID:  1,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, logger.TestLogger(t), WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.NoError(t, err)

	_, err = v.validateToken(context.Background(), "")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingToken)
}

func TestJWTBasedAuth_JWKSKeyRotation(t *testing.T) {
	keyA := generateTestRSAKey(t, "key-A")
	keyB := generateTestRSAKey(t, "key-B")

	jwksServer := newTestJWKSServer(t, keyA)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	// Token signed with key-A succeeds
	claimsA := validTestClaims(t, issuer, audience)
	tokenA := createTestJWT(t, keyA, claimsA)
	resultA, err := v.validateToken(context.Background(), tokenA)
	require.NoError(t, err)
	assert.Equal(t, "org_test123", resultA.OrgID)

	// Simulate key rotation: JWKS now serves only key-B
	jwksServer.setKeys(keyB)

	// Allow the refresh interval to elapse so the next miss triggers a fetch
	time.Sleep(2 * time.Millisecond)

	// Token signed with key-B succeeds after JWKS refresh
	claimsB := validTestClaims(t, issuer, audience)
	claimsB["org_id"] = "org_after_rotation"
	claimsB["authorization_details"] = []interface{}{
		map[string]interface{}{
			"type":  "request_digest",
			"value": "abc123def456",
		},
		map[string]interface{}{
			"type":  "workflow_owner",
			"value": testJWTExpectedWorkflowOwner(t, 1, "org_after_rotation"),
		},
	}
	tokenB := createTestJWT(t, keyB, claimsB)
	resultB, err := v.validateToken(context.Background(), tokenB)
	require.NoError(t, err)
	assert.Equal(t, "org_after_rotation", resultB.OrgID)
}

func TestJWTBasedAuth_AuthorizationDetailsFromTypedArray(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org_single",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"authorization_details": []interface{}{
			map[string]interface{}{"type": "request_digest", "value": "single_digest"},
			map[string]interface{}{"type": "workflow_owner", "value": testJWTExpectedWorkflowOwner(t, 1, "org_single")},
		},
	}
	tokenString := createTestJWT(t, rsaKey, claims)

	result, err := v.validateToken(context.Background(), tokenString)
	require.NoError(t, err)
	assert.Equal(t, "org_single", result.OrgID)
	assert.Equal(t, "single_digest", result.RequestDigest)
	assert.Equal(t, testJWTExpectedWorkflowOwner(t, 1, "org_single"), result.WorkflowOwner)
}

func TestJWTBasedAuth_UnsupportedAlgorithm(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	// Create a token signed with HMAC instead of RSA
	claims := validTestClaims(t, issuer, audience)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = rsaKey.kid
	tokenString, err := token.SignedString([]byte("hmac-secret"))
	require.NoError(t, err)

	_, err = v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTBasedAuth_JWKSServerUnavailable(t *testing.T) {
	// Start a server that always returns 500
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	issuer := server.URL + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	rsaKey := generateTestRSAKey(t, "key-1")
	claims := validTestClaims(t, issuer, audience)
	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrJWKSKeyNotFound)
}

func TestJWTBasedAuth_StartRefreshesJWKSPeriodically(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	v, err := NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL:           jwksServer.URL() + "/",
		Audience:            "https://api.test.chain.link",
		TenantID:            1,
		JWKSRefreshInterval: 10 * time.Millisecond,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, logger.TestLogger(t), WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.NoError(t, err)

	require.NoError(t, v.Start(t.Context()))
	jwksServer.waitForHits(t, 2, time.Second)
	require.NoError(t, v.Close())
}

func TestNewJWTBasedAuth_InvalidConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	var v *jwtBasedAuth

	_, err := NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL: "",
		Audience:  "https://api.test.chain.link",
		TenantID:  1,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, lggr, WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "issuer URL is required")

	_, err = NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL: "https://example.auth0.com/",
		Audience:  "",
		TenantID:  1,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, lggr, WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "audience is required")

	v, err = NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL:           "https://example.auth0.com/",
		Audience:            "https://api.test.chain.link",
		TenantID:            0,
		JWKSRefreshInterval: time.Millisecond,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, lggr, WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.NoError(t, err)
	require.NotNil(t, v)
	defer func() { _ = v.Close() }()
	require.Equal(t, uint64(1), v.expectedTenantID)
}

func TestNewJWTBasedAuth_UsesVaultJWTAuthEnabledLimiter_Disabled(t *testing.T) {
	setDefaultGetter(t, `{}`)

	v, err := NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL: "https://example.auth0.com/",
		Audience:  "https://api.test.chain.link",
		TenantID:  1,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, logger.TestLogger(t))
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-1",
		Method: vaulttypes.MethodSecretsList,
		Auth:   "token",
	}

	_, err = v.AuthorizeRequest(t.Context(), req)
	require.Error(t, err)
	require.ErrorContains(t, err, "JWTBasedAuth is disabled")
}

func TestNewJWTBasedAuth_UsesVaultJWTAuthEnabledLimiter_Enabled(t *testing.T) {
	setDefaultGetter(t, `{"global":{"VaultJWTAuthEnabled":true}}`)

	v, err := NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL: "https://example.auth0.com/",
		Audience:  "https://api.test.chain.link",
		TenantID:  1,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, logger.TestLogger(t))
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-1",
		Method: vaulttypes.MethodSecretsList,
	}

	_, err = v.AuthorizeRequest(t.Context(), req)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid JWT auth token")
	require.ErrorContains(t, err, ErrMissingToken.Error())
}

func TestJWTBasedAuth_AuthorizeCreateRequestFromRawJSON(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	derivedOrg123Owner := testJWTExpectedWorkflowOwner(t, 1, "org-123")
	rawRequest := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.create","params":{"request_id":"req-1","encrypted_secrets":[{"id":{"key":"7611","namespace":"main","owner":"%s"},"encrypted_value":"cipher+/=="}]}}`, derivedOrg123Owner))
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"scope":                           OAuthScopeVaultSecretsCreate,
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
			map[string]interface{}{
				"type":  "workflow_owner",
				"value": derivedOrg123Owner,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	authResult, err := v.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "org-123", authResult.OrgID())
	require.Equal(t, derivedOrg123Owner, authResult.WorkflowOwner())
	require.Equal(t, digest, authResult.Digest())
}

func TestJWTBasedAuth_RejectsMissingOAuthScope(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	rawRequest := []byte(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.create","params":{"request_id":"req-1","encrypted_secrets":[{"id":{"key":"7611","namespace":"main","owner":"org-123"},"encrypted_value":"cipher+/=="}]}}`)
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	_, err = v.AuthorizeRequest(t.Context(), req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrMissingVaultOAuthScope)
}

func TestJWTBasedAuth_RejectsMismatchedOAuthScope(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	rawRequest := []byte(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.delete","params":{"request_id":"req-1","ids":[{"key":"7611","namespace":"main","owner":"org-123"}]}}`)
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		"scope":                           OAuthScopeVaultSecretsCreate,
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	_, err = v.AuthorizeRequest(t.Context(), req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrVaultJWTScopeDenied)
}

func TestJWTBasedAuth_RejectsMultipleVaultOAuthScopes(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	rawRequest := []byte(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.create","params":{"request_id":"req-1","encrypted_secrets":[{"id":{"key":"7611","namespace":"main","owner":"org-123"},"encrypted_value":"cipher+/=="}]}}`)
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		"scope":                           OAuthScopeVaultSecretsCreate + " " + OAuthScopeVaultSecretsDelete,
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	_, err = v.AuthorizeRequest(t.Context(), req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrVaultJWTMultipleOAuthScopes)
}

func TestJWTBasedAuth_AcceptsOpenIDPlusSingleVaultScope(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	derivedOrg123Owner := testJWTExpectedWorkflowOwner(t, 1, "org-123")
	rawRequest := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.create","params":{"request_id":"req-1","encrypted_secrets":[{"id":{"key":"7611","namespace":"main","owner":"%s"},"encrypted_value":"cipher+/=="}]}}`, derivedOrg123Owner))
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"scope":                           "openid profile " + OAuthScopeVaultSecretsCreate,
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	authResult, err := v.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "org-123", authResult.OrgID())
}

func TestJWTBasedAuth_AcceptsPermissionsClaimInsteadOfScope(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	derivedOrg123Owner := testJWTExpectedWorkflowOwner(t, 1, "org-123")
	rawRequest := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.list","params":{"request_id":"req-1","owner":"%s","namespace":"main"}}`, derivedOrg123Owner))
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"permissions":                     []interface{}{OAuthScopeVaultSecretsList},
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	authResult, err := v.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "org-123", authResult.OrgID())
}

func TestJWTBasedAuth_AuthorizeCreateRequestWithoutWorkflowOwnerWhenIdentifiersUseOrgID(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	derivedOrg123Owner := testJWTExpectedWorkflowOwner(t, 1, "org-123")
	rawRequest := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.create","params":{"request_id":"req-1","encrypted_secrets":[{"id":{"key":"7611","namespace":"main","owner":"%s"},"encrypted_value":"cipher+/=="}]}}`, derivedOrg123Owner))
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"scope":                           OAuthScopeVaultSecretsCreate,
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	authResult, err := v.AuthorizeRequest(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "org-123", authResult.OrgID())
	require.Equal(t, derivedOrg123Owner, authResult.WorkflowOwner())
	require.Equal(t, digest, authResult.Digest())
}

func TestJWTBasedAuth_RejectsCreateRequestWithoutWorkflowOwnerWhenIdentifierOwnerDiffers(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	rawRequest := []byte(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.create","params":{"request_id":"req-1","encrypted_secrets":[{"id":{"key":"7611","namespace":"main","owner":"0xAbCd"},"encrypted_value":"cipher+/=="}]}}`)
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"scope":                           OAuthScopeVaultSecretsCreate,
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	// Owner binding is enforced by the composite Authorizer, not the JWT layer directly.
	a := NewAuthorizer(nil, v, logger.TestLogger(t))
	authResult, err := a.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.Error(t, err)
	require.ErrorContains(t, err, "encrypted secret owner at index 0 \"0xAbCd\" does not match authorized workflow owner")
}

func TestJWTBasedAuth_AuthorizeRequest_RejectsMissingTenantClaim(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	derivedOrg123Owner := testJWTExpectedWorkflowOwner(t, 1, "org-123")
	rawRequest := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.list","params":{"request_id":"req-1","owner":"%s","namespace":"main"}}`, derivedOrg123Owner))
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		"permissions":                     []interface{}{OAuthScopeVaultSecretsList},
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	_, err = v.AuthorizeRequest(t.Context(), req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrMissingTenantID)
}

func TestJWTBasedAuth_AuthorizeRequest_RejectsTenantAgainstJobSpecMismatch(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://vault.test.chain.link"

	v, err := NewJWTBasedAuth(JWTBasedAuthConfig{
		IssuerURL:           issuer,
		Audience:            audience,
		TenantID:            2,
		JWKSRefreshInterval: time.Millisecond,
	}, limits.Factory{Settings: cresettings.DefaultGetter}, logger.TestLogger(t), WithJWTBasedAuthGateLimiter(limits.NewGateLimiter(true)))
	require.NoError(t, err)

	owner := testJWTExpectedWorkflowOwner(t, 1, "org-123")
	rawRequest := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"req-1","method":"vault.secrets.delete","params":{"request_id":"req-1","ids":[{"owner":"%s","namespace":"main","key":"k"}]}}`, owner))
	req, err := jsonrpc.DecodeRequest[json.RawMessage](rawRequest, "")
	require.NoError(t, err)

	digest, err := req.Digest()
	require.NoError(t, err)

	token := createTestJWT(t, rsaKey, jwt.MapClaims{
		"iss":                             issuer,
		"aud":                             audience,
		"exp":                             jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"iat":                             jwt.NewNumericDate(time.Now()),
		"org_id":                          "org-123",
		ClaimVaultSecretManagementEnabled: "true",
		ClaimChainlinkTenantID:            "1",
		"permissions":                     []interface{}{OAuthScopeVaultSecretsDelete},
		"authorization_details": []interface{}{
			map[string]interface{}{
				"type":  "request_digest",
				"value": digest,
			},
			map[string]interface{}{
				"type":  "workflow_owner",
				"value": owner,
			},
		},
	})

	req, err = jsonrpc.DecodeRequest[json.RawMessage](rawRequest, token)
	require.NoError(t, err)

	authResult, err := v.AuthorizeRequest(t.Context(), req)
	require.Nil(t, authResult)
	require.ErrorIs(t, err, ErrJWTTenantIDJobSpecMismatch)
}

func TestJWTBasedAuth_InvalidTenantNumericClaim(t *testing.T) {
	rsaKey := generateTestRSAKey(t, "key-1")
	jwksServer := newTestJWKSServer(t, rsaKey)

	issuer := jwksServer.URL() + "/"
	audience := "https://api.test.chain.link"
	v := newTestValidator(t, issuer, audience)

	claims := validTestClaims(t, issuer, audience)
	claims[ClaimChainlinkTenantID] = "NaN"

	tokenString := createTestJWT(t, rsaKey, claims)

	_, err := v.validateToken(context.Background(), tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func setDefaultGetter(t *testing.T, payload string) {
	t.Helper()

	prev := cresettings.DefaultGetter
	t.Cleanup(func() {
		cresettings.DefaultGetter = prev
	})

	getter, err := settings.NewJSONGetter([]byte(payload))
	require.NoError(t, err)
	cresettings.DefaultGetter = getter
}
