/*
The LDAP authentication package forwards the credentials in the user session request
for authentication with a configured upstream LDAP server

This package relies on the two following local database tables:

	ldap_sessions: 	Upon successful LDAP response, creates a keyed local copy of the user email
	ldap_user_api_tokens: User created API tokens, tied to the node, storing user email.

Note: user can have only one API token at a time, and token expiration is enforced

User session and roles are cached and revalidated with the upstream service at the interval defined in
the local LDAP config through the Application.sessionReaper implementation in reaper.go.

Changes to the upstream identity server will propagate through and update local tables (web sessions, API tokens)
by either removing the entries or updating the roles. This sync happens for every auth endpoint hit, and
via the defined sync interval. One goroutine is created to coordinate the sync timing in the New function

This implementation is read only; user mutation actions such as Delete are not supported.

MFA is supported via the remote LDAP server implementation. Sufficient request time out should accommodate
for a blocking auth call while the user responds to a potential push notification callback.
*/
package ldapauth

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	UniqueMemberAttribute = "uniqueMember"
)

var ErrUserNotInUpstream = errors.New("LDAP query returned no matching users")
var ErrUserNoLDAPGroups = errors.New("user present in directory, but matching no role groups assigned")

type ldapAuthenticator struct {
	ds          sqlutil.DataSource
	ldapClient  LDAPClient
	config      config.LDAP
	lggr        logger.Logger
	auditLogger audit.AuditLogger
}

// ldapAuthenticator implements sessions.AuthenticationProvider interface
var _ sessions.AuthenticationProvider = (*ldapAuthenticator)(nil)

func NewLDAPAuthenticator(
	ds sqlutil.DataSource,
	ldapCfg config.LDAP,
	dev bool,
	lggr logger.Logger,
	auditLogger audit.AuditLogger,
) (*ldapAuthenticator, error) {
	// If not chainlink dev and not tls, error
	if !dev && !ldapCfg.ServerTLS() {
		return nil, errors.New("LDAP Authentication driver requires TLS when running in Production mode")
	}

	// Ensure all RBAC role mappings to LDAP Groups are defined, and required fields populated, or error on startup
	if ldapCfg.AdminUserGroupCN() == "" || ldapCfg.EditUserGroupCN() == "" ||
		ldapCfg.RunUserGroupCN() == "" || ldapCfg.ReadUserGroupCN() == "" {
		return nil, errors.New("LDAP Group mapping from server group name for all local RBAC role required. Set group names for `_UserGroupCN` fields")
	}
	if ldapCfg.ServerAddress() == "" {
		return nil, errors.New("LDAP ServerAddress config required")
	}
	if ldapCfg.ReadOnlyUserLogin() == "" {
		return nil, errors.New("LDAP ReadOnlyUserLogin config required")
	}

	ldapAuth := ldapAuthenticator{
		ds:          ds,
		ldapClient:  newLDAPClient(ldapCfg),
		config:      ldapCfg,
		lggr:        lggr.Named("LDAPAuthenticationProvider"),
		auditLogger: auditLogger,
	}

	// Single override of library defined global
	ldap.DefaultTimeout = ldapCfg.QueryTimeout()

	// Test initial connection and credentials
	lggr.Infof("Attempting initial connection to configured LDAP server with bind as API user")
	conn, err := ldapAuth.ldapClient.CreateEphemeralConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to establish connection to LDAP server with provided URL and credentials: %w", err)
	}
	conn.Close()

	// Store LDAP connection config for auth/new connection per request instead of persisted connection with reconnect
	return &ldapAuth, nil
}

// FindUser will attempt to return an LDAP user with mapped role by email.
func (l *ldapAuthenticator) FindUser(ctx context.Context, email string) (sessions.User, error) {
	email = strings.ToLower(email)

	// First check for the supported local admin users table
	var foundLocalAdminUser sessions.User
	checkErr := l.ds.GetContext(ctx, &foundLocalAdminUser, "SELECT * FROM users WHERE lower(email) = lower($1)", email)
	if checkErr == nil {
		return foundLocalAdminUser, nil
	}
	// If error is not nil, there was either an issue or no local users found
	if !errors.Is(checkErr, sql.ErrNoRows) {
		// If the error is not that no local user was found, log and exit
		l.lggr.Errorf("error searching users table: %v", checkErr)
		return sessions.User{}, errors.New("error Finding user")
	}

	// First query for user "is active" property if defined
	usersActive, err := l.validateUsersActive([]string{email})
	if err != nil {
		if errors.Is(err, ErrUserNotInUpstream) {
			return sessions.User{}, ErrUserNotInUpstream
		}
		l.lggr.Errorf("error in validateUsers call: %v", err)
		return sessions.User{}, errors.New("error running query to validate user active")
	}
	if !usersActive[0] {
		return sessions.User{}, errors.New("user not active")
	}

	conn, err := l.ldapClient.CreateEphemeralConnection()
	if err != nil {
		l.lggr.Errorf("error in LDAP dial: ", err)
		return sessions.User{}, errors.New("unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	// User email and role are the only upstream data that needs queried for.
	// List query user groups using the provided email, on success is a list of group the uniquemember belongs to
	// data is readily available
	escapedEmail := ldap.EscapeFilter(email)
	searchBaseDN := fmt.Sprintf("%s, %s", l.config.GroupsDN(), l.config.BaseDN())
	filterQuery := fmt.Sprintf("(&(uniquemember=%s=%s,%s,%s))", l.config.BaseUserAttr(), escapedEmail, l.config.UsersDN(), l.config.BaseDN())
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, int(l.config.QueryTimeout().Seconds()), false,
		filterQuery,
		[]string{"cn"},
		nil,
	)

	// Query the server
	result, err := conn.Search(searchRequest)
	if err != nil {
		l.lggr.Errorf("error searching users in LDAP query: %v", err)
		return sessions.User{}, errors.New("error searching users in LDAP directory")
	}

	if len(result.Entries) == 0 {
		// Provided email is not present in upstream LDAP server, local admin CLI auth is supported
		// So query and check the users table as well before failing
		var localUserRole sessions.UserRole
		if err = l.ds.GetContext(ctx, &localUserRole, "SELECT role FROM users WHERE email = $1", email); err != nil {
			// Above query for local user unsuccessful, return error
			l.lggr.Warnf("No local users table user found with email %s", email)
			return sessions.User{}, errors.New("no users found with provided email")
		}

		// If the above query to the local users table was successful, return that local user's role
		return sessions.User{
			Email: email,
			Role:  localUserRole,
		}, nil
	}

	// Populate found user by email and role based on matched group names
	userRole, err := l.groupSearchResultsToUserRole(result.Entries)
	if err != nil {
		l.lggr.Warnf("User '%s' found but no matching assigned groups in LDAP to assume role", email)
		return sessions.User{}, err
	}

	// Convert search result to sessions.User type with required fields
	return sessions.User{
		Email: email,
		Role:  userRole,
	}, nil
}

// FindUserByAPIToken retrieves a possible stored user and role from the ldap_user_api_tokens table store
func (l *ldapAuthenticator) FindUserByAPIToken(ctx context.Context, apiToken string) (sessions.User, error) {
	if !l.config.UserApiTokenEnabled() {
		return sessions.User{}, errors.New("API token is not enabled ")
	}

	// Query the ldap user API token table for given token, user role and email are cached so
	// no further upstream LDAP query is performed, sessions and tokens are synced against the upstream server
	// via the UpstreamSyncInterval config and reaper.go sync implementation
	var foundUserToken struct {
		UserEmail string
		UserRole  sessions.UserRole
		Valid     bool
	}
	err := l.ds.GetContext(ctx, &foundUserToken,
		"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM ldap_user_api_tokens WHERE token_key = $1",
		apiToken, l.config.UserAPITokenDuration().Duration(),
	)
	if err != nil {
		return sessions.User{}, err
	}
	if !foundUserToken.Valid { // API Token expired, purge
		if _, execErr := l.ds.ExecContext(ctx, "DELETE FROM ldap_user_api_tokens WHERE token_key = $1", apiToken); execErr != nil {
			l.lggr.Errorf("error purging stale ldap API token session: %v", execErr)
		}
		return sessions.User{}, sessions.ErrUserSessionExpired
	}

	return sessions.User{
		Email: foundUserToken.UserEmail,
		Role:  foundUserToken.UserRole,
	}, nil
}

// ListUsers will load and return all active users in applicable LDAP groups, extended with local admin users as well
func (l *ldapAuthenticator) ListUsers(ctx context.Context) ([]sessions.User, error) {
	// For each defined role/group, query for the list of group members to gather the full list of possible users
	users := []sessions.User{}
	var err error

	conn, err := l.ldapClient.CreateEphemeralConnection()
	if err != nil {
		l.lggr.Errorf("error in LDAP dial: ", err)
		return users, errors.New("unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	// Query for list of uniqueMember IDs present in Admin group
	adminUsers, err := l.ldapGroupMembersListToUser(conn, l.config.AdminUserGroupCN(), sessions.UserRoleAdmin)
	if err != nil {
		l.lggr.Errorf("error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("unable to list group users")
	}
	// Query for list of uniqueMember IDs present in Edit group
	editUsers, err := l.ldapGroupMembersListToUser(conn, l.config.EditUserGroupCN(), sessions.UserRoleEdit)
	if err != nil {
		l.lggr.Errorf("error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("unable to list group users")
	}
	// Query for list of uniqueMember IDs present in Run group
	runUsers, err := l.ldapGroupMembersListToUser(conn, l.config.RunUserGroupCN(), sessions.UserRoleRun)
	if err != nil {
		l.lggr.Errorf("error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("unable to list group users")
	}
	// Query for list of uniqueMember IDs present in Read group
	readUsers, err := l.ldapGroupMembersListToUser(conn, l.config.ReadUserGroupCN(), sessions.UserRoleView)
	if err != nil {
		l.lggr.Errorf("error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("unable to list group users")
	}

	// Aggregate full list
	users = append(users, adminUsers...)
	users = append(users, editUsers...)
	users = append(users, runUsers...)
	users = append(users, readUsers...)

	// Dedupe preserving order of highest role
	uniqueRef := make(map[string]struct{})
	dedupedUsers := []sessions.User{}
	for _, user := range users {
		if _, ok := uniqueRef[user.Email]; !ok {
			uniqueRef[user.Email] = struct{}{}
			dedupedUsers = append(dedupedUsers, user)
		}
	}

	// If no active attribute to check is defined, user simple being assigned the group is enough, return full list
	if l.config.ActiveAttribute() == "" {
		return dedupedUsers, nil
	}

	// Now optionally validate that all uniqueMembers are active in the org/LDAP server
	emails := []string{}
	for _, user := range dedupedUsers {
		emails = append(emails, user.Email)
	}
	activeUsers, err := l.validateUsersActive(emails)
	if err != nil {
		l.lggr.Errorf("error validating supplied user list: ", err)
		return users, errors.New("error validating supplied user list")
	}

	// Filter non active users
	returnUsers := []sessions.User{}
	for i, active := range activeUsers {
		if active {
			returnUsers = append(returnUsers, dedupedUsers[i])
		}
	}

	// Extend with local admin users
	var localAdminUsers []sessions.User
	sql := "SELECT * FROM users ORDER BY email ASC;"
	if err := l.ds.SelectContext(ctx, &localAdminUsers, sql); err != nil {
		l.lggr.Errorf("error extending upstream LDAP users with local admin users in users table: ", err)
	} else {
		returnUsers = append(returnUsers, localAdminUsers...)
	}

	return returnUsers, nil
}

// ldapGroupMembersListToUser queries the LDAP server given a conn for a list of uniqueMember who are part of the parameterized group
func (l *ldapAuthenticator) ldapGroupMembersListToUser(conn LDAPConn, groupNameCN string, roleToAssign sessions.UserRole) ([]sessions.User, error) {
	users, err := ldapGroupMembersListToUser(
		conn, groupNameCN, roleToAssign, l.config.GroupsDN(),
		l.config.BaseDN(), l.config.QueryTimeout(),
		l.lggr,
	)
	if err != nil {
		l.lggr.Errorf("error listing members of group (%s): %v", groupNameCN, err)
		return users, errors.New("error searching group members in LDAP directory")
	}
	return users, nil
}

// AuthorizedUserWithSession will return the API user associated with the Session ID if it
// exists and hasn't expired, and update session's LastUsed field. The state of the upstream LDAP server
// is polled and synced at the defined interval via a SleeperTask
func (l *ldapAuthenticator) AuthorizedUserWithSession(ctx context.Context, sessionID string) (sessions.User, error) {
	if len(sessionID) == 0 {
		return sessions.User{}, errors.New("session ID cannot be empty")
	}
	// Query the ldap_sessions table for given session ID, user role and email are cached so
	// no further upstream LDAP query is performed
	var foundSession struct {
		UserEmail string
		UserRole  sessions.UserRole
		Valid     bool
	}
	if err := l.ds.GetContext(ctx, &foundSession,
		"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM ldap_sessions WHERE id = $1",
		sessionID, l.config.SessionTimeout().Duration(),
	); err != nil {
		return sessions.User{}, sessions.ErrUserSessionExpired
	}
	if !foundSession.Valid {
		// Sessions expired, purge
		if _, execErr := l.ds.ExecContext(ctx, "DELETE FROM ldap_sessions WHERE id = $1", sessionID); execErr != nil {
			l.lggr.Errorf("error purging stale ldap session: %v", execErr)
		}
		return sessions.User{}, sessions.ErrUserSessionExpired
	}
	return sessions.User{
		Email: foundSession.UserEmail,
		Role:  foundSession.UserRole,
	}, nil
}

// DeleteUser is not supported for read only LDAP
func (l *ldapAuthenticator) DeleteUser(ctx context.Context, email string) error {
	return sessions.ErrNotSupported
}

// DeleteUserSession removes an ldapSession table entry by ID
func (l *ldapAuthenticator) DeleteUserSession(ctx context.Context, sessionID string) error {
	_, err := l.ds.ExecContext(ctx, "DELETE FROM ldap_sessions WHERE id = $1", sessionID)
	return err
}

// GetUserWebAuthn returns an empty stub, MFA token prompt is handled either by the upstream
// server blocking callback, or an error code to pass a OTP
func (l *ldapAuthenticator) GetUserWebAuthn(ctx context.Context, email string) ([]sessions.WebAuthn, error) {
	return []sessions.WebAuthn{}, nil
}

// CreateSession will forward the session request credentials to the
// LDAP server, querying for a user + role response if username and
// password match. The API call is blocking with timeout, so a sufficient timeout
// should allow the user to respond to potential MFA push notifications
func (l *ldapAuthenticator) CreateSession(ctx context.Context, sr sessions.SessionRequest) (string, error) {
	conn, err := l.ldapClient.CreateEphemeralConnection()
	if err != nil {
		return "", errors.New("unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	var returnErr error

	// Attempt to LDAP Bind with user provided credentials
	escapedEmail := ldap.EscapeFilter(strings.ToLower(sr.Email))
	searchBaseDN := fmt.Sprintf("%s=%s,%s,%s", l.config.BaseUserAttr(), escapedEmail, l.config.UsersDN(), l.config.BaseDN())
	if err = conn.Bind(searchBaseDN, sr.Password); err != nil {
		l.lggr.Infof("Error binding user authentication request in LDAP Bind: %v", err)
		returnErr = errors.New("unable to log in with LDAP server. Check credentials")
	}

	// Bind was successful meaning user and credentials are present in LDAP directory
	// Reuse FindUser functionality to fetch user roles used to create ldap_session entry
	// with cached user email and role
	foundUser, err := l.FindUser(ctx, escapedEmail)
	if err != nil {
		l.lggr.Infof("Successful user login, but error querying for user groups: user: %s, error %v", escapedEmail, err)
		returnErr = errors.New("log in successful, but no assigned groups to assume role")
	}

	isLocalUser := false
	if returnErr != nil {
		// Unable to log in against LDAP server, attempt fallback local auth with credentials, case of local CLI Admin account
		// Successful local user sessions can not be managed by the upstream server and have expiration handled by the reaper sync module
		foundUser, returnErr = l.localLoginFallback(ctx, sr)
		isLocalUser = true
	}

	// If err is still populated, return
	if returnErr != nil {
		return "", returnErr
	}

	l.lggr.Infof("Successful LDAP login request for user %s - %s", sr.Email, foundUser.Role)

	// Save session, user, and role to database. Given a session ID for future queries, the LDAP server will not be queried
	// Sessions are set to expire after the duration + creation date elapsed, and are synced on an interval against the upstream
	// LDAP server
	session := sessions.NewSession()
	_, err = l.ds.ExecContext(
		ctx,
		"INSERT INTO ldap_sessions (id, user_email, user_role, localauth_user, created_at) VALUES ($1, $2, $3, $4, now())",
		session.ID,
		strings.ToLower(sr.Email),
		foundUser.Role,
		isLocalUser,
	)
	if err != nil {
		l.lggr.Errorf("unable to create new session in ldap_sessions table %v", err)
		return "", fmt.Errorf("error creating local LDAP session: %w", err)
	}

	l.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]interface{}{"email": sr.Email})

	return session.ID, nil
}

// ClearNonCurrentSessions removes all ldap_sessions but the id passed in.
func (l *ldapAuthenticator) ClearNonCurrentSessions(ctx context.Context, sessionID string) error {
	_, err := l.ds.ExecContext(ctx, "DELETE FROM ldap_sessions where id != $1", sessionID)
	return err
}

// CreateUser is not supported for read only LDAP
func (l *ldapAuthenticator) CreateUser(ctx context.Context, user *sessions.User) error {
	return sessions.ErrNotSupported
}

// UpdateRole is not supported for read only LDAP
func (l *ldapAuthenticator) UpdateRole(ctx context.Context, email, newRole string) (sessions.User, error) {
	return sessions.User{}, sessions.ErrNotSupported
}

// SetPassword for remote users is not supported via the read only LDAP implementation, however change password
// in the context of updating a local admin user's password is required
func (l *ldapAuthenticator) SetPassword(ctx context.Context, user *sessions.User, newPassword string) error {
	// Ensure specified user is part of the local admins user table
	var localAdminUser sessions.User
	sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
	if err := l.ds.GetContext(ctx, &localAdminUser, sql, user.Email); err != nil {
		l.lggr.Infof("Can not change password, local user with email not found in users table: %s, err: %v", user.Email, err)
		return sessions.ErrNotSupported
	}

	// User is local admin, save new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	sql = "UPDATE users SET hashed_password = $1, updated_at = now() WHERE email = $2 RETURNING *"
	if err := l.ds.GetContext(ctx, user, sql, hashedPassword, user.Email); err != nil {
		l.lggr.Errorf("unable to set password for user: %s, err: %v", user.Email, err)
		return errors.New("unable to save password")
	}
	return nil
}

// TestPassword tests if an LDAP login bind can be performed with provided credentials, returns nil if success
func (l *ldapAuthenticator) TestPassword(ctx context.Context, email string, password string) error {
	conn, err := l.ldapClient.CreateEphemeralConnection()
	if err != nil {
		return errors.New("unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	// Attempt to LDAP Bind with user provided credentials
	escapedEmail := ldap.EscapeFilter(strings.ToLower(email))
	searchBaseDN := fmt.Sprintf("%s=%s,%s,%s", l.config.BaseUserAttr(), escapedEmail, l.config.UsersDN(), l.config.BaseDN())
	err = conn.Bind(searchBaseDN, password)
	if err == nil {
		return nil
	}
	l.lggr.Infof("Error binding user authentication request in TestPassword call LDAP Bind: %v", err)

	// Fall back to test local users table in case of supported local CLI users as well
	var hashedPassword string
	if err := l.ds.GetContext(ctx, &hashedPassword, "SELECT hashed_password FROM users WHERE lower(email) = lower($1)", email); err != nil {
		return errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, hashedPassword) {
		return errors.New("invalid credentials")
	}

	return nil
}

// CreateAndSetAuthToken generates a new credential token with the user role
func (l *ldapAuthenticator) CreateAndSetAuthToken(ctx context.Context, user *sessions.User) (*auth.Token, error) {
	newToken := auth.NewToken()

	err := l.SetAuthToken(ctx, user, newToken)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// SetAuthToken updates the user to use the given Authentication Token.
func (l *ldapAuthenticator) SetAuthToken(ctx context.Context, user *sessions.User, token *auth.Token) error {
	if !l.config.UserApiTokenEnabled() {
		return errors.New("API token is not enabled ")
	}

	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return fmt.Errorf("LDAPAuth SetAuthToken hashed secret error: %w", err)
	}

	err = sqlutil.TransactDataSource(ctx, l.ds, nil, func(tx sqlutil.DataSource) error {
		// Is this user a local CLI Admin or upstream LDAP user?
		// Check presence in local users table. Set localauth_user column true if present.
		// This flag omits the session/token from being purged by the sync daemon/reaper.go
		isLocalCLIAdmin := false
		err = l.ds.QueryRowxContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&isLocalCLIAdmin)
		if err != nil {
			return fmt.Errorf("error checking user presence in users table: %w", err)
		}

		// Remove any existing API tokens
		if _, err = l.ds.ExecContext(ctx, "DELETE FROM ldap_user_api_tokens WHERE user_email = $1", user.Email); err != nil {
			return fmt.Errorf("error executing DELETE FROM ldap_user_api_tokens: %w", err)
		}
		// Create new API token for user
		_, err = l.ds.ExecContext(
			ctx,
			"INSERT INTO ldap_user_api_tokens (user_email, user_role, localauth_user, token_key, token_salt, token_hashed_secret, created_at) VALUES ($1, $2, $3, $4, $5, $6, now())",
			user.Email,
			user.Role,
			isLocalCLIAdmin,
			token.AccessKey,
			salt,
			hashedSecret,
		)
		if err != nil {
			return fmt.Errorf("failed insert into ldap_user_api_tokens: %w", err)
		}
		return nil
	})
	if err != nil {
		return errors.New("error creating API token")
	}

	l.auditLogger.Audit(audit.APITokenCreated, map[string]interface{}{"user": user.Email})
	return nil
}

// DeleteAuthToken clears and disables the users Authentication Token.
func (l *ldapAuthenticator) DeleteAuthToken(ctx context.Context, user *sessions.User) error {
	_, err := l.ds.ExecContext(ctx, "DELETE FROM ldap_user_api_tokens WHERE email = $1")
	return err
}

// SaveWebAuthn is not supported for read only LDAP
func (l *ldapAuthenticator) SaveWebAuthn(ctx context.Context, token *sessions.WebAuthn) error {
	return sessions.ErrNotSupported
}

// Sessions returns all sessions limited by the parameters.
func (l *ldapAuthenticator) Sessions(ctx context.Context, offset, limit int) ([]sessions.Session, error) {
	var sessions []sessions.Session
	sql := `SELECT * FROM ldap_sessions ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err := l.ds.SelectContext(ctx, &sessions, sql, limit, offset); err != nil {
		return sessions, nil
	}
	return sessions, nil
}

// FindExternalInitiator supports the 'Run' role external intiator header auth functionality
func (l *ldapAuthenticator) FindExternalInitiator(ctx context.Context, eia *auth.Token) (*bridges.ExternalInitiator, error) {
	exi := &bridges.ExternalInitiator{}
	err := l.ds.GetContext(ctx, exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}

// localLoginFallback tests the credentials provided against the 'local' authentication method
// This covers the case of local CLI API calls requiring local login separate from the LDAP server
func (l *ldapAuthenticator) localLoginFallback(ctx context.Context, sr sessions.SessionRequest) (sessions.User, error) {
	var user sessions.User
	sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
	err := l.ds.GetContext(ctx, &user, sql, sr.Email)
	if err != nil {
		return user, err
	}
	if !constantTimeEmailCompare(strings.ToLower(sr.Email), strings.ToLower(user.Email)) {
		l.auditLogger.Audit(audit.AuthLoginFailedEmail, map[string]interface{}{"email": sr.Email})
		return user, errors.New("invalid email")
	}

	if !utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		l.auditLogger.Audit(audit.AuthLoginFailedPassword, map[string]interface{}{"email": sr.Email})
		return user, errors.New("invalid password")
	}

	return user, nil
}

// validateUsersActive performs an additional LDAP server query for the supplied emails, checking the
// returned user data for an 'active' property defined optionally in the config.
// Returns same length bool 'valid' array, indexed by sorted email
func (l *ldapAuthenticator) validateUsersActive(emails []string) ([]bool, error) {
	validUsers := make([]bool, len(emails))
	// If active attribute to check is not defined in config, skip
	if l.config.ActiveAttribute() == "" {
		// fill with valids
		for i := range emails {
			validUsers[i] = true
		}
		return validUsers, nil
	}

	conn, err := l.ldapClient.CreateEphemeralConnection()
	if err != nil {
		l.lggr.Errorf("error in LDAP dial: ", err)
		return validUsers, errors.New("unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	// Build the full email list query to pull all 'isActive' information for each user specified in one query
	filterQuery := "(|"
	for _, email := range emails {
		escapedEmail := ldap.EscapeFilter(email)
		filterQuery = fmt.Sprintf("%s(%s=%s)", filterQuery, l.config.BaseUserAttr(), escapedEmail)
	}
	filterQuery = fmt.Sprintf("(&%s))", filterQuery)
	searchBaseDN := fmt.Sprintf("%s,%s", l.config.UsersDN(), l.config.BaseDN())
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, int(l.config.QueryTimeout().Seconds()), false,
		filterQuery,
		[]string{l.config.BaseUserAttr(), l.config.ActiveAttribute()},
		nil,
	)
	// Query LDAP server for the ActiveAttribute property of each specified user
	results, err := conn.Search(searchRequest)
	if err != nil {
		l.lggr.Errorf("error searching user in LDAP query: %v", err)
		return validUsers, errors.New("error searching users in LDAP directory")
	}

	// Ensure user response entries
	if len(results.Entries) == 0 {
		return validUsers, ErrUserNotInUpstream
	}

	// Pull expected ActiveAttribute value from list of string possible values
	// keyed on email for final step to return flag bool list where order is preserved
	emailToActiveMap := make(map[string]bool)
	for _, result := range results.Entries {
		isActiveAttribute := result.GetAttributeValue(l.config.ActiveAttribute())
		uidAttribute := result.GetAttributeValue(l.config.BaseUserAttr())
		emailToActiveMap[uidAttribute] = isActiveAttribute == l.config.ActiveAttributeAllowedValue()
	}
	for i, email := range emails {
		active, ok := emailToActiveMap[email]
		if ok && active {
			validUsers[i] = true
		}
	}

	return validUsers, nil
}

// ldapGroupMembersListToUser queries the LDAP server given a conn for a list of uniqueMember who are part of the parameterized group. Reused by sync.go
func ldapGroupMembersListToUser(
	conn LDAPConn,
	groupNameCN string,
	roleToAssign sessions.UserRole,
	groupsDN string,
	baseDN string,
	queryTimeout time.Duration,
	lggr logger.Logger,
) ([]sessions.User, error) {
	users := []sessions.User{}
	// Prepare and query the GroupsDN for the specified group name
	searchBaseDN := fmt.Sprintf("%s, %s", groupsDN, baseDN)
	filterQuery := fmt.Sprintf("(&(cn=%s))", groupNameCN)
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, int(queryTimeout.Seconds()), false,
		filterQuery,
		[]string{UniqueMemberAttribute},
		nil,
	)
	result, err := conn.Search(searchRequest)
	if err != nil {
		lggr.Errorf("error searching group members in LDAP query: %v", err)
		return users, errors.New("error searching group members in LDAP directory")
	}

	// The result.Entry query response here is for the 'group' type of LDAP resource. The result should be a single entry, containing
	// a single Attribute named 'uniqueMember' containing a list of string Values. These Values are strings that should be returned in
	// the format "uid=test.user@example.com,ou=users,dc=example,dc=com". The 'uid' is then manually parsed here as the library does
	// not expose the functionality
	if len(result.Entries) != 1 {
		lggr.Errorf("unexpected length of query results for group user members, expected one got %d", len(result.Entries))
		return users, errors.New("error searching group members in LDAP directory")
	}

	// Get string list of members from 'uniqueMember' attribute
	uniqueMemberValues := result.Entries[0].GetAttributeValues(UniqueMemberAttribute)
	for _, uniqueMemberEntry := range uniqueMemberValues {
		parts := strings.Split(uniqueMemberEntry, ",") // Split attribute value on comma (uid, ou, dc parts)
		uidComponent := ""
		for _, part := range parts { // Iterate parts for "uid="
			if strings.HasPrefix(part, "uid=") {
				uidComponent = part
				break
			}
		}
		if uidComponent == "" {
			lggr.Errorf("unexpected LDAP group query response for unique members - expected list of LDAP Values for uniqueMember containing LDAP strings in format uid=test.user@example.com,ou=users,dc=example,dc=com. Got %s", uniqueMemberEntry)
			continue
		}
		// Map each user email to the sessions.User struct
		userEmail := strings.TrimPrefix(uidComponent, "uid=")
		users = append(users, sessions.User{
			Email: userEmail,
			Role:  roleToAssign,
		})
	}
	return users, nil
}

// groupSearchResultsToUserRole takes a list of LDAP group search result entries and returns the associated
// internal user role based on the group name mappings defined in the configuration
func (l *ldapAuthenticator) groupSearchResultsToUserRole(ldapGroups []*ldap.Entry) (sessions.UserRole, error) {
	return GroupSearchResultsToUserRole(
		ldapGroups,
		l.config.AdminUserGroupCN(),
		l.config.EditUserGroupCN(),
		l.config.RunUserGroupCN(),
		l.config.ReadUserGroupCN(),
	)
}

func GroupSearchResultsToUserRole(ldapGroups []*ldap.Entry, adminCN string, editCN string, runCN string, readCN string) (sessions.UserRole, error) {
	// If defined Admin group name is present in groups search result, return UserRoleAdmin
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == adminCN {
			return sessions.UserRoleAdmin, nil
		}
	}
	// Check edit role
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == editCN {
			return sessions.UserRoleEdit, nil
		}
	}
	// Check run role
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == runCN {
			return sessions.UserRoleRun, nil
		}
	}
	// Check view role
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == readCN {
			return sessions.UserRoleView, nil
		}
	}
	// No role group found, error
	return sessions.UserRoleView, ErrUserNoLDAPGroups
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
