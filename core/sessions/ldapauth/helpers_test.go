package ldapauth

import (
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
)

// Returns an instantiated ldapAuthenticator struct without validation for testing
func NewTestLDAPAuthenticator(
	ds sqlutil.DataSource,
	ldapCfg config.LDAP,
	lggr logger.Logger,
	auditLogger audit.AuditLogger,
) (*ldapAuthenticator, error) {
	ldapAuth := ldapAuthenticator{
		ds:          ds,
		ldapClient:  newLDAPClient(ldapCfg),
		config:      ldapCfg,
		lggr:        lggr.Named("LDAPAuthenticationProvider"),
		auditLogger: auditLogger,
	}

	return &ldapAuth, nil
}

// Default server group name mappings for test config and mocked ldap search results
const (
	NodeAdminsGroupCN   = "NodeAdmins"
	NodeEditorsGroupCN  = "NodeEditors"
	NodeRunnersGroupCN  = "NodeRunners"
	NodeReadOnlyGroupCN = "NodeReadOnly"
)

// Implement a setter function within the _test file so that the ldapauth_test module can set the unexported field with a mock
func (l *ldapAuthenticator) SetLDAPClient(newClient LDAPClient) {
	l.ldapClient = newClient
}

// Implements config.LDAP
type TestConfig struct {
}

func (t *TestConfig) ServerAddress() string {
	return "ldaps://MOCK"
}

func (t *TestConfig) ReadOnlyUserLogin() string {
	return "mock-readonly"
}

func (t *TestConfig) ReadOnlyUserPass() string {
	return "mock-password"
}

func (t *TestConfig) ServerTLS() bool {
	return false
}

func (t *TestConfig) SessionTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(time.Duration(0))
}

func (t *TestConfig) QueryTimeout() time.Duration {
	return time.Duration(0)
}

func (t *TestConfig) UserAPITokenDuration() commonconfig.Duration {
	return *commonconfig.MustNewDuration(time.Duration(0))
}

func (t *TestConfig) BaseUserAttr() string {
	return "uid"
}

func (t *TestConfig) BaseDN() string {
	return "dc=custom,dc=example,dc=com"
}

func (t *TestConfig) UsersDN() string {
	return "ou=users"
}

func (t *TestConfig) GroupsDN() string {
	return "ou=groups"
}

func (t *TestConfig) ActiveAttribute() string {
	return "organizationalStatus"
}

func (t *TestConfig) ActiveAttributeAllowedValue() string {
	return "ACTIVE"
}

func (t *TestConfig) AdminUserGroupCN() string {
	return NodeAdminsGroupCN
}

func (t *TestConfig) EditUserGroupCN() string {
	return NodeEditorsGroupCN
}

func (t *TestConfig) RunUserGroupCN() string {
	return NodeRunnersGroupCN
}

func (t *TestConfig) ReadUserGroupCN() string {
	return NodeReadOnlyGroupCN
}

func (t *TestConfig) UserApiTokenEnabled() bool {
	return true
}

func (t *TestConfig) UpstreamSyncInterval() commonconfig.Duration {
	return *commonconfig.MustNewDuration(time.Duration(0))
}

func (t *TestConfig) UpstreamSyncRateLimit() commonconfig.Duration {
	return *commonconfig.MustNewDuration(time.Duration(0))
}
