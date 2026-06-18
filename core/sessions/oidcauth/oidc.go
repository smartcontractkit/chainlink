/*
The OIDC module handles authentication by redirecting to a
Open ID Connect Identity Provider.
*/
package oidcauth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	webauth "github.com/smartcontractkit/chainlink/v2/core/web/auth"
)

const (
	SQLSelectUserbyEmail = "SELECT * FROM users WHERE lower(email) = lower($1)"
)

var ErrUserNoOIDCGroups = errors.New("user claims response from identity server received, but no matching role group names in claim")

type oidcAuthenticator struct {
	ds           sqlutil.DataSource
	config       config.OIDC
	provider     *oidc.Provider
	oidcConfig   *oidc.Config
	oauth2Config *oauth2.Config
	lggr         logger.Logger
	auditLogger  audit.AuditLogger
	deviceFlows  *deviceFlowStore
}

// ExchangeTokenRequest represents the expected JSON payload from the frontend
type ExchangeTokenRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"`
}

// ExchangeTokenResponse represents the response sent to the frontend
type ExchangeTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// oidcAuthenticator implements sessions.AuthenticationProvider interface
var _ clsessions.AuthenticationProvider = (*oidcAuthenticator)(nil)

// validateOIDCConfig checks the required OIDC fields at startup. ClientSecret is
// intentionally NOT required: confidential clients (Okta "Web" apps) set it and
// it is sent on the token exchange, while public clients (Okta "Native" apps,
// which the CLI device flow needs) leave it empty and rely on PKCE instead.
func validateOIDCConfig(oidcCfg config.OIDC) error {
	if oidcCfg.AdminClaim() == "" || oidcCfg.EditClaim() == "" ||
		oidcCfg.RunClaim() == "" || oidcCfg.ReadClaim() == "" {
		return errors.New("OIDC Group name mapping for callback group claims for all local RBAC role required. Set group names for `_Claim` fields")
	}
	if oidcCfg.ClientID() == "" {
		return errors.New("OIDC ClientID config required")
	}
	if oidcCfg.ProviderURL() == "" {
		return errors.New("OIDC ProviderURL config required")
	}
	if oidcCfg.RedirectURL() == "" {
		return errors.New("OIDC RedirectURL config required")
	}
	if oidcCfg.ClaimName() == "" {
		return errors.New("OIDC ClaimName config required")
	}
	return nil
}

func NewOIDCAuthenticator(
	ds sqlutil.DataSource,
	oidcCfg config.OIDC,
	lggr logger.Logger,
	auditLogger audit.AuditLogger,
) (*oidcAuthenticator, error) {
	// Ensure all RBAC role mappings to OIDC Id claims are defined, and required fields populated, or error on startup
	lggr.Debugf("OIDC CFG:\n %#v\n", oidcCfg)
	if err := validateOIDCConfig(oidcCfg); err != nil {
		return nil, err
	}

	var provider *oidc.Provider
	var oidcConfig *oidc.Config
	var oauth2Config *oauth2.Config

	ctx := context.Background()
	// Initialize provider based on config domain, this contains a blocking call to as part of the OpenID Connect discovery process
	provider, err := oidc.NewProvider(ctx, oidcCfg.ProviderURL())
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
		deviceFlows:  newDeviceFlowStore(),
	}

	return &oidcAuth, nil
}

func (oi *oidcAuthenticator) generateState() string {
	b := make([]byte, 32) // 256 bits of entropy
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		oi.lggr.Fatalf("failed to generate random bytes: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)[:43]
}

func (oi *oidcAuthenticator) handleCheckEnabled(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"enabled": true})
}

func (oi *oidcAuthenticator) handleSignIn(c *gin.Context) {
	// generate state and a PKCE verifier, store both on the session. PKCE (RFC
	// 7636) is sent unconditionally: it is additive hardening for confidential
	// clients and the sole proof-of-possession for public clients that carry no
	// secret. The challenge derived from the verifier travels to the provider;
	// the verifier stays server-side until the exchange.
	state := oi.generateState()
	verifier := oauth2.GenerateVerifier()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("pkce_verifier", verifier)
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// redirect to provider
	url := oi.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	c.Redirect(http.StatusFound, url)
}

func (oi *oidcAuthenticator) handleTokenExchange(c *gin.Context) {
	// parse and validate the incoming JSON request
	var req ExchangeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExchangeTokenResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// check state matches stored value on the session
	ginSession := sessions.Default(c)
	storedState := ginSession.Get("state")
	if storedState == nil || req.State != storedState.(string) {
		c.JSON(http.StatusBadRequest, ExchangeTokenResponse{
			Success: false,
			Message: "Invalid state parameter",
		})
		return
	}
	ginSession.Delete("state")

	// Recover the PKCE verifier stored at sign-in and bind it to this exchange.
	storedVerifier, _ := ginSession.Get("pkce_verifier").(string)
	ginSession.Delete("pkce_verifier")
	if storedVerifier == "" {
		c.JSON(http.StatusBadRequest, ExchangeTokenResponse{
			Success: false,
			Message: "Missing PKCE verifier",
		})
		return
	}

	// Begin token exchange to retrieve attested claims of authenticated user.
	// VerifierOption proves possession of the value whose challenge was sent at
	// sign-in; for public clients this replaces the client secret.
	ctx := context.Background()
	oauth2Token, err := oi.oauth2Config.Exchange(ctx, req.Code, oauth2.VerifierOption(storedVerifier))
	if err != nil {
		oi.lggr.Errorf("Failed to exchange token: %v", err)
		c.JSON(http.StatusInternalServerError, ExchangeTokenResponse{
			Success: false,
			Message: "OIDC exchange failed",
		})
		return
	}

	// Request token from provider for claims lookup and verification
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		oi.lggr.Errorf("No id_token field in oauth2 token: %v", err)
		c.String(http.StatusInternalServerError, "Missing id_token field in response")
		return
	}

	sessionID, _, err := oi.issueSessionFromIDToken(ctx, rawIDToken)
	if err != nil {
		oi.handleIssueSessionError(c, err)
		return
	}

	// save session
	ginSession.Set(webauth.SessionIDKey, sessionID)
	err = ginSession.Save()
	if err != nil {
		oi.lggr.Errorf("failed to saved session %v", err)
		c.String(http.StatusInternalServerError, "Authentication failed")
		return
	}

	c.JSON(http.StatusOK, ExchangeTokenResponse{
		Success: true,
	})
}

// errNoMatchingRole signals that a verified token carried no group claim that
// maps to a local RBAC role. Callers translate it to a 4xx (the token was
// valid, the user simply lacks an authorised group), distinct from a 5xx for
// genuine verification or storage failures.
var errNoMatchingRole = errors.New("no matching role within attested user group claims")

// issueSessionFromIDToken is the single trusted path that turns a raw OIDC
// id_token into an oidc_sessions row. Both the browser auth-code exchange and
// the CLI device flow funnel through here so verification, claim extraction,
// role mapping, and session creation never diverge between the two. It returns
// the new session ID and the authenticated email.
func (oi *oidcAuthenticator) issueSessionFromIDToken(ctx context.Context, rawIDToken string) (string, string, error) {
	// Verify signature, issuer, expiry, and audience (pinned to ClientID).
	idToken, err := oi.provider.Verifier(oi.oidcConfig).Verify(ctx, rawIDToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to verify ID token: %w", err)
	}

	var claims map[string]any
	if err = idToken.Claims(&claims); err != nil {
		return "", "", fmt.Errorf("failed to parse OIDC return claims: %w", err)
	}

	idClaims, err := oi.ExtractIDClaimValues(claims, oi.config.ClaimName())
	if err != nil {
		return "", "", fmt.Errorf("failed to extract ID claims for claim name '%s': %w", oi.config.ClaimName(), err)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", errors.New("failed to get email from claims")
	}
	oi.lggr.Tracef("Received and validated ID claims: %v\n", idClaims)

	role, err := oi.IDClaimsToUserRole(
		idClaims,
		oi.config.AdminClaim(),
		oi.config.EditClaim(),
		oi.config.RunClaim(),
		oi.config.ReadClaim(),
	)
	if err != nil {
		// Audit the full claim set so an over-broad or unexpected group grant is
		// detectable after the fact.
		oi.auditLogger.Audit(audit.AuthLoginFailed2FA, map[string]any{"email": email, "groups": idClaims})
		return "", "", fmt.Errorf("%w: %v", errNoMatchingRole, err)
	}

	clSession := clsessions.NewSession()
	_, err = oi.ds.ExecContext(
		ctx,
		"INSERT INTO oidc_sessions (id, user_email, user_role, created_at) VALUES ($1, $2, $3, now())",
		clSession.ID,
		strings.ToLower(email),
		role,
	)
	if err != nil {
		return "", "", fmt.Errorf("unable to create new session in oidc_sessions table: %w", err)
	}

	oi.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]any{"email": email, "role": role})
	return clSession.ID, email, nil
}

// handleIssueSessionError maps an issueSessionFromIDToken error to the HTTP
// status the browser exchange returns: a missing/invalid role is a client-side
// 400, everything else is a 500.
func (oi *oidcAuthenticator) handleIssueSessionError(c *gin.Context, err error) {
	if errors.Is(err, errNoMatchingRole) {
		oi.lggr.Errorf("Failed to map configured RBAC role name against received list of group claims: %v", err)
		c.String(http.StatusBadRequest, "No matching role within attested user group claims")
		return
	}
	oi.lggr.Errorf("Failed to establish OIDC session: %v", err)
	c.String(http.StatusInternalServerError, "Failed to establish OIDC session")
}

// FindUser in the context of the OIDC driver only supports local admin users
func (oi *oidcAuthenticator) FindUser(ctx context.Context, email string) (clsessions.User, error) {
	email = strings.ToLower(email)

	var foundUser clsessions.User

	if err := oi.ds.GetContext(ctx, &foundUser, SQLSelectUserbyEmail, email); err != nil {
		// If the error is not that no local user was found, log and exit
		if errors.Is(err, sql.ErrNoRows) {
			return clsessions.User{}, errors.New("user not found")
		}

		oi.lggr.Errorf("error searching users table: %v", err)
		return clsessions.User{}, errors.New("error finding user")
	}

	return foundUser, nil
}

// FindUserByAPIToken retrieves a possible stored user and role from the oidc_user_api_tokens table store
func (oi *oidcAuthenticator) FindUserByAPIToken(ctx context.Context, apiToken string) (clsessions.User, error) {
	if !oi.config.UserAPITokenEnabled() {
		return clsessions.User{}, errors.New("API token is not enabled")
	}

	var foundUser clsessions.User
	err := sqlutil.TransactDataSource(ctx, oi.ds, nil, func(tx sqlutil.DataSource) error {
		// Query the oidc user API token table for given token, user role and email are cached so
		// no further upstream OIDC query is performed, sessions and tokens are synced against the upstream server
		// via the UpstreamSyncInterval config and reaper.go sync implementation
		var foundUserToken struct {
			UserEmail string
			UserRole  clsessions.UserRole
			Valid     bool
		}
		if err := tx.GetContext(ctx, &foundUserToken,
			"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM oidc_user_api_tokens WHERE token_key = $1",
			apiToken, oi.config.UserAPITokenDuration().Duration(),
		); err != nil {
			return err
		}
		if !foundUserToken.Valid {
			return clsessions.ErrUserSessionExpired
		}
		foundUser = clsessions.User{
			Email: foundUserToken.UserEmail,
			Role:  foundUserToken.UserRole,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, clsessions.ErrUserSessionExpired) {
			// API Token expired, purge
			if _, execErr := oi.ds.ExecContext(ctx, "DELETE FROM oidc_user_api_tokens WHERE token_key = $1", apiToken); execErr != nil {
				oi.lggr.Errorf("error purging stale oidc API token session: %v", execErr)
			}
		}
		return clsessions.User{}, err
	}
	return foundUser, nil
}

// ListUsers in the context of the OIDC driver only supports listing the local (admin) users, we don't have an identity server to query against
func (oi *oidcAuthenticator) ListUsers(ctx context.Context) ([]clsessions.User, error) {
	returnUsers := []clsessions.User{}
	if err := oi.ds.SelectContext(ctx, &returnUsers, "SELECT * FROM users ORDER BY email ASC;"); err != nil {
		oi.lggr.Errorf("error listing local users: %v", err)
	}
	return returnUsers, nil
}

// AuthorizedUserWithSession will return the API user associated with the Session ID if it
// exists and hasn't expired
func (oi *oidcAuthenticator) AuthorizedUserWithSession(ctx context.Context, sessionID string) (clsessions.User, error) {
	if len(sessionID) == 0 {
		return clsessions.User{}, errors.New("session ID cannot be empty")
	}
	var foundUser clsessions.User
	err := sqlutil.TransactDataSource(ctx, oi.ds, nil, func(tx sqlutil.DataSource) error {
		// Query the oidc_sessions table for given session ID, user role and email are saved after the id claims is provided and validated
		var foundSession struct {
			UserEmail string
			UserRole  clsessions.UserRole
			Valid     bool
		}
		if err := tx.GetContext(ctx, &foundSession,
			"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM oidc_sessions WHERE id = $1",
			sessionID, oi.config.SessionTimeout().Duration(),
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return clsessions.ErrUserSessionExpired
			}
			return err
		}
		if !foundSession.Valid {
			// Sessions expired, purge
			return clsessions.ErrUserSessionExpired
		}
		foundUser = clsessions.User{
			Email: foundSession.UserEmail,
			Role:  foundSession.UserRole,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, clsessions.ErrUserSessionExpired) {
			if _, execErr := oi.ds.ExecContext(ctx, "DELETE FROM oidc_sessions WHERE id = $1", sessionID); execErr != nil {
				oi.lggr.Errorf("error purging stale OIDC session: %v", execErr)
			}
		}
		return clsessions.User{}, err
	}
	return foundUser, nil
}

// DeleteUser is not supported for read only OIDC
func (oi *oidcAuthenticator) DeleteUser(ctx context.Context, email string) error {
	return clsessions.ErrNotSupported
}

// DeleteUserSession removes an oidcSession table entry by ID
func (oi *oidcAuthenticator) DeleteUserSession(ctx context.Context, sessionID string) error {
	_, err := oi.ds.ExecContext(ctx, "DELETE FROM oidc_sessions WHERE id = $1", sessionID)
	return err
}

// GetUserWebAuthn returns an empty stub, MFA is delegated to SAML provider
func (oi *oidcAuthenticator) GetUserWebAuthn(ctx context.Context, email string) ([]clsessions.WebAuthn, error) {
	return []clsessions.WebAuthn{}, nil
}

// CreateSession in the context of the OIDC driver handles only the local auth admin user, exposed by the default endpoint defined in the router. To initiate the SAML/OIDC
// flow, a separate /oidc-login route is defined which handles the redirect to the
// configured provider
func (oi *oidcAuthenticator) CreateSession(ctx context.Context, sr clsessions.SessionRequest) (string, error) {
	foundUser, err := oi.localLoginFallback(ctx, sr)
	if err != nil {
		return "", err
	}

	sanitizedEmail := strings.ReplaceAll(sr.Email, "\n", "")
	sanitizedEmail = strings.ReplaceAll(sanitizedEmail, "\r", "")
	oi.lggr.Infof("Successful local admin login request for user %s - %s", sanitizedEmail, foundUser.Role)

	// Save local admin session, user, and role to sessions table
	// Sessions are set to expire after the duration + creation date elapsed
	session := clsessions.NewSession()
	_, err = oi.ds.ExecContext(ctx,
		"INSERT INTO oidc_sessions (id, user_email, user_role, created_at) VALUES ($1, $2, $3, now())",
		session.ID,
		strings.ToLower(sr.Email),
		foundUser.Role,
	)
	if err != nil {
		oi.lggr.Errorf("unable to create new session in oidc_sessions table %v", err)
		return "", fmt.Errorf("error creating local OIDC session: %w", err)
	}

	oi.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]any{"email": sr.Email})

	return session.ID, nil
}

// ClearNonCurrentSessions removes all oicd_sessions but the id passed in.
func (oi *oidcAuthenticator) ClearNonCurrentSessions(ctx context.Context, sessionID string) error {
	_, err := oi.ds.ExecContext(ctx, "DELETE FROM oidc_sessions where id != $1", sessionID)
	return err
}

// CreateUser is not supported for read only OIDC
func (oi *oidcAuthenticator) CreateUser(ctx context.Context, user *clsessions.User) error {
	return clsessions.ErrNotSupported
}

// UpdateRole is not supported for read only OIDC
func (oi *oidcAuthenticator) UpdateRole(ctx context.Context, email string, newRole string) (clsessions.User, error) {
	return clsessions.User{}, clsessions.ErrNotSupported
}

// SetPassword for remote users is not supported via the read only OIDC implementation, however change password
// in the context of updating a local admin user's password is required
func (oi *oidcAuthenticator) SetPassword(ctx context.Context, user *clsessions.User, newPassword string) error {
	// Ensure specified user is part of the local admins user table
	var localAdminUser clsessions.User
	if err := oi.ds.GetContext(ctx, &localAdminUser, SQLSelectUserbyEmail, user.Email); err != nil {
		oi.lggr.Infof("Can not change password, local user with email not found in users table: %s, err: %v", user.Email, err)
		return clsessions.ErrNotSupported
	}

	// User is local admin, save new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		oi.lggr.Errorf("Error hashing user password: err: %v", err)
		return errors.New("unable to hash password")
	}
	if err := oi.ds.GetContext(ctx, user,
		"UPDATE users SET hashed_password = $1, updated_at = now() WHERE email = $2 RETURNING *",
		hashedPassword, user.Email,
	); err != nil {
		oi.lggr.Errorf("unable to set password for user: %s, err: %v", user.Email, err)
		return errors.New("unable to save password")
	}
	return nil
}

// TestPassword only supports the potential local admin user, as there is no queryable identity server for the OIDC implementation
func (oi *oidcAuthenticator) TestPassword(ctx context.Context, email string, password string) error {
	// Fall back to test local users table in case of supported local CLI users as well
	var hashedPassword string
	if err := oi.ds.GetContext(ctx, &hashedPassword, "SELECT hashed_password FROM users WHERE lower(email) = lower($1)", email); err != nil {
		return errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, hashedPassword) {
		return errors.New("invalid credentials")
	}
	return nil
}

// CreateAndSetAuthToken generates a new credential token with the user role
func (oi *oidcAuthenticator) CreateAndSetAuthToken(ctx context.Context, user *clsessions.User) (*auth.Token, error) {
	newToken := auth.NewToken()
	err := oi.SetAuthToken(ctx, user, newToken)
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

// SetAuthToken updates the user to use the given Authentication Token.
func (oi *oidcAuthenticator) SetAuthToken(ctx context.Context, user *clsessions.User, token *auth.Token) error {
	if !oi.config.UserAPITokenEnabled() {
		return errors.New("API token is not enabled ")
	}

	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return fmt.Errorf("OIDCAuth SetAuthToken hashed secret error: %w", err)
	}

	err = sqlutil.TransactDataSource(ctx, oi.ds, nil, func(tx sqlutil.DataSource) error {
		// Remove any existing API tokens
		if _, err = oi.ds.ExecContext(ctx, "DELETE FROM oidc_user_api_tokens WHERE user_email = $1", user.Email); err != nil {
			return fmt.Errorf("error executing DELETE FROM oidc_user_api_tokens: %w", err)
		}
		// Create new API token for user
		_, err = oi.ds.ExecContext(ctx,
			"INSERT INTO oidc_user_api_tokens (user_email, user_role, token_key, token_salt, token_hashed_secret, created_at) VALUES ($1, $2, $3, $4, $5, $6, now())",
			user.Email,
			user.Role,
			token.AccessKey,
			salt,
			hashedSecret,
		)
		if err != nil {
			return fmt.Errorf("failed insert into oidc_user_api_tokens: %w", err)
		}
		return nil
	})
	if err != nil {
		oi.lggr.Errorf("error creating API token: %v", err)
		return errors.New("error creating API token")
	}

	oi.auditLogger.Audit(audit.APITokenCreated, map[string]any{"user": user.Email})
	return nil
}

// DeleteAuthToken clears and disables the users Authentication Token.
func (oi *oidcAuthenticator) DeleteAuthToken(ctx context.Context, user *clsessions.User) error {
	_, err := oi.ds.ExecContext(ctx, "DELETE FROM oidc_user_api_tokens WHERE user_email = $1", user.Email)
	return err
}

// SaveWebAuthn is not supported for read only OIDC
func (oi *oidcAuthenticator) SaveWebAuthn(ctx context.Context, token *clsessions.WebAuthn) error {
	return clsessions.ErrNotSupported
}

// Sessions returns all sessions limited by the parameters.
func (oi *oidcAuthenticator) Sessions(ctx context.Context, offset, limit int) ([]clsessions.Session, error) {
	var sessions []clsessions.Session
	sql := `SELECT * FROM oidc_sessions ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err := oi.ds.SelectContext(ctx, &sessions, sql, limit, offset); err != nil {
		return sessions, nil
	}
	return sessions, nil
}

// FindExternalInitiator supports the 'Run' role external intiator header auth functionality
func (oi *oidcAuthenticator) FindExternalInitiator(ctx context.Context, eia *auth.Token) (*bridges.ExternalInitiator, error) {
	exi := &bridges.ExternalInitiator{}
	err := oi.ds.GetContext(ctx, exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}

// localLoginFallback tests the credentials provided against the 'local' authentication method
// This covers the case of local CLI API calls requiring local login separate from the OIDC server
func (oi *oidcAuthenticator) localLoginFallback(ctx context.Context, sr clsessions.SessionRequest) (clsessions.User, error) {
	var user clsessions.User
	err := oi.ds.GetContext(ctx, &user, SQLSelectUserbyEmail, sr.Email)
	if err != nil {
		return user, err
	}
	if !constantTimeEmailCompare(strings.ToLower(sr.Email), strings.ToLower(user.Email)) {
		oi.auditLogger.Audit(audit.AuthLoginFailedEmail, map[string]any{"email": sr.Email})
		return user, errors.New("invalid email")
	}

	if !utils.CheckPasswordHash(sr.Password, string(user.HashedPassword)) {
		oi.auditLogger.Audit(audit.AuthLoginFailedPassword, map[string]any{"email": sr.Email})
		return user, errors.New("invalid password")
	}

	return user, nil
}

func (oi *oidcAuthenticator) IDClaimsToUserRole(idClaims []string, adminClaim string, editClaim string, runClaim string, readClaim string) (clsessions.UserRole, error) {
	// If defined Admin group name is present in id claims, return UserRoleAdmin
	if slices.Contains(idClaims, adminClaim) {
		return clsessions.UserRoleAdmin, nil
	}
	// Check edit role
	if slices.Contains(idClaims, editClaim) {
		return clsessions.UserRoleEdit, nil
	}
	// Check run role
	if slices.Contains(idClaims, runClaim) {
		return clsessions.UserRoleRun, nil
	}
	// Check view role
	if slices.Contains(idClaims, readClaim) {
		return clsessions.UserRoleView, nil
	}
	// No role group found, error
	return clsessions.UserRoleView, ErrUserNoOIDCGroups
}

// extractIDClaimValues extracts groups from the claims using the specified key
func (oi *oidcAuthenticator) ExtractIDClaimValues(claims map[string]any, key string) ([]string, error) {
	claimValues, ok := claims[key]
	if !ok {
		return nil, fmt.Errorf("claim '%s' not found in ID token", key)
	}

	// Handle different types of claim values
	switch v := claimValues.(type) {
	case []any:
		val := make([]string, 0, len(v))
		for _, item := range v {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type for item in '%s': expected string, got %T", key, item)
			}
			val = append(val, str)
		}
		return val, nil
	case []string:
		return v, nil
	case string:
		return []string{v}, nil
	default:
		return nil, fmt.Errorf("claim '%s' is not a string or array: got %T", key, v)
	}
}

func constantTimeEmailCompare(left, right string) bool {
	const constantTimeEmailLength = 256
	length := mathutil.Max(constantTimeEmailLength, len(left), len(right))
	leftBytes := make([]byte, length)
	rightBytes := make([]byte, length)
	copy(leftBytes, left)
	copy(rightBytes, right)
	return subtle.ConstantTimeCompare(leftBytes, rightBytes) == 1
}

func (oi *oidcAuthenticator) ExtendRouter(api *gin.RouterGroup) error {
	api.GET("/oidc-enabled", oi.handleCheckEnabled)
	api.GET("/oidc-login", oi.handleSignIn)
	api.POST("/oidc-exchange", oi.handleTokenExchange)
	// Device authorization grant for headless clients (the CLI).
	api.POST("/oidc-device/start", oi.handleDeviceStart)
	api.POST("/oidc-device/poll", oi.handleDevicePoll)

	return nil
}
