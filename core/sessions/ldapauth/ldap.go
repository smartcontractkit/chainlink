package ldapauth

/*

The LDAP authentication module forwards the credentials in the user session request
for authentication with a configured upstream LDAP server

This module relies on the two following local database tables:
	ldap_sessions: 	Upon successful LDAP response, creates a keyed local copy of the user
					email
	ldap_user_api_tokens: User created API tokens, tied to the node, storing user email.
						  Note: user can have only one API token at a time, and token
						  expiration is enforced

User session and roles are cached and revalidated with the upstream service at the interval defined in
the local LDAP config through the Application.sessionReaper implementation in reaper.go

This implementation is read only; user mutation actions such as Delete are not supported.

MFA is supported via the remote LDAP server implementation. Sufficient request time out should accomodate
for a blocking auth call while the user responds to a potential push notification callback.

*/

import (
	"crypto/subtle"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/mathutil"
)

const (
	LDAPUniqueMemberAttribute = "uniqueMember"
)

// implements sessions.AuthenticationProvider interface
type ldapAuthenticator struct {
	q           pg.Q
	config      config.LDAP
	lggr        logger.Logger
	auditLogger audit.AuditLogger
}

var _ sessions.AuthenticationProvider = (*ldapAuthenticator)(nil)

func NewLDAPAuthenticator(
	db *sqlx.DB,
	pgCfg pg.QConfig,
	ldapCfg config.LDAP,
	dev bool,
	lggr logger.Logger,
	auditLogger audit.AuditLogger,
) (sessions.AuthenticationProvider, error) {
	namedLogger := lggr.Named("LDAPAuthenticationProvider")

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
		q:           pg.NewQ(db, namedLogger, pgCfg),
		config:      ldapCfg,
		lggr:        lggr.Named("LDAPAuthenticationProvider"),
		auditLogger: auditLogger,
	}

	// Single override of library defined global
	ldap.DefaultTimeout = ldapCfg.QueryTimeout()

	// Test initial connection and credentials
	lggr.Infof("Attempting initial connection to configured LDAP server with bind as API user")
	conn, err := ldapAuth.dialAndConnect()
	if err != nil {
		return nil, errors.Errorf("Unable to establish connection to LDAP server with provided URL and credentials: %v", err)
	}
	conn.Close()

	// Store LDAP connection config for auth/new connection per request instead of persisted connection with reconnect
	return &ldapAuth, nil
}

// FindUser will attempt to return an LDAP user with mapped role by email.
func (l *ldapAuthenticator) FindUser(email string) (sessions.User, error) {
	email = strings.ToLower(email)
	foundUser := sessions.User{}

	// First query for user "is active" property if defined
	usersActive, err := l.validateUsersActive([]string{email})
	if err != nil {
		return foundUser, errors.New("Error running query to validate user")
	}
	if usersActive[0] == false {
		return foundUser, errors.New("User not active")
	}

	// Establish ephemeral connection
	conn, err := l.dialAndConnect()
	if err != nil {
		l.lggr.Errorf("Error in LDAP dial: ", err)
		return foundUser, errors.New("Unable to establish connection to LDAP server with provided URL and credentials")
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
		l.lggr.Errorf("Error searching users in LDAP query: %v", err)
		return foundUser, errors.New("Error searching users in LDAP directory")
	}

	if len(result.Entries) == 0 {
		return foundUser, errors.New("No users found with provided email")
	}

	// Populate found user by email and role based on matched group names
	userRole, err := l.groupSearchResultsToUserRole(result.Entries)
	if err != nil {
		l.lggr.Warnf("User '%s' found but no matching assigned groups in LDAP to assume role", email)
		return sessions.User{}, err
	}

	// Convert search result to sessions.User type with required fields
	foundUser = sessions.User{
		Email: email,
		Role:  userRole,
	}

	return foundUser, nil
}

// FindUserByAPIToken retrieves a possible stored user and role from the ldap_user_api_tokens table store
func (l *ldapAuthenticator) FindUserByAPIToken(apiToken string) (sessions.User, error) {
	if !l.config.UserApiTokenEnabled() {
		return sessions.User{}, errors.New("API token is not enabled ")
	}

	var foundUser sessions.User
	err := l.q.Transaction(func(tx pg.Queryer) error {
		// Query the ldap user API token table for given token, user role and email are cached so
		// no further upstream LDAP query is performed, sessions and tokens are synced against the upstream server
		// via the UpstreamSyncInterval config and reaper.go sync implementation
		var foundUserToken struct {
			UserEmail string
			UserRole  sessions.UserRole
			Valid     bool
		}
		if err := tx.Get(&foundUserToken,
			"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM ldap_user_api_tokens WHERE token_key = $1",
			apiToken, l.config.UserAPITokenDuration().Duration(),
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
		if err == sessions.ErrUserSessionExpired {
			// API Token expired, purge
			if _, err := l.q.Exec("DELETE FROM ldap_user_api_tokens WHERE token_key = $1", apiToken); err != nil {
				l.lggr.Errorf("Error purging stale ldap API token session: %v", err)
			}
		}
		return sessions.User{}, err
	}
	return foundUser, nil
}

// ListUsers will load and return all user rows from the db.
func (l *ldapAuthenticator) ListUsers() ([]sessions.User, error) {
	// For each defined role/group, query for the list of group members to gather the full list of possible users
	users := []sessions.User{}
	var err error

	// Establish ephemeral connection
	conn, err := l.dialAndConnect()
	if err != nil {
		l.lggr.Errorf("Error in LDAP dial: ", err)
		return users, errors.New("Unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	// Query for list of uniqueMember IDs present in Admin group
	adminUsers, err := l.ldapGroupMembersListToUser(conn, l.config.AdminUserGroupCN(), sessions.UserRoleAdmin)
	if err != nil {
		l.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("Unable to list group users")
	}
	// Query for list of uniqueMember IDs present in Edit group
	editUsers, err := l.ldapGroupMembersListToUser(conn, l.config.EditUserGroupCN(), sessions.UserRoleEdit)
	if err != nil {
		l.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("Unable to list group users")
	}
	// Query for list of uniqueMember IDs present in Edit group
	runUsers, err := l.ldapGroupMembersListToUser(conn, l.config.RunUserGroupCN(), sessions.UserRoleRun)
	if err != nil {
		l.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("Unable to list group users")
	}
	// Query for list of uniqueMember IDs present in Edit group
	readUsers, err := l.ldapGroupMembersListToUser(conn, l.config.ReadUserGroupCN(), sessions.UserRoleView)
	if err != nil {
		l.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
		return users, errors.New("Unable to list group users")
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
		l.lggr.Errorf("Error validating supplied user list: ", err)
		return users, errors.New("Error validating supplied user list")
	}

	// Filter non active users
	returnUsers := []sessions.User{}
	for i, active := range activeUsers {
		if active {
			returnUsers = append(returnUsers, dedupedUsers[i])
		}
	}

	return returnUsers, nil
}

// ldapGroupMembersListToUser queries the LDAP server given a conn for a list of uniqueMember who are part of the parameterized group
func (l *ldapAuthenticator) ldapGroupMembersListToUser(conn *ldap.Conn, groupNameCN string, roleToAssign sessions.UserRole) ([]sessions.User, error) {
	users := []sessions.User{}
	// Prepare and query the GroupsDN for the specified group name
	searchBaseDN := fmt.Sprintf("%s, %s", l.config.GroupsDN(), l.config.BaseDN())
	filterQuery := fmt.Sprintf("(&(cn=%s))", groupNameCN)
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, int(l.config.QueryTimeout().Seconds()), false,
		filterQuery,
		[]string{LDAPUniqueMemberAttribute},
		nil,
	)
	result, err := conn.Search(searchRequest)
	if err != nil {
		l.lggr.Errorf("Error searching group members in LDAP query: %v", err)
		return users, errors.New("Error searching group members in LDAP directory")
	}
	// Resulting entries are unique members for the group, map each user to the sessions.User struct
	for _, user := range result.Entries {
		users = append(users, sessions.User{
			Email: user.GetAttributeValue(l.config.BaseUserAttr()),
			Role:  roleToAssign,
		})
	}
	return users, nil
}

// AuthorizedUserWithSession will return the API user associated with the Session ID if it
// exists and hasn't expired, and update session's LastUsed field. The state of the upstream LDAP server
// is polled and synced at the defined interval via a SleeperTask
func (l *ldapAuthenticator) AuthorizedUserWithSession(sessionID string) (sessions.User, error) {
	if len(sessionID) == 0 {
		return sessions.User{}, errors.New("Session ID cannot be empty")
	}
	var foundUser sessions.User
	err := l.q.Transaction(func(tx pg.Queryer) error {
		// Query the ldap_sessions table for given session ID, user role and email are cached so
		// no further upstream LDAP query is performed
		var foundSession struct {
			UserEmail string
			UserRole  sessions.UserRole
			Valid     bool
		}
		if err := tx.Get(&foundSession,
			"SELECT user_email, user_role, created_at + $2 >= now() as valid FROM ldap_sessions WHERE id = $1",
			sessionID, l.config.SessionTimeout().Duration(),
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
		if err == sessions.ErrUserSessionExpired {
			if _, err := l.q.Exec("DELETE FROM ldap_sessions WHERE id = $1", sessionID); err != nil {
				l.lggr.Errorf("Error purging stale ldap session: %v", err)
			}
		}
		return sessions.User{}, err
	}
	return foundUser, nil
}

// DeleteUser is not supported for read only LDAP
func (l *ldapAuthenticator) DeleteUser(email string) error {
	return sessions.ErrNotSupported
}

// DeleteUserSession removes an ldapSession table entry by ID
func (l *ldapAuthenticator) DeleteUserSession(sessionID string) error {
	_, err := l.q.Exec("DELETE FROM ldap_sessions WHERE id = $1", sessionID)
	return err
}

// GetUserWebAuthn returns an empty stub, MFA token prompt is handled either by the upstream
// server blocking callback, or an error code to pass a OTP
func (l *ldapAuthenticator) GetUserWebAuthn(email string) ([]sessions.WebAuthn, error) {
	return []sessions.WebAuthn{}, nil
}

// CreateSession will forward the session request credentials to the
// LDAP server, querying for a user + role response if username and
// password match. The API call is blocking with timeout, so a sufficient timeout
// should allow the user to respond to potential MFA push notifications
func (l *ldapAuthenticator) CreateSession(sr sessions.SessionRequest) (string, error) {
	// Establish ephemeral connection
	conn, err := l.dialAndConnect()
	if err != nil {
		return "", errors.New("Unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	var returnErr error

	// Attempt to LDAP Bind with user provided credentials
	escapedEmail := ldap.EscapeFilter(strings.ToLower(sr.Email))
	searchBaseDN := fmt.Sprintf("%s=%s,%s,%s", l.config.BaseUserAttr(), escapedEmail, l.config.UsersDN(), l.config.BaseDN())
	if err := conn.Bind(searchBaseDN, sr.Password); err != nil {
		l.lggr.Infof("Error binding user authentication request in LDAP Bind: %v", err)
		returnErr = errors.New("Unable to log in with LDAP server. Check credentials")
	}

	// Bind was successful meaning user and credentials are present in LDAP directory
	// Reuse FindUser functionality to fetch user roles used to create ldap_session entry
	// with cached user email and role
	foundUser, err := l.FindUser(escapedEmail)
	if err != nil {
		l.lggr.Infof("Successful user login, but error querying for user groups: user: %s, error %v", escapedEmail, err)
		returnErr = errors.New("Log in successful, but no assigned groups to assume role")
	}

	if returnErr != nil {
		// Unable to log in against LDAP server, attempt fallback local auth with credentials, case of local CLI Admin account
		foundUser, returnErr = l.localLoginFallback(sr)
	}

	// If err is still populate, return
	if returnErr != nil {
		return "", returnErr
	}

	l.lggr.Infof("Successful LDAP login request for user %s - %s", sr.Email, foundUser.Role)

	// Save session, user, and role to database. Given a session ID for future queries, the LDAP server will not be queried
	// Sessions are set to expire after the duration + creation date elapsed, and are synced on an interval against the upstream
	// LDAP server
	session := sessions.NewSession()
	_, err = l.q.Exec(
		"INSERT INTO ldap_sessions (id, user_email, user_role, created_at) VALUES ($1, $2, $3, now())",
		session.ID,
		strings.ToLower(sr.Email),
		foundUser.Role,
	)
	if err != nil {
		return "", errors.Wrap(err, "unable to create new session in ldap_sessions table")
	}

	l.auditLogger.Audit(audit.AuthLoginSuccessNo2FA, map[string]interface{}{"email": sr.Email})

	return session.ID, nil
}

// ClearNonCurrentSessions removes all ldap_sessions but the id passed in.
func (l *ldapAuthenticator) ClearNonCurrentSessions(sessionID string) error {
	_, err := l.q.Exec("DELETE FROM ldap_sessions where id != $1", sessionID)
	return err
}

// CreateUser is not supported for read only LDAP
func (l *ldapAuthenticator) CreateUser(user *sessions.User) error {
	return sessions.ErrNotSupported
}

// UpdateRole is not supported for read only LDAP
func (l *ldapAuthenticator) UpdateRole(email, newRole string) (sessions.User, error) {
	return sessions.User{}, sessions.ErrNotSupported
}

// SetPassword is not supported for read only LDAP
func (l *ldapAuthenticator) SetPassword(user *sessions.User, newPassword string) error {
	return sessions.ErrNotSupported
}

// TestPassword tests if an LDAP login bind can be performed with provided credentials, returns nil if success
func (l *ldapAuthenticator) TestPassword(email string, password string) error {
	// Establish ephemeral connection
	conn, err := l.dialAndConnect()
	if err != nil {
		return errors.New("Unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()
	// Attempt to LDAP Bind with user provided credentials
	escapedEmail := ldap.EscapeFilter(strings.ToLower(email))
	searchBaseDN := fmt.Sprintf("%s=%s,%s,%s", l.config.BaseUserAttr(), escapedEmail, l.config.UsersDN(), l.config.BaseDN())
	if err := conn.Bind(searchBaseDN, password); err != nil {
		l.lggr.Infof("Error binding user authentication request in TestPassword call LDAP Bind: %v", err)
		return errors.New("Invalid credentials")
	}
	return nil
}

// CreateAndSetAuthToken generates a new credential token with the user role
func (l *ldapAuthenticator) CreateAndSetAuthToken(user *sessions.User) (*auth.Token, error) {
	newToken := auth.NewToken()

	err := l.SetAuthToken(user, newToken)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// SetAuthToken updates the user to use the given Authentication Token.
func (l *ldapAuthenticator) SetAuthToken(user *sessions.User, token *auth.Token) error {
	if !l.config.UserApiTokenEnabled() {
		return errors.New("API token is not enabled ")
	}

	salt := utils.NewSecret(utils.DefaultSecretSize)
	hashedSecret, err := auth.HashedSecret(token, salt)
	if err != nil {
		return errors.Wrap(err, "LDAPAuth SetAuthToken hashed secret error")
	}

	err = l.q.Transaction(func(tx pg.Queryer) error {
		// Remove any existing API tokens
		if _, err := l.q.Exec("DELETE FROM ldap_user_api_tokens WHERE user_email = $1", user.Email); err != nil {
			return errors.Wrap(err, "Error executing DELETE FROM ldap_user_api_tokens")
		}
		// Create new API token for user
		_, err = l.q.Exec(
			"INSERT INTO ldap_user_api_tokens (user_email, user_role, token_key, token_salt, token_hashed_secret, created_at) VALUES ($1, $2, $3, $4, $5, now())",
			user.Email,
			user.Role,
			token.AccessKey,
			salt,
			hashedSecret,
		)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error removing potential existing token and new token creation")
	}

	l.auditLogger.Audit(audit.APITokenCreated, map[string]interface{}{"user": user.Email})
	return nil
}

// DeleteAuthToken clears and disables the users Authentication Token.
func (l *ldapAuthenticator) DeleteAuthToken(user *sessions.User) error {
	_, err := l.q.Exec("DELETE FROM ldap_user_api_tokens WHERE email = $1")
	return err
}

// SaveWebAuthn is not supported for read only LDAP
func (l *ldapAuthenticator) SaveWebAuthn(token *sessions.WebAuthn) error {
	return sessions.ErrNotSupported
}

// Sessions returns all sessions limited by the parameters.
func (l *ldapAuthenticator) Sessions(offset, limit int) ([]sessions.Session, error) {
	var sessions []sessions.Session
	sql := `SELECT * FROM ldap_sessions ORDER BY created_at, id LIMIT $1 OFFSET $2;`
	if err := l.q.Select(&sessions, sql, limit, offset); err != nil {
		return sessions, nil
	}
	return sessions, nil
}

// FindExternalInitiator supports the 'Run' role external intiator header auth functionality
func (l *ldapAuthenticator) FindExternalInitiator(eia *auth.Token) (*bridges.ExternalInitiator, error) {
	exi := &bridges.ExternalInitiator{}
	err := l.q.Get(exi, `SELECT * FROM external_initiators WHERE access_key = $1`, eia.AccessKey)
	return exi, err
}

// localLoginFallback tests the credentials provided against the 'local' authentication method
// This covers the case of local CLI API calls requiring local login separate from the LDAP server
func (l *ldapAuthenticator) localLoginFallback(sr sessions.SessionRequest) (sessions.User, error) {
	var user sessions.User
	sql := "SELECT * FROM users WHERE lower(email) = lower($1)"
	err := l.q.Get(&user, sql, sr.Email)
	if err != nil {
		return user, err
	}
	if !constantTimeEmailCompare(strings.ToLower(sr.Email), strings.ToLower(user.Email)) {
		l.auditLogger.Audit(audit.AuthLoginFailedEmail, map[string]interface{}{"email": sr.Email})
		return user, errors.New("Invalid email")
	}

	if !utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		l.auditLogger.Audit(audit.AuthLoginFailedPassword, map[string]interface{}{"email": sr.Email})
		return user, errors.New("Invalid password")
	}

	return user, nil
}

// dialAndConnect returns a valid, active LDAP connection for querying
func (l *ldapAuthenticator) dialAndConnect() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(l.config.ServerAddress())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Dial LDAP Server")
	}
	// Root level root user auth with credentials provided from config
	bindStr := l.config.BaseUserAttr() + "=" + l.config.ReadOnlyUserLogin() + "," + l.config.BaseDN()
	if err := conn.Bind(bindStr, l.config.ReadOnlyUserPass()); err != nil {
		return nil, errors.Wrap(err, "Unable to login as initial root LDAP user")
	}
	return conn, nil
}

// validateUsersActive performs an additional LDAP server query for the supplied email, checking the
// return user data for an 'active' property defined optionally in the config. Returns same length bool array
func (l *ldapAuthenticator) validateUsersActive(emails []string) ([]bool, error) {
	validUsers := make([]bool, len(emails))
	// If active attribute to check is not defined in config, skip
	if l.config.ActiveAttribute() != "" {
		// fill with valids
		for i, _ := range emails {
			validUsers[i] = true
		}
		return validUsers, nil
	}

	// Establish ephemeral connection
	conn, err := l.dialAndConnect()
	if err != nil {
		l.lggr.Errorf("Error in LDAP dial: ", err)
		return validUsers, errors.New("Unable to establish connection to LDAP server with provided URL and credentials")
	}
	defer conn.Close()

	// Build the full or "|" query to pull all information for each user specified in one query
	orQuery := ""
	for _, email := range emails {
		escapedEmail := ldap.EscapeFilter(email)
		orQuery += fmt.Sprintf("(uniquemember=%s=%s,%s,%s)", l.config.BaseUserAttr(), escapedEmail, l.config.UsersDN(), l.config.BaseDN())
	}
	searchBaseDN := fmt.Sprintf("%s, %s", l.config.GroupsDN(), l.config.BaseDN())
	filterQuery := fmt.Sprintf("(|%s)", orQuery)
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, int(l.config.QueryTimeout().Seconds()), false,
		filterQuery,
		[]string{l.config.ActiveAttribute()},
		nil,
	)
	// Query LDAP server for the ActiveAttribute property of each specified user
	results, err := conn.Search(searchRequest)
	if err != nil {
		l.lggr.Errorf("Error searching user in LDAP query: %v", err)
		return validUsers, errors.New("Error searching users in LDAP directory")
	}

	// Ensure user response entries
	if len(results.Entries) == 0 {
		return validUsers, errors.New("No users matching email query")
	}

	// Pull expected ActiveAttribute value from list of string possible values
	// and set return bool array with Active truthiness value
	for i, result := range results.Entries {
		attributeValues := result.GetAttributeValue(l.config.ActiveAttribute())
		if attributeValues == l.config.ActiveAttributeAllowedValue() {
			validUsers[i] = true
		}
	}
	return validUsers, nil
}

// groupSearchResultsToUserRole takes a list of LDAP group search result entries and returns the associated
// internal user role based on the group name mappings defined in the configuration
func (l *ldapAuthenticator) groupSearchResultsToUserRole(ldapGroups []*ldap.Entry) (sessions.UserRole, error) {
	// If defined Admin group name is present in groups search result, return UserRoleAdmin
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == l.config.AdminUserGroupCN() {
			return sessions.UserRoleAdmin, nil
		}
	}
	// Check edit role
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == l.config.EditUserGroupCN() {
			return sessions.UserRoleEdit, nil
		}
	}
	// Check run role
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == l.config.RunUserGroupCN() {
			return sessions.UserRoleRun, nil
		}
	}
	// Check view role
	for _, group := range ldapGroups {
		if group.GetAttributeValue("cn") == l.config.ReadUserGroupCN() {
			return sessions.UserRoleView, nil
		}
	}
	// No role group found, error
	return sessions.UserRoleView, errors.New("User present in directory, but matching no role groups assigned")
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
