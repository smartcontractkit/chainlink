package ldapauth

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/sqlx"
)

type LDAPServerStateSyncer struct {
	q      pg.Q
	config config.LDAP
	lggr   logger.Logger
}

// NewLDAPServerStateSync creates a reaper that cleans stale sessions from the store.
func NewLDAPServerStateSync(
	db *sqlx.DB,
	pgCfg pg.QConfig,
	config config.LDAP,
	lggr logger.Logger,
) utils.SleeperTask {
	namedLogger := lggr.Named("LDAPServerStateSync")
	return utils.NewSleeperTask(&LDAPServerStateSyncer{
		pg.NewQ(db, namedLogger, pgCfg),
		config,
		namedLogger,
	})
}

func (ldSync *LDAPServerStateSyncer) Name() string {
	return "LDAPServerStateSync"
}

func (ldSync *LDAPServerStateSyncer) Work() {
	// Purge expired ldap_sessions and ldap_user_api_tokens
	recordCreationStaleThreshold := ldSync.config.UpstreamSyncInterval().Before(ldSync.config.SessionTimeout().Before(time.Now()))
	err := ldSync.deleteStaleSessions(recordCreationStaleThreshold)
	if err != nil {
		ldSync.lggr.Error("unable to expire local LDAP sessions: ", err)
	}
	recordCreationStaleThreshold = ldSync.config.UserAPITokenDuration().Before(ldSync.config.SessionTimeout().Before(time.Now()))
	err = ldSync.deleteStaleAPITokens(recordCreationStaleThreshold)
	if err != nil {
		ldSync.lggr.Error("unable to expire user API tokens: ", err)
	}

	// For each defined role/group, query for the list of group members to gather the full list of possible users
	users := []sessions.User{}

	// Establish ephemeral connection
	conn, err := ldap.DialURL(ldSync.config.ServerAddress())
	if err != nil {
		ldSync.lggr.Errorf("Failed to Dial LDAP Server", err)
	}
	// Root level root user auth with credentials provided from config
	bindStr := ldSync.config.BaseUserAttr() + "=" + ldSync.config.ReadOnlyUserLogin() + "," + ldSync.config.BaseDN()
	if err := conn.Bind(bindStr, ldSync.config.ReadOnlyUserPass()); err != nil {
		ldSync.lggr.Errorf("Unable to login as initial root LDAP user", err)
	}
	defer conn.Close()

	// Query for list of uniqueMember IDs present in Admin group
	adminUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.AdminUserGroupCN(), sessions.UserRoleAdmin)
	if err != nil {
		ldSync.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
	}
	// Query for list of uniqueMember IDs present in Edit group
	editUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.EditUserGroupCN(), sessions.UserRoleEdit)
	if err != nil {
		ldSync.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
	}
	// Query for list of uniqueMember IDs present in Edit group
	runUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.RunUserGroupCN(), sessions.UserRoleRun)
	if err != nil {
		ldSync.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
	}
	// Query for list of uniqueMember IDs present in Edit group
	readUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.ReadUserGroupCN(), sessions.UserRoleView)
	if err != nil {
		ldSync.lggr.Errorf("Error in ldapGroupMembersListToUser: ", err)
	}

	users = append(users, adminUsers...)
	users = append(users, editUsers...)
	users = append(users, runUsers...)
	users = append(users, readUsers...)

	// Dedupe preserving order of highest role (sorted)
	// Preserve members as a map for future lookup
	upstreamUserStateMap := make(map[string]sessions.User)
	for _, user := range users {
		if _, ok := upstreamUserStateMap[user.Email]; !ok {
			upstreamUserStateMap[user.Email] = user
		}
	}

	// upstreamUserStateMap is now the most up to date source of truth
	// Now sync database sessions and roles with new data
	err = ldSync.q.Transaction(func(tx pg.Queryer) error {
		// First, purge users present in the local ldap_sessions table but not in the upstream server
		type LDAPSession struct {
			UserEmail string
			UserRole  sessions.UserRole
		}
		var existingSessions []LDAPSession
		if err := tx.Select(&existingSessions, "SELECT user_email, user_role FROM ldap_sessions WHERE localauth_user = false"); err != nil {
			return errors.Wrap(err, "Unable to query ldap_sessions table")
		}
		var existingAPITokens []LDAPSession
		if err := tx.Select(&existingAPITokens, "SELECT user_email, user_role FROM ldap_user_api_tokens WHERE localauth_user = false"); err != nil {
			return errors.Wrap(err, "Unable to query ldap_user_api_tokens table")
		}

		// Create existing sessions and API tokens lookup map for later
		existingSessionsMap := make(map[string]LDAPSession)
		for _, sess := range existingSessions {
			existingSessionsMap[sess.UserEmail] = sess
		}
		existingAPITokensMap := make(map[string]LDAPSession)
		for _, sess := range existingAPITokens {
			existingAPITokensMap[sess.UserEmail] = sess
		}

		// Populate list of session emails present in the local session table but not in the upstream state
		emailsToPurge := []interface{}{}
		for _, ldapSession := range existingSessions {
			if _, ok := upstreamUserStateMap[ldapSession.UserEmail]; !ok {
				emailsToPurge = append(emailsToPurge, ldapSession.UserEmail)
			}
		}
		// Likewise for API Tokens table
		apiTokenEmailsToPurge := []interface{}{}
		for _, ldapSession := range existingAPITokens {
			if _, ok := upstreamUserStateMap[ldapSession.UserEmail]; !ok {
				apiTokenEmailsToPurge = append(apiTokenEmailsToPurge, ldapSession.UserEmail)
			}
		}

		// Remove any active sessions this user may have
		if len(emailsToPurge) > 0 {
			placeholders := make([]string, len(emailsToPurge))
			for i := range emailsToPurge {
				placeholders[i] = fmt.Sprintf("$%d", i+1)
			}
			query := fmt.Sprintf("DELETE FROM ldap_sessions WHERE user_email IN (%s)", strings.Join(placeholders, ", "))
			_, err := ldSync.q.Exec(query, emailsToPurge...)
			if err != nil {
				return err
			}
		}

		// Remove any active API tokens this user may have
		if len(apiTokenEmailsToPurge) > 0 {
			placeholders := make([]string, len(apiTokenEmailsToPurge))
			for i := range apiTokenEmailsToPurge {
				placeholders[i] = fmt.Sprintf("$%d", i+1)
			}
			query := fmt.Sprintf("DELETE FROM ldap_user_api_tokens WHERE user_email IN (%s)", strings.Join(placeholders, ", "))
			_, err = ldSync.q.Exec(query, apiTokenEmailsToPurge...)
			if err != nil {
				return err
			}
		}

		// For each user session row, update role to match state of user map from upstream source
		queryWhenClause := ""
		emailValues := []interface{}{}
		// Prepare CASE WHEN query statement with parameterized argument $n placeholders and matching role based on index
		for email, user := range upstreamUserStateMap {
			// Only build on SET CASE statement per local session and API token role, not for each upstream user value
			_, sessionOk := existingSessionsMap[email]
			_, tokenOk := existingAPITokensMap[email]
			if !sessionOk && !tokenOk {
				continue
			}
			emailValues = append(emailValues, email)
			queryWhenClause += fmt.Sprintf("WHEN user_email = $%d THEN '%s' ", len(emailValues), user.Role)
		}

		// Set new role state for all rows in single Exec
		query := fmt.Sprintf("UPDATE ldap_sessions SET user_role = CASE %s ELSE user_role END", queryWhenClause)
		_, err := ldSync.q.Exec(query, emailValues...)
		if err != nil {
			return err
		}

		// Update role of API tokens as well
		query = fmt.Sprintf("UPDATE ldap_user_api_tokens SET user_role = CASE %s ELSE user_role END", queryWhenClause)
		_, err = ldSync.q.Exec(query, emailValues...)
		if err != nil {
			return err
		}

		ldSync.lggr.Info("local ldap_sessions and ldap_user_api_tokens table successfully synced with upstream LDAP state")
		return nil
	})
	if err != nil {
		ldSync.lggr.Errorf("Error syncing local database state: ", err)
	}
}

// deleteStaleSessions deletes all ldap_sessions before the passed time.
func (ldSync *LDAPServerStateSyncer) deleteStaleSessions(before time.Time) error {
	_, err := ldSync.q.Exec("DELETE FROM ldap_sessions WHERE created_at < $1", before)
	return err
}

// deleteStaleAPITokens deletes all ldap_user_api_tokens before the passed time.
func (ldSync *LDAPServerStateSyncer) deleteStaleAPITokens(before time.Time) error {
	_, err := ldSync.q.Exec("DELETE FROM ldap_user_api_tokens WHERE created_at < $1", before)
	return err
}

// ldapGroupMembersListToUser queries the LDAP server given a conn for a list of uniqueMember who are part of the parameterized group
func (l *LDAPServerStateSyncer) ldapGroupMembersListToUser(conn *ldap.Conn, groupNameCN string, roleToAssign sessions.UserRole) ([]sessions.User, error) {
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

	// The result.Entry query response here is for the 'group' type of LDAP resource. The result should be a single entry, containing
	// a single Attribute named 'uniqueMember' containing a list of string Values. These Values are strings that should be returned in
	// the format "uid=test.user@example.com,ou=users,dc=example,dc=com". The 'uid' is then manually parsed here as the library does
	// not expose the functionality
	if len(result.Entries) != 1 {
		l.lggr.Errorf("Unexpected length of query results for group user members, expected one got %d", len(result.Entries))
		return users, errors.New("Error searching group members in LDAP directory")
	}

	// Get string list of members from 'uniqueMember' attribute
	uniqueMemberValues := result.Entries[0].GetAttributeValues(LDAPUniqueMemberAttribute)
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
			l.lggr.Errorf("Unexpected LDAP group query response for unique members - expected list of LDAP Values for uniqueMember containing LDAP strings in format uid=test.user@example.com,ou=users,dc=example,dc=com. Got %s", uniqueMemberEntry)
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
