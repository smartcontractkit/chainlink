package oidcauth_test

import (
	// "errors"
	"fmt"
	"strings"
	"testing"
	"time"

	// "time"

	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	// "github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/oidcauth"
	// "github.com/smartcontractkit/chainlink/v2/core/sessions/oidc/mocks"
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
	db.GetContext(ctx, user1, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)

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
		db.GetContext(ctx, u, sql, strings.ToLower(u.Email), u.HashedPassword, u.Role)
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
			log.Error("user not found in ListUsers result: %#v", u.Email)
			panic("match not found")
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
	db.GetContext(ctx, user1, sql, strings.ToLower(user1.Email), user1.HashedPassword, user1.Role)

	sessionRequest := sessions.SessionRequest{
		Email:    user1.Email,
		Password: cltest.Password,
	}

	_, err := oidcAuthProvider.CreateSession(ctx, sessionRequest)
	require.NoError(t, err)
}
