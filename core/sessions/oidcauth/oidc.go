/*
The OIDC module handles implementation of the initial auth flow redirects by requiring
conustruction with a reference to the HTTP router, verifying the attestation against
the config's Provider service, and then creating a sesion based on the OIDC response
'groups' claim, with roles mapped within the config.

This module configures and spins up its own gin http api router to handle the single callback endpoint.

NewApplication would have needed to rely on the created router engine, but the router
engine requires the application to have already been created. It works better to have a
second, standalone router and http listener initiated conditionally here when the OIDC
driver is used, where the callback routes are self contained and managed within this module.
*/
package oidcauth

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"golang.org/x/oauth2"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	RouterRateLimitterPeriod = 1 * time.Minute
	RouterRateLimitterLimit  = 1000
)

var ErrUserNoOIDCGroups = errors.New("user claims response from identity server recieved, but no matching role group names in claim")

type oidcAuthenticator struct {
	q            pg.Q
	config       config.OIDC
	provider     *oidc.Provider
	oidcConfig   *oidc.Config
	oauth2Config *oauth2.Config
	lggr         logger.Logger
	auditLogger  audit.AuditLogger
}

// oidcAuthenticator implements sessions.AuthenticationProvider interface
var _ sessions.AuthenticationProvider = (*oidcAuthenticator)(nil)

func NewOIDCAuthenticator(
	db *sqlx.DB,
	pgCfg pg.QConfig,
	oidcCfg config.OIDC,
	lggr logger.Logger,
	auditLogger audit.AuditLogger,
) (*oidcAuthenticator, error) {
	namedLogger := lggr.Named("OIDCAuthenticationProvider")

	// Ensure all RBAC role mappings to OIDC Groups are defined, and required fields populated, or error on startup
	if oidcCfg.AdminUserGroupClaim() == "" || oidcCfg.EditUserGroupClaim() == "" ||
		oidcCfg.RunUserGroupClaim() == "" || oidcCfg.ReadUserGroupClaim() == "" {
		return nil, errors.New("OIDC Group name mapping for callback group claims for all local RBAC role required. Set group names for `_UserGroupClaim` fields")
	}
	if oidcCfg.ClientID() == "" {
		return nil, errors.New("OIDC ClientID config required")
	}
	if oidcCfg.ClientSecret() == "" {
		return nil, errors.New("OIDC ClientSecret config required")
	}
	if oidcCfg.ProviderDomain() == "" {
		return nil, errors.New("OIDC ProviderDomain config required")
	}
	if oidcCfg.HTTPPort() == 0 {
		return nil, errors.New("OIDC HTTPPort config required")
	}

	var provider *oidc.Provider
	var oidcConfig *oidc.Config
	var oauth2Config *oauth2.Config

	ctx := context.Background()
	// Initialize provider based on config domain, this contains a blocking call to as part of the OpenID Connect discovery process
	provider, err := oidc.NewProvider(ctx, oidcCfg.ProviderDomain()+oidcCfg.OAuth2ProviderRouteSuffix())
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
		RedirectURL:  oidcCfg.OIDCCallbackURL() + oidcCfg.OIDCCallbackURLSuffix(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}

	// Create Authenticator struct, with internal HTTP handlers
	ldapAuth := oidcAuthenticator{
		q:            pg.NewQ(db, namedLogger, pgCfg),
		config:       oidcCfg,
		provider:     provider,
		oidcConfig:   oidcConfig,
		oauth2Config: oauth2Config,
		lggr:         lggr.Named("OIDCAuthenticationProvider"),
		auditLogger:  auditLogger,
	}

	// Create a new, separate gin engine to register and listen for the OIDC callback request containing
	// the user claims and groups, set up ratelimitter
	oidcCallbackEngine := gin.New()
	api := oidcCallbackEngine.Group("/", rateLimiter(RouterRateLimitterPeriod, RouterRateLimitterLimit))
	api.GET("/auth/oidc-login", ginHandlerFromHTTP(ldapAuth.handleLoginProviderRedirect))
	api.GET(oidcCfg.OIDCCallbackURLSuffix(), ginHandlerFromHTTP(ldapAuth.handleOIDCCallback))

	lggr.Infof("Initialized OIDC HTTP router and routes on %d", fmt.Sprintf("%d", oidcCfg.HTTPPort()))
	oidcCallbackEngine.Run(":" + fmt.Sprintf("%d", oidcCfg.HTTPPort()))

	return &ldapAuth, nil
}

func (oi *oidcAuthenticator) handleLoginProviderRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oi.oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline), http.StatusFound)
}

func (oi *oidcAuthenticator) handleOIDCCallback(w http.ResponseWriter, r *http.Request) {
	// Verify initial error or required query params
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		oi.lggr.Warnf("Recieved error in OIDC response: %s", errMsg)
		http.Error(w, "Error in OIDC response", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	// Begin token exchange to retrieve attested claims of authenticated user
	oauth2Token, err := oi.oauth2Config.Exchange(ctx, code)
	if err != nil {
		oi.lggr.Errorf("Failed to exchange token: %v", err)
		http.Error(w, "OIDC exchange failed", http.StatusInternalServerError)
		return
	}

	// Request token from provider for claims lookup and verification
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		oi.lggr.Errorf("No id_token field in oauth2 token: %v", err)
		http.Error(w, "Missing id_token field in response", http.StatusInternalServerError)
		return
	}

	// Verify claim and retrieve attested user groups
	idToken, err := oi.provider.Verifier(oi.oidcConfig).Verify(ctx, rawIDToken)
	if err != nil {
		oi.lggr.Errorf("Failed to verify ID token: %v", err)
		http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
		return
	}

	var claims struct {
		Email      string   `json:"email"`
		Groups     []string `json:"groups"`
		UserGroups []string `json:"userGroups"`
		exp        []string `json:"exp"`
	}
	if err := idToken.Claims(&claims); err != nil {
		oi.lggr.Errorf("Failed to parse OIDC return claims: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	oi.lggr.Infof("Recieved and validated OIDC claims: %v\n", claims)

	// Map the groups and insert a newly created session paired with role mapping for user
	role, err := groupClaimsToUserRole(
		claims.Groups,
		oi.config.AdminUserGroupClaim(),
		oi.config.EditUserGroupClaim(),
		oi.config.RunUserGroupClaim(),
		oi.config.ReadUserGroupClaim(),
	)
	if err != nil {
		oi.lggr.Errorf("Failed to map configured RBAC role name against recieved list of group claims: %v", err)
		http.Error(w, "No matching role within attested user group claims", http.StatusBadRequest)
		return
	}

	// Save new user authenticated session and role to oidc_sessions table
	// Sessions are set to expire after the duration + creation date elapsed
	session := sessions.NewSession()
	_, err = oi.q.Exec(
		"INSERT INTO oidc_sessions (id, user_email, user_role, created_at) VALUES ($1, $2, $3, $4, now())",
		session.ID,
		strings.ToLower(claims.Email),
		role,
	)
	if err != nil {
		oi.lggr.Errorf("unable to create new session in oidc_sessions table %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
	}

	oi.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]interface{}{"email": claims.Email})

	// Redirect to operator UI
	// Set authenticated response and session cookie
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(map[string]interface{}{
		"Authenticated": true,
		"Session":       session.ID,
	})
	w.Write(jsonData)
}

// FindUser in the context of the OIDC driver only supports local admin users
func (oi *oidcAuthenticator) FindUser(email string) (sessions.User, error) {
	email = strings.ToLower(email)
	foundUser := sessions.User{}

	var foundLocalAdminUser sessions.User
	checkErr := oi.q.Transaction(func(tx pg.Queryer) error {
		sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
		return tx.Get(&foundLocalAdminUser, sql, email)
	})
	if checkErr == nil {
		return foundLocalAdminUser, nil
	}
	// If error is not nil, there was either an issue or no local users found
	if !errors.Is(checkErr, sql.ErrNoRows) {
		// If the error is not that no local user was found, log and exit
		oi.lggr.Errorf("error searching users table: %v", checkErr)
		return sessions.User{}, errors.New("error Finding user")
	}

	return foundUser, nil
}

// FindUserByAPIToken retrieves a possible stored user and role from the oidc_user_api_tokens table store
func (oi *oidcAuthenticator) FindUserByAPIToken(apiToken string) (sessions.User, error) {
	if !oi.config.UserApiTokenEnabled() {
		return sessions.User{}, errors.New("API token is not enabled ")
	}

	var foundUser sessions.User
	err := oi.q.Transaction(func(tx pg.Queryer) error {
		// Query the oidc user API token table for given token, user role and email are cached so
		// no further upstream OIDC query is performed, sessions and tokens are synced against the upstream server
		// via the UpstreamSyncInterval config and reaper.go sync implementation
		var foundUserToken struct {
			UserEmail string
			UserRole  sessions.UserRole
			Valid     bool
		}
		if err := tx.Get(&foundUserToken,
			"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM oidc_user_api_tokens WHERE token_key = $1",
			apiToken, oi.config.UserAPITokenDuration().Duration(),
		); err != nil {
			return err
		}
		if !foundUserToken.Valid {
			return sessions.ErrUserSessionExpired
		}
		foundUser = sessions.User{
			Email: foundUserToken.UserEmail,
			Role:  foundUserToken.UserRole,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, sessions.ErrUserSessionExpired) {
			// API Token expired, purge
			if _, execErr := oi.q.Exec("DELETE FROM oidc_user_api_tokens WHERE token_key = $1", apiToken); err != nil {
				oi.lggr.Errorf("error purging stale oidc API token session: %v", execErr)
			}
		}
		return sessions.User{}, err
	}
	return foundUser, nil
}

// ListUsers in the context of the OIDC driver only supports listing the local (admin) users, we don't have an identity server to query against
func (oi *oidcAuthenticator) ListUsers() ([]sessions.User, error) {
	returnUsers := []sessions.User{}
	if err := oi.q.Transaction(func(tx pg.Queryer) error {
		sql := "SELECT * FROM users ORDER BY email ASC;"
		return tx.Select(&returnUsers, sql)
	}); err != nil {
		oi.lggr.Errorf("error listing local users: ", err)
	}
	return returnUsers, nil
}

// AuthorizedUserWithSession will return the API user associated with the Session ID if it
// exists and hasn't expired, and update session's LastUsed field
func (oi *oidcAuthenticator) AuthorizedUserWithSession(sessionID string) (sessions.User, error) {
	if len(sessionID) == 0 {
		return sessions.User{}, errors.New("session ID cannot be empty")
	}
	var foundUser sessions.User
	err := oi.q.Transaction(func(tx pg.Queryer) error {
		// Query the oidc_sessions table for given session ID, user role and email are saved after the SAML groups claim is provided and validated
		var foundSession struct {
			UserEmail string
			UserRole  sessions.UserRole
			Valid     bool
		}
		if err := tx.Get(&foundSession,
			"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM oidc_sessions WHERE id = $1",
			sessionID, oi.config.SessionTimeout().Duration(),
		); err != nil {
			return sessions.ErrUserSessionExpired
		}
		if !foundSession.Valid {
			// Sessions expired, purge
			return sessions.ErrUserSessionExpired
		}
		foundUser = sessions.User{
			Email: foundSession.UserEmail,
			Role:  foundSession.UserRole,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, sessions.ErrUserSessionExpired) {
			if _, execErr := oi.q.Exec("DELETE FROM oidc_sessions WHERE id = $1", sessionID); err != nil {
				oi.lggr.Errorf("error purging stale OIDC session: %v", execErr)
			}
		}
		return sessions.User{}, err
	}
	return foundUser, nil
}

// DeleteUser is not supported for read only OIDC
func (oi *oidcAuthenticator) DeleteUser(email string) error {
	return sessions.ErrNotSupported
}

// DeleteUserSession removes an oidcSession table entry by ID
func (oi *oidcAuthenticator) DeleteUserSession(sessionID string) error {
	_, err := oi.q.Exec("DELETE FROM oidc_sessions WHERE id = $1", sessionID)
	return err
}

// GetUserWebAuthn returns an empty stub, MFA is delegated to SAML provider
func (oi *oidcAuthenticator) GetUserWebAuthn(email string) ([]sessions.WebAuthn, error) {
	return []sessions.WebAuthn{}, nil
}

// CreateSession in the context of the OIDC driver handles only the local auth admin user, exposed by the default endpoint defined in the router. To initiate the SAML/OIDC
// flow, a separate /oidc-login route is defined which handles the redirect to the
// configured provider
func (oi *oidcAuthenticator) CreateSession(sr sessions.SessionRequest) (string, error) {
	foundUser, err := oi.localLoginFallback(sr)
	if err != nil {
		return "", err
	}

	oi.lggr.Infof("Successful local admin login request for user %s - %s", sr.Email, foundUser.Role)

	// Save local admin session, user, and role to sessions table
	// Sessions are set to expire after the duration + creation date elapsed
	session := sessions.NewSession()
	_, err = oi.q.Exec(
		"INSERT INTO oidc_sessions (id, user_email, user_role, created_at) VALUES ($1, $2, $3, $4, now())",
		session.ID,
		strings.ToLower(sr.Email),
		foundUser.Role,
	)
	if err != nil {
		oi.lggr.Errorf("unable to create new session in oidc_sessions table %v", err)
		return "", fmt.Errorf("error creating local LDAP session: %w", err)
	}

	oi.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]interface{}{"email": sr.Email})

	return session.ID, nil
}

// ClearNonCurrentSessions removes all oicd_sessions but the id passed in.
func (oi *oidcAuthenticator) ClearNonCurrentSessions(sessionID string) error {
	_, err := oi.q.Exec("DELETE FROM oicd_sessions where id != $1", sessionID)
	return err
}

// CreateUser is not supported for read only OIDC
func (oi *oidcAuthenticator) CreateUser(user *sessions.User) error {
	return sessions.ErrNotSupported
}

// UpdateRole is not supported for read only OIDC
func (oi *oidcAuthenticator) UpdateRole(email, newRole string) (sessions.User, error) {
	return sessions.User{}, sessions.ErrNotSupported
}

// SetPassword for remote users is not supported via the read only OIDC implementation, however change password
// in the context of updating a local admin user's password is required
func (oi *oidcAuthenticator) SetPassword(user *sessions.User, newPassword string) error {
	// Ensure specified user is part of the local admins user table
	var localAdminUser sessions.User
	if err := oi.q.Transaction(func(tx pg.Queryer) error {
		sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
		return tx.Get(&localAdminUser, sql, user.Email)
	}); err != nil {
		oi.lggr.Infof("Can not change password, local user with email not found in users table: %s, err: %v", user.Email, err)
		return sessions.ErrNotSupported
	}

	// User is local admin, save new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := oi.q.Transaction(func(tx pg.Queryer) error {
		sql := "UPDATE users SET hashed_password = $1, updated_at = now() WHERE email = $2 RETURNING *"
		return tx.Get(user, sql, hashedPassword, user.Email)
	}); err != nil {
		oi.lggr.Errorf("unable to set password for user: %s, err: %v", user.Email, err)
		return errors.New("unable to save password")
	}
	return nil
}

// TestPassword only supports the potential local admin user, as there is no queryable identity server for the OIDC implementation
func (oi *oidcAuthenticator) TestPassword(email string, password string) error {
	// Fall back to test local users table in case of supported local CLI users as well
	var hashedPassword string
	if err := oi.q.Get(&hashedPassword, "SELECT hashed_password FROM users WHERE lower(email) = lower($1)", email); err != nil {
		return errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, hashedPassword) {
		return errors.New("invalid credentials")
	}
	return nil
}

// CreateAndSetAuthToken generates a new credential token with the user role
func (oi *oidcAuthenticator) CreateAndSetAuthToken(user *sessions.User) (*auth.Token, error) {
	newToken := auth.NewToken()
	err := oi.SetAuthToken(user, newToken)
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

// SetAuthToken updates the user to use the given Authentication Token.
func (oi *oidcAuthenticator) SetAuthToken(user *sessions.User, token *auth.Token) error {
	if !oi.config.UserApiTokenEnabled() {
		return errors.New("API token is not enabled ")
	}

	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return fmt.Errorf("OIDCAuth SetAuthToken hashed secret error: %w", err)
	}

	err = oi.q.Transaction(func(tx pg.Queryer) error {
		// Remove any existing API tokens
		if _, err = oi.q.Exec("DELETE FROM oidc_user_api_tokens WHERE user_email = $1", user.Email); err != nil {
			return fmt.Errorf("error executing DELETE FROM oidc_user_api_tokens: %w", err)
		}
		// Create new API token for user
		_, err = oi.q.Exec(
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
		return errors.New("error creating API token")
	}

	oi.auditLogger.Audit(audit.APITokenCreated, map[string]interface{}{"user": user.Email})
	return nil
}

// DeleteAuthToken clears and disables the users Authentication Token.
func (oi *oidcAuthenticator) DeleteAuthToken(user *sessions.User) error {
	_, err := oi.q.Exec("DELETE FROM oidc_user_api_tokens WHERE email = $1")
	return err
}

// SaveWebAuthn is not supported for read only OIDC
func (oi *oidcAuthenticator) SaveWebAuthn(token *sessions.WebAuthn) error {
	return sessions.ErrNotSupported
}

// Sessions returns all sessions limited by the parameters.
func (oi *oidcAuthenticator) Sessions(offset, limit int) ([]sessions.Session, error) {
	var sessions []sessions.Session
	sql := `SELECT * FROM oidc_sessions ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err := oi.q.Select(&sessions, sql, limit, offset); err != nil {
		return sessions, nil
	}
	return sessions, nil
}

// FindExternalInitiator supports the 'Run' role external intiator header auth functionality
func (oi *oidcAuthenticator) FindExternalInitiator(eia *auth.Token) (*bridges.ExternalInitiator, error) {
	exi := &bridges.ExternalInitiator{}
	err := oi.q.Get(exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}

// localLoginFallback tests the credentials provided against the 'local' authentication method
// This covers the case of local CLI API calls requiring local login separate from the OIDC server
func (oi *oidcAuthenticator) localLoginFallback(sr sessions.SessionRequest) (sessions.User, error) {
	var user sessions.User
	sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
	err := oi.q.Get(&user, sql, sr.Email)
	if err != nil {
		return user, err
	}
	if !constantTimeEmailCompare(strings.ToLower(sr.Email), strings.ToLower(user.Email)) {
		oi.auditLogger.Audit(audit.AuthLoginFailedEmail, map[string]interface{}{"email": sr.Email})
		return user, errors.New("invalid email")
	}

	if !utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		oi.auditLogger.Audit(audit.AuthLoginFailedPassword, map[string]interface{}{"email": sr.Email})
		return user, errors.New("invalid password")
	}

	return user, nil
}

func groupClaimsToUserRole(oidcGroupClaims []string, adminGroupName string, editGroupName string, runGroupName string, readGroupName string) (sessions.UserRole, error) {
	// If defined Admin group name is present in groups claim, return UserRoleAdmin
	for _, group := range oidcGroupClaims {
		if group == adminGroupName {
			return sessions.UserRoleAdmin, nil
		}
	}
	// Check edit role
	for _, group := range oidcGroupClaims {
		if group == editGroupName {
			return sessions.UserRoleEdit, nil
		}
	}
	// Check run role
	for _, group := range oidcGroupClaims {
		if group == runGroupName {
			return sessions.UserRoleRun, nil
		}
	}
	// Check view role
	for _, group := range oidcGroupClaims {
		if group == readGroupName {
			return sessions.UserRoleView, nil
		}
	}
	// No role group found, error
	return sessions.UserRoleView, ErrUserNoOIDCGroups
}

const constantTimeEmailLength = 256

func constantTimeEmailCompare(left, right string) bool {
	length := mathutil.Max(constantTimeEmailLength, len(left), len(right))
	leftBytes := make([]byte, length)
	rightBytes := make([]byte, length)
	copy(leftBytes, left)
	copy(rightBytes, right)
	return subtle.ConstantTimeCompare(leftBytes, rightBytes) == 1
}

func ginHandlerFromHTTP(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func rateLimiter(period time.Duration, limit int64) gin.HandlerFunc {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}
	return mgin.NewMiddleware(limiter.New(store, rate))
}
