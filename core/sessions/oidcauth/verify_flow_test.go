package oidcauth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

const testClientID = "test-client-id"

// countOIDCSessions returns the number of rows for an email in oidc_sessions.
func countOIDCSessions(t *testing.T, oi *oidcAuthenticator, email string) int {
	t.Helper()
	var n int
	require.NoError(t, oi.ds.GetContext(context.Background(), &n,
		"SELECT count(*) FROM oidc_sessions WHERE lower(user_email) = lower($1)", email))
	return n
}

// TestIssueSessionFromIDToken_Valid runs a correctly signed, correctly
// audienced token through the real verifier and asserts a session row lands.
func TestIssueSessionFromIDToken_Valid(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	idp := newMockIDP(t, testClientID)
	idp.email = "valid-user@example.com"
	idp.groups = []string{AdminClaim}
	oi := newAuthenticatorForIDP(t, idp, db)

	raw := idp.signIDToken(t)
	sessionID, email, err := oi.issueSessionFromIDToken(context.Background(), raw)
	require.NoError(t, err)
	assert.NotEmpty(t, sessionID)
	assert.Equal(t, "valid-user@example.com", email)
	assert.Equal(t, 1, countOIDCSessions(t, oi, "valid-user@example.com"))
}

// TestIssueSessionFromIDToken_WrongAudience is the critical auth-bypass guard:
// a token minted for a different client must be rejected and create no session.
func TestIssueSessionFromIDToken_WrongAudience(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	idp := newMockIDP(t, testClientID)
	idp.email = "attacker@example.com"
	idp.audience = "some-other-client" // token for a different app
	oi := newAuthenticatorForIDP(t, idp, db)

	raw := idp.signIDToken(t)
	_, _, err := oi.issueSessionFromIDToken(context.Background(), raw)
	require.Error(t, err, "token with mismatched audience must be rejected")
	assert.Contains(t, err.Error(), "verify")
	assert.Equal(t, 0, countOIDCSessions(t, oi, "attacker@example.com"))
}

// TestIssueSessionFromIDToken_Expired asserts the verifier enforces expiry.
func TestIssueSessionFromIDToken_Expired(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	idp := newMockIDP(t, testClientID)
	idp.email = "expired-user@example.com"
	idp.lifetime = -time.Hour // already expired
	oi := newAuthenticatorForIDP(t, idp, db)

	raw := idp.signIDToken(t)
	_, _, err := oi.issueSessionFromIDToken(context.Background(), raw)
	require.Error(t, err, "expired token must be rejected")
	assert.Equal(t, 0, countOIDCSessions(t, oi, "expired-user@example.com"))
}

// TestIssueSessionFromIDToken_NoMatchingGroup asserts a validly signed token
// whose groups map to no RBAC role is rejected with errNoMatchingRole and
// creates no session.
func TestIssueSessionFromIDToken_NoMatchingGroup(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	idp := newMockIDP(t, testClientID)
	idp.email = "nogroup-user@example.com"
	idp.groups = []string{"SomeUnrelatedGroup"}
	oi := newAuthenticatorForIDP(t, idp, db)

	raw := idp.signIDToken(t)
	_, _, err := oi.issueSessionFromIDToken(context.Background(), raw)
	require.ErrorIs(t, err, errNoMatchingRole)
	assert.Equal(t, 0, countOIDCSessions(t, oi, "nogroup-user@example.com"))
}

// TestIssueSessionFromIDToken_RoleMapping asserts each group claim maps to the
// expected RBAC role on the stored session.
func TestIssueSessionFromIDToken_RoleMapping(t *testing.T) {
	t.Parallel()
	cases := []struct {
		group    string
		wantRole string
	}{
		{AdminClaim, "admin"},
		{EditorClaim, "edit"},
		{RunnerClaim, "run"},
		{ReadClaim, "view"},
	}
	for _, tc := range cases {
		t.Run(tc.group, func(t *testing.T) {
			t.Parallel()
			db := pgtest.NewSqlxDB(t)
			idp := newMockIDP(t, testClientID)
			email := tc.group + "@example.com"
			idp.email = email
			idp.groups = []string{tc.group}
			oi := newAuthenticatorForIDP(t, idp, db)

			_, _, err := oi.issueSessionFromIDToken(context.Background(), idp.signIDToken(t))
			require.NoError(t, err)

			var role string
			require.NoError(t, oi.ds.GetContext(context.Background(), &role,
				"SELECT user_role FROM oidc_sessions WHERE lower(user_email) = lower($1)", email))
			assert.Equal(t, tc.wantRole, role)
		})
	}
}

// TestDeviceFlow_EndToEnd drives the full device flow against the mock IdP: the
// node polls the mock token endpoint, verifies the returned id_token, and
// records a completed flow with a real oidc_sessions row.
func TestDeviceFlow_EndToEnd(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	idp := newMockIDP(t, testClientID)
	idp.email = "device-user@example.com"
	idp.groups = []string{EditorClaim}
	oi := newAuthenticatorForIDP(t, idp, db)

	// Kick off a device authorization at the provider, then run the same
	// blocking poll the node uses in production.
	da, err := oi.oauth2Config.DeviceAuth(context.Background())
	require.NoError(t, err)

	state := &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}
	oi.pollDeviceToken("e2e-handle", state, da)

	state.mu.Lock()
	defer state.mu.Unlock()
	require.True(t, state.done)
	require.NoError(t, state.err, "device flow should complete and issue a session")
	assert.NotEmpty(t, state.sessionID)
	assert.Equal(t, "device-user@example.com", state.email)
	assert.Equal(t, 1, countOIDCSessions(t, oi, "device-user@example.com"))
}

// TestDeviceFlow_WrongAudienceRejected asserts the device path uses the same
// verifier: a token for the wrong client fails the flow and creates no session.
func TestDeviceFlow_WrongAudienceRejected(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	idp := newMockIDP(t, testClientID)
	idp.email = "device-attacker@example.com"
	idp.audience = "other-client"
	oi := newAuthenticatorForIDP(t, idp, db)

	da, err := oi.oauth2Config.DeviceAuth(context.Background())
	require.NoError(t, err)

	state := &deviceFlowState{expiresAt: time.Now().Add(time.Minute)}
	oi.pollDeviceToken("e2e-handle-2", state, da)

	state.mu.Lock()
	defer state.mu.Unlock()
	require.True(t, state.done)
	require.Error(t, state.err)
	assert.Empty(t, state.sessionID)
	assert.Equal(t, 0, countOIDCSessions(t, oi, "device-attacker@example.com"))
}
