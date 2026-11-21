package oidcauth

import (
	"strconv"
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

	require.NoError(t, putPendingAuth(ctx, db, "state-1", "verifier-1", "nonce-1", "203.0.113.1", exp))

	verifier, nonce, ok, err := takePendingAuth(ctx, db, "state-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "verifier-1", verifier)
	assert.Equal(t, "nonce-1", nonce)

	_, _, ok, err = takePendingAuth(ctx, db, "state-1")
	require.NoError(t, err)
	assert.False(t, ok)

	require.NoError(t, putPendingAuth(ctx, db, "state-expired", "verifier-2", "nonce-2", "203.0.113.1", time.Now().Add(-time.Second)))
	_, _, ok, err = takePendingAuth(ctx, db, "state-expired")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestPendingAuth_PerClientCap(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(pendingAuthTTL)

	for i := range maxPendingAuthPerClient {
		require.NoError(t, putPendingAuth(ctx, db, "s-"+strconv.Itoa(i), "v", "n", "2001:db8::/64", exp))
	}
	err := putPendingAuth(ctx, db, "s-overflow", "v", "n", "2001:db8::/64", exp)
	require.ErrorIs(t, err, errTooManyPendingAuthPerClient)
	require.NoError(t, putPendingAuth(ctx, db, "s-other", "v", "n", "203.0.113.1", exp))

	// Taking a row frees its slot.
	_, _, ok, err := takePendingAuth(ctx, db, "s-0")
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, putPendingAuth(ctx, db, "s-overflow", "v", "n", "2001:db8::/64", exp))
}

func TestPendingAuth_GlobalCap(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	// Fill to the cap in one statement rather than 10k round trips.
	_, err := db.ExecContext(ctx, `
		INSERT INTO oidc_pending_auth (state, verifier, nonce, client_key, expires_at)
		SELECT 'bulk-' || g, 'v', 'n', 'k-' || (g % 1000), now() + interval '10 minutes'
		FROM generate_series(1, $1) g`, maxPendingAuth)
	require.NoError(t, err)
	err = putPendingAuth(ctx, db, "one-more", "v", "n", "fresh-client", time.Now().Add(pendingAuthTTL))
	require.ErrorIs(t, err, errTooManyPendingAuth)
}

func testReservation(handle, ip string, exp time.Time) deviceFlowReservation {
	return deviceFlowReservation{
		handle:    handle,
		clientIP:  ip,
		ipBucket:  clientBucket(parseHostAddr(ip)),
		expiresAt: exp,
		pollGrace: minDevicePollGrace,
	}
}

func TestReserveDeviceFlow_PerIPCap(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(time.Minute)

	for i := range maxDeviceFlowsPerIP {
		require.NoError(t, reserveDeviceFlow(ctx, db, testReservation(handleN(i), "203.0.113.10", exp)))
	}
	err := reserveDeviceFlow(ctx, db, testReservation("overflow", "203.0.113.10", exp))
	require.ErrorIs(t, err, errTooManyDeviceFlowsPerIP)
	require.NoError(t, reserveDeviceFlow(ctx, db, testReservation("other", "198.51.100.1", exp)))
}

// TestReserveDeviceFlow_IPv6SharesSlash64 is the address-rotation case: every
// address in one /64 counts against the same five slots.
func TestReserveDeviceFlow_IPv6SharesSlash64(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(time.Minute)

	for i := range maxDeviceFlowsPerIP {
		require.NoError(t, reserveDeviceFlow(ctx, db, testReservation(handleN(i), "2001:db8::"+strconv.Itoa(i+1), exp)))
	}
	err := reserveDeviceFlow(ctx, db, testReservation("overflow", "2001:db8::ffff", exp))
	require.ErrorIs(t, err, errTooManyDeviceFlowsPerIP)
}

// TestReserveDeviceFlow_AbandonedReleaseSlot: flows nobody polls stop holding
// slots once their poll grace passes, so a start-and-walk-away flood cannot
// hold the pool for the full five minutes.
func TestReserveDeviceFlow_AbandonedReleaseSlot(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(time.Minute)

	for i := range maxDeviceFlowsPerIP {
		require.NoError(t, reserveDeviceFlow(ctx, db, testReservation(handleN(i), "203.0.113.20", exp)))
	}
	require.ErrorIs(t, reserveDeviceFlow(ctx, db, testReservation("blocked", "203.0.113.20", exp)), errTooManyDeviceFlowsPerIP)

	// Keep one flow alive by polling it, then age the others past the grace.
	res, err := consumeDeviceFlow(ctx, db, handleN(0))
	require.NoError(t, err)
	require.True(t, res.known)
	_, err = db.ExecContext(ctx, `UPDATE oidc_device_flows SET poll_deadline = now() - interval '1 second' WHERE handle <> $1`, handleN(0))
	require.NoError(t, err)

	live, err := deviceFlowLive(ctx, db, handleN(1))
	require.NoError(t, err)
	assert.False(t, live, "an unpolled flow is no longer live")
	res, err = consumeDeviceFlow(ctx, db, handleN(1))
	require.NoError(t, err)
	assert.False(t, res.known, "an unpolled flow cannot be redeemed")

	require.NoError(t, reserveDeviceFlow(ctx, db, testReservation("unblocked", "203.0.113.20", exp)))
	live, err = deviceFlowLive(ctx, db, handleN(0))
	require.NoError(t, err)
	assert.True(t, live, "a polled flow survives the sweep")
}

// TestReserveDeviceFlow_SweepDeletesUnclaimedSession: a flow that finished
// but was never redeemed takes its minted session with it.
func TestReserveDeviceFlow_SweepDeletesUnclaimedSession(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(time.Minute)

	require.NoError(t, reserveDeviceFlow(ctx, db, testReservation("orphan", "203.0.113.30", exp)))
	_, err := db.ExecContext(ctx, `INSERT INTO oidc_sessions (id, user_email, user_role, created_at) VALUES ('orphan-sess', 'o@example.com', 'view', now())`)
	require.NoError(t, err)
	require.NoError(t, recordDeviceFlowResult(ctx, db, "orphan", "orphan-sess", "o@example.com", "view", ""))
	_, err = db.ExecContext(ctx, `UPDATE oidc_device_flows SET poll_deadline = now() - interval '1 second' WHERE handle = 'orphan'`)
	require.NoError(t, err)

	require.NoError(t, reserveDeviceFlow(ctx, db, testReservation("next", "203.0.113.31", exp)))
	var n int
	require.NoError(t, db.GetContext(ctx, &n, `SELECT count(*) FROM oidc_sessions WHERE id = 'orphan-sess'`))
	assert.Equal(t, 0, n)
}

func TestConsumeDeviceFlow_SingleWinner(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	ctx := t.Context()
	exp := time.Now().Add(time.Minute)
	require.NoError(t, reserveDeviceFlow(ctx, db, testReservation("shared-handle", "203.0.113.9", exp)))

	res, err := consumeDeviceFlow(ctx, db, "shared-handle")
	require.NoError(t, err)
	assert.False(t, res.terminal)
	assert.True(t, res.known, "an unfinished flow stays pending")
	assert.Empty(t, res.sessionID)
	require.NoError(t, res.flowErr)

	require.NoError(t, recordDeviceFlowResult(ctx, db, "shared-handle", "sess-winner", "a@example.com", string(clsessions.UserRoleEdit), ""))

	res, err = consumeDeviceFlow(ctx, db, "shared-handle")
	require.NoError(t, err)
	require.True(t, res.terminal)
	require.True(t, res.known)
	assert.Equal(t, "sess-winner", res.sessionID)
	assert.Equal(t, "a@example.com", res.email)
	assert.Equal(t, clsessions.UserRoleEdit, res.role)
	assert.Equal(t, "203.0.113.9", res.startIP)
	require.NoError(t, res.flowErr)

	res, err = consumeDeviceFlow(ctx, db, "shared-handle")
	require.NoError(t, err)
	assert.False(t, res.terminal)
	assert.False(t, res.known, "a second poll must not observe the completed flow")
}
