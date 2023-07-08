package ldapauth

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldap"
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
	// Purge expired ldap_sessions
	recordCreationStaleThreshold := ldSync.config.UpstreamSyncInterval().Before(
		ldSync.config.SessionTimeout().Before(time.Now()))
	err := ldSync.deleteStaleSessions(recordCreationStaleThreshold)
	if err != nil {
		ldSync.lggr.Error("unable to reap stale sessions: ", err)
	}

	// For each defined role/group, query for the list of group members to gather the full list of possible users
	users := []sessions.User{}

	// Establish ephemeral connection
	conn, err := ldap.DialURL(ldSync.config.ServerAddress())
	if err != nil {
		ldSync.lggr.Errorf("Failed to Dial LDAP Server", err)
	}
	// Root level root user auth with credentials provided from config
	bindStr := ldSync.config.BaseUserAttr() + "=" + ldSync.config.ReadOnlyUserLogin() + "," + ldSync.config.BaseDn()
	if err := conn.Bind(bindStr, ldSync.config.ReadOnlyUserPass()); err != nil {
		ldSync.lggr.Errorf("Unable to login as initial root LDAP user", err)
	}
	defer conn.Close()

	// Query for list of uniqueMember IDs present in Admin group
	adminUsers, err := ldSync.LDAPGroupMembersListToUser(conn, ldSync.config.AdminUserGroupCn(), sessions.UserRoleAdmin)
	if err != nil {
		ldSync.lggr.Errorf("Error in LDAPGroupMembersListToUser: ", err)
	}
	// Query for list of uniqueMember IDs present in Edit group
	editUsers, err := ldSync.LDAPGroupMembersListToUser(conn, ldSync.config.EditUserGroupCn(), sessions.UserRoleEdit)
	if err != nil {
		ldSync.lggr.Errorf("Error in LDAPGroupMembersListToUser: ", err)
	}
	// Query for list of uniqueMember IDs present in Edit group
	runUsers, err := ldSync.LDAPGroupMembersListToUser(conn, ldSync.config.RunUserGroupCn(), sessions.UserRoleRun)
	if err != nil {
		ldSync.lggr.Errorf("Error in LDAPGroupMembersListToUser: ", err)
	}
	// Query for list of uniqueMember IDs present in Edit group
	readUsers, err := ldSync.LDAPGroupMembersListToUser(conn, ldSync.config.ReadUserGroupCn(), sessions.UserRoleView)
	if err != nil {
		ldSync.lggr.Errorf("Error in LDAPGroupMembersListToUser: ", err)
	}

	users = append(users, adminUsers...)
	users = append(users, editUsers...)
	users = append(users, runUsers...)
	users = append(users, readUsers...)

	// Dedupe preserving order of highest role
	uniqueRef := make(map[string]struct{})
	upstreamUserState := []sessions.User{}
	for _, user := range users {
		if _, ok := uniqueRef[user.Email]; !ok {
			uniqueRef[user.Email] = struct{}{}
			upstreamUserState = append(upstreamUserState, user)
		}
	}

	// upstreamUserState is now the most up to date source of truth
	// Update state of local ldap tables, purging users not present in
	// up to date list, and downgrading roles where applicable
	err = ldSync.q.Transaction(func(tx pg.Queryer) error {
		emailList := []interface{}{}
		for _, user := range upstreamUserState {
			emailList = append(emailList, user.Email)
		}
		placeholders := make([]string, len(emailList))
		for i := range emailList {
			placeholders[i] = "?"
		}
		query := fmt.Sprintf("DELETE FROM ldap_sessions WHERE email NOT IN (%s)", strings.Join(placeholders, ", "))
		_, err := ldSync.q.Exec(query, emailList...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ldSync.lggr.Errorf("Error syncing local database state: ", err)
	}
}

// DeleteStaleSessions deletes all ldap_sessions before the passed time.
func (ldSync *LDAPServerStateSyncer) deleteStaleSessions(before time.Time) error {
	_, err := ldSync.q.Exec("DELETE FROM ldap_sessions WHERE last_used < $1", before)
	return err
}

// LDAPGroupMembersListToUser queries the LDAP server given a conn for a list of uniqueMember who are part of the parameterized group
func (l *LDAPServerStateSyncer) LDAPGroupMembersListToUser(conn *ldap.Conn, groupNameCN string, roleToAssign sessions.UserRole) ([]sessions.User, error) {
	users := []sessions.User{}
	// Prepare and query the GroupsDN for the specified group name
	searchBaseDN := fmt.Sprintf("%s, %s", l.config.GroupsDn(), l.config.BaseDn())
	filterQuery := fmt.Sprintf("(&(cn=%s))", groupNameCN)
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, int(l.config.QueryTimeout().Seconds()), false,
		filterQuery,
		[]string{"uniqueMember"},
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
