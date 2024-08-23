package ldapauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
)

type LDAPServerStateSyncer struct {
	ds           sqlutil.DataSource
	ldapClient   LDAPClient
	config       config.LDAP
	lggr         logger.Logger
	nextSyncTime time.Time
}

// NewLDAPServerStateSync creates a reaper that cleans stale sessions from the store.
func NewLDAPServerStateSync(
	ds sqlutil.DataSource,
	config config.LDAP,
	lggr logger.Logger,
) *utils.SleeperTask {
	namedLogger := lggr.Named("LDAPServerStateSync")
	serverSync := LDAPServerStateSyncer{
		ds:           ds,
		ldapClient:   newLDAPClient(config),
		config:       config,
		lggr:         namedLogger,
		nextSyncTime: time.Time{},
	}
	// If enabled, start a background task that calls the Sync/Work function on an
	// interval without needing an auth event to trigger it
	// Use IsInstant to check 0 value to omit functionality.
	if !config.UpstreamSyncInterval().IsInstant() {
		lggr.Info("LDAP Config UpstreamSyncInterval is non-zero, sync functionality will be called on a timer, respecting the UpstreamSyncRateLimit value")
		serverSync.StartWorkOnTimer()
	} else {
		// Ensure upstream server state is synced on startup manually if interval check not set
		serverSync.Work()
	}

	// Start background Sync call task reactive to auth related events
	serverSyncSleeperTask := utils.NewSleeperTask(&serverSync)
	return serverSyncSleeperTask
}

func (ldSync *LDAPServerStateSyncer) Name() string {
	return "LDAPServerStateSync"
}

func (ldSync *LDAPServerStateSyncer) StartWorkOnTimer() {
	time.AfterFunc(ldSync.config.UpstreamSyncInterval().Duration(), ldSync.StartWorkOnTimer)
	ldSync.Work()
}

func (ldSync *LDAPServerStateSyncer) Work() {
	ctx := context.Background() // TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	// Purge expired ldap_sessions and ldap_user_api_tokens
	recordCreationStaleThreshold := ldSync.config.SessionTimeout().Before(time.Now())
	err := ldSync.deleteStaleSessions(ctx, recordCreationStaleThreshold)
	if err != nil {
		ldSync.lggr.Error("unable to expire local LDAP sessions: ", err)
	}
	recordCreationStaleThreshold = ldSync.config.UserAPITokenDuration().Before(time.Now())
	err = ldSync.deleteStaleAPITokens(ctx, recordCreationStaleThreshold)
	if err != nil {
		ldSync.lggr.Error("unable to expire user API tokens: ", err)
	}

	// Optional rate limiting check to limit the amount of upstream LDAP server queries performed
	if !ldSync.config.UpstreamSyncRateLimit().IsInstant() {
		if !time.Now().After(ldSync.nextSyncTime) {
			return
		}

		// Enough time has elapsed to sync again, store the time for when next sync is allowed and begin sync
		ldSync.nextSyncTime = time.Now().Add(ldSync.config.UpstreamSyncRateLimit().Duration())
	}

	ldSync.lggr.Info("Begin Upstream LDAP provider state sync after checking time against config UpstreamSyncInterval and UpstreamSyncRateLimit")

	// For each defined role/group, query for the list of group members to gather the full list of possible users
	users := []sessions.User{}

	conn, err := ldSync.ldapClient.CreateEphemeralConnection()
	if err != nil {
		ldSync.lggr.Error("Failed to Dial LDAP Server: ", err)
		return
	}
	// Root level root user auth with credentials provided from config
	bindStr := ldSync.config.BaseUserAttr() + "=" + ldSync.config.ReadOnlyUserLogin() + "," + ldSync.config.BaseDN()
	if err = conn.Bind(bindStr, ldSync.config.ReadOnlyUserPass()); err != nil {
		ldSync.lggr.Error("Unable to login as initial root LDAP user: ", err)
	}
	defer conn.Close()

	// Query for list of uniqueMember IDs present in Admin group
	adminUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.AdminUserGroupCN(), sessions.UserRoleAdmin)
	if err != nil {
		ldSync.lggr.Error("Error in ldapGroupMembersListToUser: ", err)
		return
	}
	// Query for list of uniqueMember IDs present in Edit group
	editUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.EditUserGroupCN(), sessions.UserRoleEdit)
	if err != nil {
		ldSync.lggr.Error("Error in ldapGroupMembersListToUser: ", err)
		return
	}
	// Query for list of uniqueMember IDs present in Edit group
	runUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.RunUserGroupCN(), sessions.UserRoleRun)
	if err != nil {
		ldSync.lggr.Error("Error in ldapGroupMembersListToUser: ", err)
		return
	}
	// Query for list of uniqueMember IDs present in Edit group
	readUsers, err := ldSync.ldapGroupMembersListToUser(conn, ldSync.config.ReadUserGroupCN(), sessions.UserRoleView)
	if err != nil {
		ldSync.lggr.Error("Error in ldapGroupMembersListToUser: ", err)
		return
	}

	users = append(users, adminUsers...)
	users = append(users, editUsers...)
	users = append(users, runUsers...)
	users = append(users, readUsers...)

	// Dedupe preserving order of highest role (sorted)
	// Preserve members as a map for future lookup
	upstreamUserStateMap := make(map[string]sessions.User)
	dedupedEmails := []string{}
	for _, user := range users {
		if _, ok := upstreamUserStateMap[user.Email]; !ok {
			upstreamUserStateMap[user.Email] = user
			dedupedEmails = append(dedupedEmails, user.Email)
		}
	}

	// For each unique user in list of active sessions, check for 'Is Active' propery if defined in the config. Some LDAP providers
	// list group members that are no longer marked as active
	usersActiveFlags, err := ldSync.validateUsersActive(dedupedEmails, conn)
	if err != nil {
		ldSync.lggr.Error("Error validating supplied user list: ", err)
	}
	// Remove users in the upstreamUserStateMap source of truth who are part of groups but marked as deactivated/no-active
	for i, active := range usersActiveFlags {
		if !active {
			delete(upstreamUserStateMap, dedupedEmails[i])
		}
	}

	// upstreamUserStateMap is now the most up to date source of truth
	// Now sync database sessions and roles with new data
	err = sqlutil.TransactDataSource(ctx, ldSync.ds, nil, func(tx sqlutil.DataSource) error {
		// First, purge users present in the local ldap_sessions table but not in the upstream server
		type LDAPSession struct {
			UserEmail string
			UserRole  sessions.UserRole
		}
		var existingSessions []LDAPSession
		if err = tx.SelectContext(ctx, &existingSessions, "SELECT user_email, user_role FROM ldap_sessions WHERE localauth_user = false"); err != nil {
			return fmt.Errorf("unable to query ldap_sessions table: %w", err)
		}
		var existingAPITokens []LDAPSession
		if err = tx.SelectContext(ctx, &existingAPITokens, "SELECT user_email, user_role FROM ldap_user_api_tokens WHERE localauth_user = false"); err != nil {
			return fmt.Errorf("unable to query ldap_user_api_tokens table: %w", err)
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
			_, err = tx.ExecContext(ctx, "DELETE FROM ldap_sessions WHERE user_email = ANY($1)", pq.Array(emailsToPurge))
			if err != nil {
				return err
			}
		}

		// Remove any active API tokens this user may have
		if len(apiTokenEmailsToPurge) > 0 {
			_, err = tx.ExecContext(ctx, "DELETE FROM ldap_user_api_tokens WHERE user_email = ANY($1)", pq.Array(apiTokenEmailsToPurge))
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

		// If there are remaining user entries to update
		if len(emailValues) != 0 {
			// Set new role state for all rows in single Exec
			query := fmt.Sprintf("UPDATE ldap_sessions SET user_role = CASE %s ELSE user_role END", queryWhenClause)
			_, err = tx.ExecContext(ctx, query, emailValues...)
			if err != nil {
				return err
			}

			// Update role of API tokens as well
			query = fmt.Sprintf("UPDATE ldap_user_api_tokens SET user_role = CASE %s ELSE user_role END", queryWhenClause)
			_, err = tx.ExecContext(ctx, query, emailValues...)
			if err != nil {
				return err
			}
		}

		ldSync.lggr.Info("local ldap_sessions and ldap_user_api_tokens table successfully synced with upstream LDAP state")
		return nil
	})
	if err != nil {
		ldSync.lggr.Error("Error syncing local database state: ", err)
	}
	ldSync.lggr.Info("Upstream LDAP sync complete")
}

// deleteStaleSessions deletes all ldap_sessions before the passed time.
func (ldSync *LDAPServerStateSyncer) deleteStaleSessions(ctx context.Context, before time.Time) error {
	_, err := ldSync.ds.ExecContext(ctx, "DELETE FROM ldap_sessions WHERE created_at < $1", before)
	return err
}

// deleteStaleAPITokens deletes all ldap_user_api_tokens before the passed time.
func (ldSync *LDAPServerStateSyncer) deleteStaleAPITokens(ctx context.Context, before time.Time) error {
	_, err := ldSync.ds.ExecContext(ctx, "DELETE FROM ldap_user_api_tokens WHERE created_at < $1", before)
	return err
}

// ldapGroupMembersListToUser queries the LDAP server given a conn for a list of uniqueMember who are part of the parameterized group
func (ldSync *LDAPServerStateSyncer) ldapGroupMembersListToUser(conn LDAPConn, groupNameCN string, roleToAssign sessions.UserRole) ([]sessions.User, error) {
	users, err := ldapGroupMembersListToUser(
		conn, groupNameCN, roleToAssign, ldSync.config.GroupsDN(),
		ldSync.config.BaseDN(), ldSync.config.QueryTimeout(),
		ldSync.lggr,
	)
	if err != nil {
		ldSync.lggr.Errorf("Error listing members of group (%s): %v", groupNameCN, err)
		return users, errors.New("error searching group members in LDAP directory")
	}
	return users, nil
}

// validateUsersActive performs an additional LDAP server query for the supplied emails, checking the
// returned user data for an 'active' property defined optionally in the config.
// Returns same length bool 'valid' array, order preserved
func (ldSync *LDAPServerStateSyncer) validateUsersActive(emails []string, conn LDAPConn) ([]bool, error) {
	validUsers := make([]bool, len(emails))
	// If active attribute to check is not defined in config, skip
	if ldSync.config.ActiveAttribute() == "" {
		// pre fill with valids
		for i := range emails {
			validUsers[i] = true
		}
		return validUsers, nil
	}

	// Build the full email list query to pull all 'isActive' information for each user specified in one query
	filterQuery := "(|"
	for _, email := range emails {
		escapedEmail := ldap.EscapeFilter(email)
		filterQuery = fmt.Sprintf("%s(%s=%s)", filterQuery, ldSync.config.BaseUserAttr(), escapedEmail)
	}
	filterQuery = fmt.Sprintf("(&%s))", filterQuery)
	searchBaseDN := fmt.Sprintf("%s,%s", ldSync.config.UsersDN(), ldSync.config.BaseDN())
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, int(ldSync.config.QueryTimeout().Seconds()), false,
		filterQuery,
		[]string{ldSync.config.BaseUserAttr(), ldSync.config.ActiveAttribute()},
		nil,
	)
	// Query LDAP server for the ActiveAttribute property of each specified user
	results, err := conn.Search(searchRequest)
	if err != nil {
		ldSync.lggr.Errorf("Error searching user in LDAP query: %v", err)
		return validUsers, errors.New("error searching users in LDAP directory")
	}
	// Ensure user response entries
	if len(results.Entries) == 0 {
		return validUsers, errors.New("no users matching email query")
	}

	// Pull expected ActiveAttribute value from list of string possible values
	// keyed on email for final step to return flag bool list where order is preserved
	emailToActiveMap := make(map[string]bool)
	for _, result := range results.Entries {
		isActiveAttribute := result.GetAttributeValue(ldSync.config.ActiveAttribute())
		uidAttribute := result.GetAttributeValue(ldSync.config.BaseUserAttr())
		emailToActiveMap[uidAttribute] = isActiveAttribute == ldSync.config.ActiveAttributeAllowedValue()
	}
	for i, email := range emails {
		active, ok := emailToActiveMap[email]
		if ok && active {
			validUsers[i] = true
		}
	}

	return validUsers, nil
}
