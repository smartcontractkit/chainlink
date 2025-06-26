package oidcauth_test

import (
	// "errors"
	// "fmt"
	"testing"
	// "time"

	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	// "github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
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
