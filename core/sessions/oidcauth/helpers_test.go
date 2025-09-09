package oidcauth

import (
	"context"
	"log"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
)

// Returns an instantiated OIDCAuthenticator struct without validation for testing
func NewTestOIDCAuthenticator(
	ds sqlutil.DataSource,
	oidcCfg config.OIDC,
	lggr logger.Logger,
	auditLogger audit.AuditLogger,
) (*oidcAuthenticator, error) {
	var provider *oidc.Provider
	var oidcConfig *oidc.Config
	var oauth2Config *oauth2.Config
	ctx := context.Background()
	// Initialize provider based on config domain, this contains a blocking call to as part of the OpenID Connect discovery process
	provider, err := NewMockProvider(ctx, oidcCfg.ProviderURL())
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	// Construct oidc and oath callback configs for oidcAuth struct
	oidcConfig = &oidc.Config{
		ClientID: oidcCfg.ClientID(),
	}
	oauth2Config = &oauth2.Config{
		ClientID:     oidcCfg.ClientID(),
		ClientSecret: oidcCfg.ClientSecret(),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  oidcCfg.RedirectURL(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", oidcCfg.ClaimName()},
	}

	// Create Authenticator struct, with internal HTTP handlers
	oidcAuth := oidcAuthenticator{
		ds:           ds,
		config:       oidcCfg,
		provider:     provider,
		oidcConfig:   oidcConfig,
		oauth2Config: oauth2Config,
		lggr:         lggr.Named("OIDCAuthenticationProvider"),
		auditLogger:  auditLogger,
	}

	return &oidcAuth, nil
}

func NewMockProvider(ctx context.Context, issuer string) (*oidc.Provider, error) {
	x := oidc.ProviderConfig{
		IssuerURL:     issuer,
		AuthURL:       issuer + "/v1/authorize",
		TokenURL:      issuer + "/v1/token",
		DeviceAuthURL: issuer + "/v1/clients",
		UserInfoURL:   issuer + "/v1/userinfo",
		JWKSURL:       issuer + "/v1/keys",
		Algorithms:    []string{""},
	}

	return x.NewProvider(ctx), nil
}

// Default server group name mappings for test config and mocked ldap search results
const (
	ClaimName   = "groups"
	AdminClaim  = "NodeAdmins"
	EditorClaim = "NodeEditors"
	RunnerClaim = "NodeRunners"
	ReadClaim   = "NodeReadOnly"
)

// Implements config.OIDC
type TestConfig struct {
}

func (t *TestConfig) ClientID() string {
	return "abcd1234"
}

func (t *TestConfig) ClientSecret() string {
	return "abcd1234"
}

func (t *TestConfig) ProviderURL() string {
	return "https://id.example.com/oauth2/default"
}

func (t *TestConfig) RedirectURL() string {
	return "http://localhost:8080/signin"
}

func (t *TestConfig) ClaimName() string {
	return ClaimName
}

func (t *TestConfig) AdminClaim() string {
	return AdminClaim
}

func (t *TestConfig) EditClaim() string {
	return EditorClaim
}

func (t *TestConfig) RunClaim() string {
	return RunnerClaim
}

func (t *TestConfig) ReadClaim() string {
	return ReadClaim
}

func (t *TestConfig) SessionTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(15 * time.Minute)
}

func (t *TestConfig) UserAPITokenEnabled() bool {
	return true
}

func (t *TestConfig) UserAPITokenDuration() commonconfig.Duration {
	return *commonconfig.MustNewDuration(24 * time.Hour)
}
