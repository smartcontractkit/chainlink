package oidcauth_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/oidcauth"
)

// Setup oidc Auth authenticator
func setupAuthenticationProvider(t *testing.T) (*sqlx.DB, sessions.AuthenticationProvider) {
	t.Helper()

	cfg := oidcauth.TestConfig{}
	db := pgtest.NewSqlxDB(t)
	oidcAuthProvider, err := oidcauth.NewTestOIDCAuthenticator(db, &cfg, logger.TestLogger(t), &audit.AuditLoggerService{})
	if err != nil {
		t.Fatalf("Error constructing NewTestoidcAuthenticator: %v\n", err)
	}

	return db, oidcAuthProvider
}

func TestORM_FindUser_Empty(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	// Init OIDC authenticator
	_, oidcAuthProvider := setupAuthenticationProvider(t)
	// Find user
	_, err := oidcAuthProvider.FindUser(ctx, "user@test.com")
	require.ErrorContains(t, err, "user not found")
}

func TestORM_FindUser_Single(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	// Init OIDC authenticator
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)
	require.NoError(t, err)

	// Find user
	foundUser, err := oidcAuthProvider.FindUser(ctx, user1.Email)
	if err != nil {
		fmt.Println("error %#w", err)
	}
	require.NoError(t, err)
	require.Equal(t, foundUser.Email, strings.ToLower(user1.Email))
	require.Equal(t, foundUser.Role, user1.Role)
}

func TestORM_FindUserByAPIToken_Success(t *testing.T) {
	ctx := testutils.Context(t)
	// Init OIDC authenticator
	db, oidcAuthProvider := setupAuthenticationProvider(t)

	testEmail := "test@test.com"
	apiToken := "example"
	_, err := db.Exec("INSERT INTO oidc_user_api_tokens values ($1, 'edit', $2, '', '', now())", testEmail, apiToken)
	require.NoError(t, err)

	// Find user
	foundUser, err := oidcAuthProvider.FindUserByAPIToken(ctx, apiToken)
	require.NoError(t, err)
	require.Equal(t, foundUser.Email, testEmail)
	require.Equal(t, sessions.UserRoleEdit, foundUser.Role)
}

func TestORM_FindUserByAPIToken_Expired(t *testing.T) {
	ctx := testutils.Context(t)
	// Init OIDC authenticator
	cfg := oidcauth.TestConfig{}
	db, oidcAuthProvider := setupAuthenticationProvider(t)

	testEmail := "test@test.com"
	apiToken := "example"
	expiredTime := time.Now().Add(-cfg.UserAPITokenDuration().Duration())
	_, err := db.Exec("INSERT INTO oidc_user_api_tokens values ($1, 'edit', $2, '', '', $3)", testEmail, apiToken, expiredTime)
	require.NoError(t, err)

	// Token found but expired. expect error
	_, err = oidcAuthProvider.FindUserByAPIToken(ctx, apiToken)
	require.Equal(t, sessions.ErrUserSessionExpired, err)
}

func TestORM_ListUsers(t *testing.T) {
	ctx := testutils.Context(t)
	// Init OIDC authenticator
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	users := []sessions.User{
		cltest.MustRandomUser(t),
		cltest.MustRandomUser(t),
		cltest.MustRandomUser(t),
		cltest.MustRandomUser(t),
	}

	for _, u := range users {
		// create user
		sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
		_, err := db.ExecContext(ctx, sql, strings.ToLower(u.Email), u.HashedPassword, u.Role)
		require.NoError(t, err)
	}

	// List User
	list, err := oidcAuthProvider.ListUsers(ctx)
	require.NoError(t, err)

	// Check users above were returned
	for _, u := range users {
		match := false
		for _, f := range list {
			if f.Email == u.Email {
				match = true
			}
		}
		if !match {
			t.Errorf("user not found in ListUsers result: %#v", u.Email)
		}
	}
}

func TestORM_CreateSession(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)
	require.NoError(t, err)

	// create session for the user
	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}
	_, err = oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)
}

func TestORM_DeleteSession(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)
	require.NoError(t, err)

	// create session for the user
	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}
	sid, err := oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)

	err = oidcAuthProvider.DeleteUserSession(ctx, sid)
	require.NoError(t, err)
}

func TestORM_ClearNonConcurrentSession(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)
	require.NoError(t, err)

	// create session for the user
	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}
	sid, err := oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)

	err = oidcAuthProvider.ClearNonCurrentSessions(ctx, sid)
	require.NoError(t, err)
}

func Test_AuthorizeUserWithSession_Success(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)
	require.NoError(t, err)

	// create session for the user
	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}
	sid, err := oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)

	// get user from session, expect ok
	user, err := oidcAuthProvider.AuthorizedUserWithSession(ctx, sid)
	require.NoError(t, err)

	require.Equal(t, user1.Email, user.Email)
	require.Equal(t, user1.Role, user.Role)
}

func Test_AuthorizeUserWithSession_Expired(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	cfg := oidcauth.TestConfig{}
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)
	require.NoError(t, err)

	// create session for the user
	session := sessions.NewSession()

	// token expired 4 hours ago
	expiredTime := time.Now().Add(-cfg.SessionTimeout().Duration() - 4*time.Hour)
	_, err = db.ExecContext(ctx,
		"INSERT INTO oidc_sessions (id, user_email, user_role, created_at) VALUES ($1, $2, $3, $4)",
		session.ID,
		strings.ToLower(user1.Email),
		user1.Role,
		expiredTime,
	)
	require.NoError(t, err)

	// get user from session, expect error
	_, err = oidcAuthProvider.AuthorizedUserWithSession(ctx, session.ID)
	require.Equal(t, err, sessions.ErrUserSessionExpired)
}

func Test_AuthorizeUserWithSession_SessionRoleMatchesUserRole(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, sessions.UserRoleView)
	require.NoError(t, err)

	// create session for the user
	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}
	sid, err := oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)

	// get user from session id
	user, err := oidcAuthProvider.AuthorizedUserWithSession(ctx, sid)
	require.NoError(t, err)
	require.Equal(t, user1.Email, user.Email)
	require.Equal(t, sessions.UserRoleView, user.Role)
}

func TestORM_CreateSession_LocalAdminFallbackLogin(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	db, oidcAuthProvider := setupAuthenticationProvider(t)
	user1 := cltest.MustRandomUser(t)

	// create user
	sql := "INSERT INTO users (email, hashed_password, role, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *"
	_, err := db.ExecContext(ctx, sql, strings.ToLower(user1.Email), user1.HashedPassword, "admin")
	require.NoError(t, err)

	// create session with correct password, expect ok
	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}
	_, err = oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)

	// create session with a invalid password, expect error
	sessionRequest = sessions.SessionRequest{
		Email:    user1.Email,
		Password: "incorrect",
	}
	_, err = oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.ErrorContains(t, err, "invalid password")
}

func Test_IDClaimsToUserRole(t *testing.T) {
	t.Parallel()
	cfg := oidcauth.TestConfig{}
	db := pgtest.NewSqlxDB(t)
	oidcAuthProvider, err := oidcauth.NewTestOIDCAuthenticator(db, &cfg, logger.TestLogger(t), &audit.AuditLoggerService{})
	require.NoError(t, err)

	tests := []struct {
		name       string
		idClaims   []string
		adminClaim string
		editClaim  string
		runClaim   string
		readClaim  string
		wantRole   sessions.UserRole
		wantErr    error
	}{
		{
			name:       "Admin role",
			idClaims:   []string{"admin_group", "other_group"},
			adminClaim: "admin_group",
			editClaim:  "edit_group",
			runClaim:   "run_group",
			readClaim:  "read_group",
			wantRole:   sessions.UserRoleAdmin,
			wantErr:    nil,
		},
		{
			name:       "Edit role",
			idClaims:   []string{"edit_group", "other_group"},
			adminClaim: "admin_group",
			editClaim:  "edit_group",
			runClaim:   "run_group",
			readClaim:  "read_group",
			wantRole:   sessions.UserRoleEdit,
			wantErr:    nil,
		},
		{
			name:       "Run role",
			idClaims:   []string{"run_group", "other_group"},
			adminClaim: "admin_group",
			editClaim:  "edit_group",
			runClaim:   "run_group",
			readClaim:  "read_group",
			wantRole:   sessions.UserRoleRun,
			wantErr:    nil,
		},
		{
			name:       "View role",
			idClaims:   []string{"read_group", "other_group"},
			adminClaim: "admin_group",
			editClaim:  "edit_group",
			runClaim:   "run_group",
			readClaim:  "read_group",
			wantRole:   sessions.UserRoleView,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRole, err := oidcAuthProvider.IDClaimsToUserRole(tt.idClaims, tt.adminClaim, tt.editClaim, tt.runClaim, tt.readClaim)
			if !errors.Is(err, nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("err %v", err)
			}

			if gotRole != tt.wantRole {
				t.Errorf("mismatch got %v want %v", gotRole, tt.wantRole)
			}
		})
	}
}

func Test_ExtractIDClaimValues(t *testing.T) {
	t.Parallel()
	cfg := oidcauth.TestConfig{}
	db := pgtest.NewSqlxDB(t)
	oidcAuthProvider, err := oidcauth.NewTestOIDCAuthenticator(db, &cfg, logger.TestLogger(t), &audit.AuditLoggerService{})
	require.NoError(t, err)
	tests := []struct {
		name    string
		claims  map[string]interface{}
		key     string
		want    []string
		wantErr error
	}{
		{
			name:    "String array claim",
			claims:  map[string]interface{}{"groups": []string{"group1", "group2"}},
			key:     "groups",
			want:    []string{"group1", "group2"},
			wantErr: nil,
		},
		{
			name:    "Interface array claim",
			claims:  map[string]interface{}{"groups": []interface{}{"group1", "group2"}},
			key:     "groups",
			want:    []string{"group1", "group2"},
			wantErr: nil,
		},
		{
			name:    "Single string claim",
			claims:  map[string]interface{}{"groups": "group1"},
			key:     "groups",
			want:    []string{"group1"},
			wantErr: nil,
		},
		{
			name:    "Key not found",
			claims:  map[string]interface{}{"other": []string{"group1"}},
			key:     "groups",
			want:    nil,
			wantErr: errors.New("claim 'groups' not found in ID token"),
		},
		{
			name:    "Invalid item type in array",
			claims:  map[string]interface{}{"groups": []interface{}{"group1", 42}},
			key:     "groups",
			want:    nil,
			wantErr: fmt.Errorf("invalid type for item in 'groups': expected string, got %T", 42),
		},
		{
			name:    "Invalid claim type",
			claims:  map[string]interface{}{"groups": 42},
			key:     "groups",
			want:    nil,
			wantErr: fmt.Errorf("claim 'groups' is not a string or array: got %T", 42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := oidcAuthProvider.ExtractIDClaimValues(tt.claims, tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractIDClaimValues() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr != nil && gotErr != nil {
				if gotErr.Error() != tt.wantErr.Error() {
					t.Errorf("ExtractIDClaimValues() gotErr = %v, want %v", gotErr, tt.wantErr)
				}
			} else if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ExtractIDClaimValues() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
