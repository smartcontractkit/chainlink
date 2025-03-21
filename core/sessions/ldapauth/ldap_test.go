package ldapauth_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/ldapauth"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/ldapauth/mocks"
)

// Setup LDAP Auth authenticator
func setupAuthenticationProvider(t *testing.T, ldapClient ldapauth.LDAPClient) (*sqlx.DB, sessions.AuthenticationProvider) {
	t.Helper()

	cfg := ldapauth.TestConfig{}
	db := pgtest.NewSqlxDB(t)
	ldapAuthProvider, err := ldapauth.NewTestLDAPAuthenticator(db, &cfg, logger.TestLogger(t), &audit.AuditLoggerService{})
	if err != nil {
		t.Fatalf("Error constructing NewTestLDAPAuthenticator: %v\n", err)
	}

	// Override the LDAPClient responsible for returning the *ldap.Conn struct with Mock
	ldapAuthProvider.SetLDAPClient(ldapClient)
	return db, ldapAuthProvider
}

func TestORM_FindUser_Empty(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// User not in upstream, return no entry
	expectedResults := ldap.SearchResult{}

	// On search performed for validateUsersActive
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&expectedResults, nil)

	// Not in upstream, no local admin users, expect error
	_, err := ldapAuthProvider.FindUser(ctx, "unknown-user")
	require.ErrorContains(t, err, "LDAP query returned no matching users")
}

func TestORM_FindUser_NoGroups(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// User present in Upstream but no groups assigned
	user1 := cltest.MustRandomUser(t)
	expectedResults := ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: "cn=User One,ou=Users,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name:   "organizationalStatus",
						Values: []string{"ACTIVE"},
					},
					{
						Name:   "uid",
						Values: []string{user1.Email},
					},
				},
			},
		},
	}

	// On search performed for validateUsersActive
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&expectedResults, nil)

	// No Groups, expect error
	_, err := ldapAuthProvider.FindUser(ctx, user1.Email)
	require.ErrorContains(t, err, "user present in directory, but matching no role groups assigned")
}

func TestORM_FindUser_NotActive(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// User present in Upstream but not active
	user1 := cltest.MustRandomUser(t)
	expectedResults := ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: "cn=User One,ou=Users,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name:   "organizationalStatus",
						Values: []string{"INACTIVE"},
					},
					{
						Name:   "uid",
						Values: []string{user1.Email},
					},
				},
			},
		},
	}

	// On search performed for validateUsersActive
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&expectedResults, nil)

	// User not active, expect error
	_, err := ldapAuthProvider.FindUser(ctx, user1.Email)
	require.ErrorContains(t, err, "user not active")
}

func TestORM_FindUser_Single(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// User present and valid
	user1 := cltest.MustRandomUser(t)
	expectedResults := ldap.SearchResult{ // Users query
		Entries: []*ldap.Entry{
			{
				DN: "cn=User One,ou=Users,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name:   "organizationalStatus",
						Values: []string{"ACTIVE"},
					},
					{
						Name:   "uid",
						Values: []string{user1.Email},
					},
				},
			},
		},
	}
	expectedGroupResults := ldap.SearchResult{ // Groups query
		Entries: []*ldap.Entry{
			{
				DN: "cn=NodeEditors,ou=Users,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name:   "cn",
						Values: []string{"NodeEditors"},
					},
				},
			},
		},
	}

	// On search performed for validateUsersActive
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&expectedResults, nil).Once()

	// Second call on user groups search
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&expectedGroupResults, nil).Once()

	// User active, and has editor group. Expect success
	user, err := ldapAuthProvider.FindUser(ctx, user1.Email)
	require.NoError(t, err)
	require.Equal(t, user1.Email, user.Email)
	require.Equal(t, sessions.UserRoleEdit, user.Role)
}

func TestORM_FindUser_FallbackMatchLocalAdmin(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	// Initilaize LDAP Authentication Provider with mock client
	mockLdapClient := mocks.NewLDAPClient(t)
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// Not in upstream, but utilize text fixture admin user presence in test DB. Succeed
	user, err := ldapAuthProvider.FindUser(ctx, cltest.APIEmailAdmin)
	require.NoError(t, err)
	require.Equal(t, cltest.APIEmailAdmin, user.Email)
	require.Equal(t, sessions.UserRoleAdmin, user.Role)
}

func TestORM_FindUserByAPIToken_Success(t *testing.T) {
	ctx := testutils.Context(t)
	// Initilaize LDAP Authentication Provider with mock client
	mockLdapClient := mocks.NewLDAPClient(t)
	db, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// Ensure valid tokens return a user with role
	testEmail := "test@test.com"
	apiToken := "example"
	_, err := db.Exec("INSERT INTO ldap_user_api_tokens values ($1, 'edit', false, $2, '', '', now())", testEmail, apiToken)
	require.NoError(t, err)

	// Found user by API token in specific ldap_user_api_tokens table
	user, err := ldapAuthProvider.FindUserByAPIToken(ctx, apiToken)
	require.NoError(t, err)
	require.Equal(t, testEmail, user.Email)
	require.Equal(t, sessions.UserRoleEdit, user.Role)
}

func TestORM_FindUserByAPIToken_Expired(t *testing.T) {
	ctx := testutils.Context(t)
	cfg := ldapauth.TestConfig{}

	// Initilaize LDAP Authentication Provider with mock client
	mockLdapClient := mocks.NewLDAPClient(t)
	db, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// Ensure valid tokens return a user with role
	testEmail := "test@test.com"
	apiToken := "example"
	expiredTime := time.Now().Add(-cfg.UserAPITokenDuration().Duration())
	_, err := db.Exec("INSERT INTO ldap_user_api_tokens values ($1, 'edit', false, $2, '', '', $3)", testEmail, apiToken, expiredTime)
	require.NoError(t, err)

	// Token found, but expired. Expect error
	_, err = ldapAuthProvider.FindUserByAPIToken(ctx, apiToken)
	require.Equal(t, sessions.ErrUserSessionExpired, err)
}

func TestORM_ListUsers_Full(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	user1 := cltest.MustRandomUser(t)
	user2 := cltest.MustRandomUser(t)
	user3 := cltest.MustRandomUser(t)
	user4 := cltest.MustRandomUser(t)
	user5 := cltest.MustRandomUser(t)
	user6 := cltest.MustRandomUser(t)

	// LDAP Group queries per role - admin
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeAdminsGroupCN),
				Attributes: []*ldap.EntryAttribute{
					{
						Name: ldapauth.UniqueMemberAttribute,
						Values: []string{
							fmt.Sprintf("uid=%s,ou=users,dc=example,dc=com", user1.Email),
							fmt.Sprintf("uid=%s,ou=users,dc=example,dc=com", user2.Email),
						},
					},
				},
			},
		},
	}, nil).Once()
	// LDAP Group queries per role - edit
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeEditorsGroupCN),
				Attributes: []*ldap.EntryAttribute{
					{
						Name: ldapauth.UniqueMemberAttribute,
						Values: []string{
							fmt.Sprintf("uid=%s,ou=users,dc=example,dc=com", user3.Email),
						},
					},
				},
			},
		},
	}, nil).Once()
	// LDAP Group queries per role - run
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: "cn=NodeRunners,ou=Groups,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name: ldapauth.UniqueMemberAttribute,
						Values: []string{
							fmt.Sprintf("uid=%s,ou=users,dc=example,dc=com", user4.Email),
							fmt.Sprintf("uid=%s,ou=users,dc=example,dc=com", user4.Email), // Test deduped
							fmt.Sprintf("uid=%s,ou=users,dc=example,dc=com", user5.Email),
						},
					},
				},
			},
		},
	}, nil).Once()
	// LDAP Group queries per role - view
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: "cn=NodeReadOnly,ou=Groups,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name: ldapauth.UniqueMemberAttribute,
						Values: []string{
							fmt.Sprintf("uid=%s,ou=users,dc=example,dc=com", user6.Email),
						},
					},
				},
			},
		},
	}, nil).Once()
	// Lastly followed by IsActive lookup
	type userActivePair struct {
		email  string
		active string
	}
	emailsActive := []userActivePair{
		{user1.Email, "ACTIVE"},
		{user2.Email, "INACTIVE"},
		{user3.Email, "ACTIVE"},
		{user4.Email, "ACTIVE"},
		{user5.Email, "INACTIVE"},
		{user6.Email, "ACTIVE"},
	}
	listUpstreamUsersQuery := ldap.SearchResult{}
	for _, upstreamUser := range emailsActive {
		listUpstreamUsersQuery.Entries = append(listUpstreamUsersQuery.Entries, &ldap.Entry{
			DN: "cn=User,ou=Users,dc=example,dc=com",
			Attributes: []*ldap.EntryAttribute{
				{
					Name:   "organizationalStatus",
					Values: []string{upstreamUser.active},
				},
				{
					Name:   "uid",
					Values: []string{upstreamUser.email},
				},
			},
		},
		)
	}
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&listUpstreamUsersQuery, nil).Once()

	// Asserts 'uid=' parsing log in  ldapGroupMembersListToUser
	// Expected full list of users above, including local admin user, excluding 'inactive' and duplicate users
	users, err := ldapAuthProvider.ListUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, users[0].Email, user1.Email)
	require.Equal(t, sessions.UserRoleAdmin, users[0].Role)
	require.Equal(t, users[1].Email, user3.Email) // User 2 inactive
	require.Equal(t, sessions.UserRoleEdit, users[1].Role)
	require.Equal(t, users[2].Email, user4.Email)
	require.Equal(t, sessions.UserRoleRun, users[2].Role)
	require.Equal(t, users[3].Email, user6.Email) // User 5 inactive
	require.Equal(t, sessions.UserRoleView, users[3].Role)
	require.Equal(t, cltest.APIEmailAdmin, users[4].Email) // Text fixture user is local admin included as well
	require.Equal(t, sessions.UserRoleAdmin, users[4].Role)
}

func TestORM_CreateSession_UpstreamBind(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// Upsream user present
	user1 := cltest.MustRandomUser(t)
	expectedResults := ldap.SearchResult{ // Users query
		Entries: []*ldap.Entry{
			{
				DN: "cn=User One,ou=Users,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name:   "organizationalStatus",
						Values: []string{"ACTIVE"},
					},
					{
						Name:   "uid",
						Values: []string{user1.Email},
					},
				},
			},
		},
	}
	expectedGroupResults := ldap.SearchResult{ // Groups query
		Entries: []*ldap.Entry{
			{
				DN: "cn=NodeEditors,ou=Users,dc=example,dc=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name:   "cn",
						Values: []string{"NodeEditors"},
					},
				},
			},
		},
	}

	// On search performed for validateUsersActive
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&expectedResults, nil).Once()

	// Second call on user groups search
	mockLdapConnProvider.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&expectedGroupResults, nil).Once()

	// User active, and has editor group. Expect success
	mockLdapConnProvider.On("Bind", mock.Anything, cltest.Password).Return(nil)
	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}

	_, err := ldapAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)
}

func TestORM_CreateSession_LocalAdminFallbackLogin(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// Fail the bind to trigger 'localLoginFallback' - local admin users should still be able to login
	// regardless of whether the authentication provider is remote or not
	mockLdapConnProvider.On("Bind", mock.Anything, cltest.Password).Return(errors.New("unable to login via LDAP server")).Once()

	// User active, and has editor group. Expect success
	sessionRequest := sessions.SessionRequest{
		Email:    cltest.APIEmailAdmin,
		Password: cltest.Password,
	}

	_, err := ldapAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)

	// Finally, assert login failing altogether
	// User active, and has editor group. Expect success
	mockLdapConnProvider.On("Bind", mock.Anything, "incorrect-password").Return(errors.New("unable to login via LDAP server")).Once()
	sessionRequest = sessions.SessionRequest{
		Email:    cltest.APIEmailAdmin,
		Password: "incorrect-password",
	}

	_, err = ldapAuthProvider.CreateSession(ctx, sessionRequest)
	require.ErrorContains(t, err, "invalid password")
}

func TestORM_SetPassword_LocalAdminFallbackLogin(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	mockLdapClient := mocks.NewLDAPClient(t)
	mockLdapConnProvider := mocks.NewLDAPConn(t)
	mockLdapClient.On("CreateEphemeralConnection").Return(mockLdapConnProvider, nil)
	mockLdapConnProvider.On("Close").Return(nil)

	// Initilaize LDAP Authentication Provider with mock client
	_, ldapAuthProvider := setupAuthenticationProvider(t, mockLdapClient)

	// Fail the bind to trigger 'localLoginFallback' - local admin users should still be able to login
	// regardless of whether the authentication provider is remote or not
	mockLdapConnProvider.On("Bind", mock.Anything, cltest.Password).Return(errors.New("unable to login via LDAP server")).Once()

	// User active, and has editor group. Expect success
	sessionRequest := sessions.SessionRequest{
		Email:    cltest.APIEmailAdmin,
		Password: cltest.Password,
	}

	_, err := ldapAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)

	// Finally, assert login failing altogether
	// User active, and has editor group. Expect success
	mockLdapConnProvider.On("Bind", mock.Anything, "incorrect-password").Return(errors.New("unable to login via LDAP server")).Once()
	sessionRequest = sessions.SessionRequest{
		Email:    cltest.APIEmailAdmin,
		Password: "incorrect-password",
	}

	_, err = ldapAuthProvider.CreateSession(ctx, sessionRequest)
	require.ErrorContains(t, err, "invalid password")
}

func TestORM_MapSearchGroups(t *testing.T) {
	t.Parallel()

	cfg := ldapauth.TestConfig{}

	tests := []struct {
		name                    string
		groupsQuerySearchResult []*ldap.Entry
		wantMappedRole          sessions.UserRole
		wantErr                 error
	}{
		{
			"user in admin group only",
			[]*ldap.Entry{
				{
					DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeAdminsGroupCN),
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "cn",
							Values: []string{ldapauth.NodeAdminsGroupCN},
						},
					},
				},
			},
			sessions.UserRoleAdmin,
			nil,
		},
		{
			"user in edit group",
			[]*ldap.Entry{
				{
					DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeEditorsGroupCN),
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "cn",
							Values: []string{ldapauth.NodeEditorsGroupCN},
						},
					},
				},
			},
			sessions.UserRoleEdit,
			nil,
		},
		{
			"user in run group",
			[]*ldap.Entry{
				{
					DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeRunnersGroupCN),
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "cn",
							Values: []string{ldapauth.NodeRunnersGroupCN},
						},
					},
				},
			},
			sessions.UserRoleRun,
			nil,
		},
		{
			"user in view role",
			[]*ldap.Entry{
				{
					DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeReadOnlyGroupCN),
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "cn",
							Values: []string{ldapauth.NodeReadOnlyGroupCN},
						},
					},
				},
			},
			sessions.UserRoleView,
			nil,
		},
		{
			"user in none",
			[]*ldap.Entry{},
			sessions.UserRole(""), // ignored, error case
			ldapauth.ErrUserNoLDAPGroups,
		},
		{
			"user in run and view",
			[]*ldap.Entry{
				{
					DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeRunnersGroupCN),
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "cn",
							Values: []string{ldapauth.NodeRunnersGroupCN},
						},
					},
				},
				{
					DN: fmt.Sprintf("cn=%s,ou=Groups,dc=example,dc=com", ldapauth.NodeReadOnlyGroupCN),
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "cn",
							Values: []string{ldapauth.NodeReadOnlyGroupCN},
						},
					},
				},
			},
			sessions.UserRoleRun, // Take highest role
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			role, err := ldapauth.GroupSearchResultsToUserRole(
				test.groupsQuerySearchResult,
				cfg.AdminUserGroupCN(),
				cfg.EditUserGroupCN(),
				cfg.RunUserGroupCN(),
				cfg.ReadUserGroupCN(),
			)
			if test.wantErr != nil {
				assert.Equal(t, test.wantErr, err)
			} else {
				assert.Equal(t, test.wantMappedRole, role)
			}
		})
	}
}
