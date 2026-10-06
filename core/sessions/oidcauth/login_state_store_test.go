package oidcauth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
)

func TestPendingAuth_SharedRoundTrip(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(pendingAuthTTL)

	require.NoError(t, putPendingAuth(ctx, db, "state-1", "verifier-1", "nonce-1", exp))

	verifier, nonce, ok, err := takePendingAuth(ctx, db, "state-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "verifier-1", verifier)
	assert.Equal(t, "nonce-1", nonce)

	_, _, ok, err = takePendingAuth(ctx, db, "state-1")
	require.NoError(t, err)
	assert.False(t, ok)

	require.NoError(t, putPendingAuth(ctx, db, "state-expired", "verifier-2", "nonce-2", time.Now().Add(-time.Second)))
	_, _, ok, err = takePendingAuth(ctx, db, "state-expired")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestReserveDeviceFlow_PerIPCap(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(time.Minute)

	for i := range maxDeviceFlowsPerIP {
		require.NoError(t, reserveDeviceFlow(ctx, db, handleN(i), "203.0.113.10", exp))
	}
	err := reserveDeviceFlow(ctx, db, "overflow", "203.0.113.10", exp)
	require.ErrorIs(t, err, errTooManyDeviceFlowsPerIP)
	require.NoError(t, reserveDeviceFlow(ctx, db, "other", "198.51.100.1", exp))
}

func TestConsumeDeviceFlow_SingleWinner(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(time.Minute)
	require.NoError(t, reserveDeviceFlow(ctx, db, "shared-handle", "203.0.113.9", exp))

	var (
		sessionID, email string
		role             clsessions.UserRole
		terminal, known  bool
		flowErr, err     error
	)
	sessionID, _, _, terminal, known, flowErr, err = consumeDeviceFlow(ctx, db, "shared-handle")
	require.NoError(t, err)
	assert.False(t, terminal)
	assert.True(t, known, "an unfinished flow stays pending")
	assert.Empty(t, sessionID)
	require.NoError(t, flowErr)

	require.NoError(t, recordDeviceFlowResult(ctx, db, "shared-handle", "sess-winner", "a@example.com", string(clsessions.UserRoleEdit), ""))

	sessionID, email, role, terminal, known, flowErr, err = consumeDeviceFlow(ctx, db, "shared-handle")
	require.NoError(t, err)
	require.True(t, terminal)
	require.True(t, known)
	assert.Equal(t, "sess-winner", sessionID)
	assert.Equal(t, "a@example.com", email)
	assert.Equal(t, clsessions.UserRoleEdit, role)
	require.NoError(t, flowErr)

	_, _, _, terminal, known, _, err = consumeDeviceFlow(ctx, db, "shared-handle")
	require.NoError(t, err)
	assert.False(t, terminal)
	assert.False(t, known, "a second poll must not observe the completed flow")
}
