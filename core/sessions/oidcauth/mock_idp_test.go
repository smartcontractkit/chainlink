package oidcauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
)

// mockIDP is a minimal OpenID Connect provider for tests. It serves a discovery
// document, a JWKS, a device authorization endpoint, and a token endpoint that
// returns RSA-signed id_tokens. It lets tests exercise the real go-oidc
// verifier (signature, issuer, audience, expiry) rather than stubbing it out.
type mockIDP struct {
	server *httptest.Server
	signer jose.Signer
	key    *rsa.PrivateKey
	kid    string

	// Knobs the tests flip to shape the id_token the token endpoint mints.
	audience       string        // "aud" claim; default = clientID
	issuerOverride string        // "iss" claim; default = server URL
	groups         []string      // groups claim
	email          string        // email claim
	lifetime       time.Duration // exp relative to now; negative => already expired

	clientID string
}

func newMockIDP(t *testing.T, clientID string) *mockIDP {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	kid := "test-key-1"
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", kid),
	)
	require.NoError(t, err)

	m := &mockIDP{
		signer:   signer,
		key:      key,
		kid:      kid,
		groups:   []string{AdminClaim},
		email:    "user@example.com",
		lifetime: time.Hour,
		clientID: clientID,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		base := m.server.URL
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":                                base,
			"authorization_endpoint":                base + "/authorize",
			"token_endpoint":                        base + "/token",
			"jwks_uri":                              base + "/keys",
			"device_authorization_endpoint":         base + "/device/authorize",
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
	})
	mux.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
			Key:       key.Public(),
			KeyID:     kid,
			Algorithm: "RS256",
			Use:       "sig",
		}}}
		_ = json.NewEncoder(w).Encode(jwks)
	})
	mux.HandleFunc("/device/authorize", func(w http.ResponseWriter, r *http.Request) {
		base := m.server.URL
		_ = json.NewEncoder(w).Encode(map[string]any{
			"device_code":      "test-device-code",
			"user_code":        "WDJB-MJHT",
			"verification_uri": base + "/activate",
			"expires_in":       300,
			"interval":         1,
		})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		idToken := m.signIDToken(t)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "test-access-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"id_token":     idToken,
		})
	})

	m.server = httptest.NewServer(mux)
	t.Cleanup(m.server.Close)
	return m
}

// signIDToken builds and signs an id_token from the current knob values.
func (m *mockIDP) signIDToken(t *testing.T) string {
	t.Helper()
	aud := m.audience
	if aud == "" {
		aud = m.clientID
	}
	iss := m.issuerOverride
	if iss == "" {
		iss = m.server.URL
	}
	now := time.Now()
	claims := map[string]any{
		"iss":    iss,
		"aud":    aud,
		"sub":    "subject-123",
		"email":  m.email,
		"groups": m.groups,
		"iat":    now.Unix(),
		"exp":    now.Add(m.lifetime).Unix(),
	}
	raw, err := jwt.Signed(m.signer).Claims(claims).Serialize()
	require.NoError(t, err)
	return raw
}

// newAuthenticatorForIDP builds an oidcAuthenticator wired to the mock IdP: a
// real go-oidc provider (so the real verifier and JWKS fetch are used) and an
// oauth2 config whose endpoints point at the mock. ds is the real test DB.
func newAuthenticatorForIDP(t *testing.T, idp *mockIDP, ds sqlutil.DataSource) *oidcAuthenticator {
	t.Helper()
	ctx := context.Background()
	provider, err := gooidc.NewProvider(ctx, idp.server.URL)
	require.NoError(t, err)

	return &oidcAuthenticator{
		ds:           ds,
		config:       &TestConfig{},
		provider:     provider,
		oidcConfig:   &gooidc.Config{ClientID: idp.clientID},
		oauth2Config: &oauth2.Config{ClientID: idp.clientID, Endpoint: provider.Endpoint()},
		lggr:         logger.TestLogger(t),
		auditLogger:  &audit.AuditLoggerService{},
		deviceFlows:  newDeviceFlowStore(),
	}
}
