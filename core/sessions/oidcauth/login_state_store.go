package oidcauth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
)

// oidcDeviceFlowLock serializes the concurrency-cap check and insert so two
// replicas cannot both pass the count and both insert.
const oidcDeviceFlowLock = 8628002

const (
	// maxPendingAuth bounds oidc_pending_auth so anonymous sign-in traffic
	// cannot grow the table without limit.
	maxPendingAuth = 10000
	// maxPendingAuthPerClient stops one client bucket from taking the whole
	// table. One operator needs one row per sign-in attempt.
	maxPendingAuthPerClient = 50
)

var (
	errTooManyPendingAuth          = errors.New("too many sign-ins in progress")
	errTooManyPendingAuthPerClient = errors.New("too many sign-ins in progress from this address")
)

// putPendingAuth stores the PKCE verifier and nonce for one browser sign-in.
// The caps are soft across replicas: two concurrent inserts can each pass the
// count, which overshoots by at most the number of replicas.
func putPendingAuth(ctx context.Context, ds sqlutil.DataSource, state, verifier, nonce, clientKey string, expiresAt time.Time) error {
	return sqlutil.TransactDataSource(ctx, ds, nil, func(tx sqlutil.DataSource) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM oidc_pending_auth WHERE expires_at <= now()`); err != nil {
			return err
		}
		var total int
		if err := tx.GetContext(ctx, &total, `SELECT count(*) FROM oidc_pending_auth`); err != nil {
			return err
		}
		if total >= maxPendingAuth {
			return errTooManyPendingAuth
		}
		var perClient int
		if err := tx.GetContext(ctx, &perClient, `SELECT count(*) FROM oidc_pending_auth WHERE client_key = $1`, clientKey); err != nil {
			return err
		}
		if perClient >= maxPendingAuthPerClient {
			return errTooManyPendingAuthPerClient
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO oidc_pending_auth (state, verifier, nonce, client_key, expires_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (state) DO UPDATE
			SET verifier = EXCLUDED.verifier,
			    nonce = EXCLUDED.nonce,
			    client_key = EXCLUDED.client_key,
			    expires_at = EXCLUDED.expires_at`,
			state, verifier, nonce, clientKey, expiresAt)
		return err
	})
}

func takePendingAuth(ctx context.Context, ds sqlutil.DataSource, state string) (verifier, nonce string, ok bool, err error) {
	var row struct {
		Verifier string
		Nonce    string
	}
	err = ds.GetContext(ctx, &row, `
		DELETE FROM oidc_pending_auth
		WHERE state = $1 AND expires_at > now()
		RETURNING verifier, nonce`, state)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, err
	}
	return row.Verifier, row.Nonce, true, nil
}

// deviceFlowReservation is the row written before the node calls the
// identity provider.
type deviceFlowReservation struct {
	handle    string
	clientIP  string
	ipBucket  string
	expiresAt time.Time
	pollGrace time.Duration
}

// sweepDeviceFlows drops flows that expired or that nobody has polled within
// their grace period, and deletes any session such a flow minted but never
// handed out.
const sweepDeviceFlows = `
	WITH gone AS (
		DELETE FROM oidc_device_flows
		WHERE expires_at <= now() OR poll_deadline <= now()
		RETURNING session_id
	)
	DELETE FROM oidc_sessions WHERE id IN (SELECT session_id FROM gone WHERE session_id <> '')`

func reserveDeviceFlow(ctx context.Context, ds sqlutil.DataSource, r deviceFlowReservation) error {
	return sqlutil.TransactDataSource(ctx, ds, nil, func(tx sqlutil.DataSource) error {
		if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1::bigint)`, int64(oidcDeviceFlowLock)); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, sweepDeviceFlows); err != nil {
			return err
		}
		var total int
		if err := tx.GetContext(ctx, &total, `SELECT count(*) FROM oidc_device_flows`); err != nil {
			return err
		}
		if total >= maxConcurrentDeviceFlows {
			return errTooManyDeviceFlows
		}
		var perBucket int
		if err := tx.GetContext(ctx, &perBucket, `SELECT count(*) FROM oidc_device_flows WHERE ip_bucket = $1`, r.ipBucket); err != nil {
			return err
		}
		if perBucket >= maxDeviceFlowsPerIP {
			return errTooManyDeviceFlowsPerIP
		}
		grace := int(r.pollGrace / time.Second)
		_, err := tx.ExecContext(ctx, `
			INSERT INTO oidc_device_flows (handle, client_ip, ip_bucket, expires_at, poll_deadline, poll_grace_seconds)
			VALUES ($1, $2, $3, $4, now() + make_interval(secs => $5::integer), $5::integer)`,
			r.handle, r.clientIP, r.ipBucket, r.expiresAt, grace)
		return err
	})
}

// setDeviceFlowTiming records the provider expiry and the poll grace once the
// provider has answered.
func setDeviceFlowTiming(ctx context.Context, ds sqlutil.DataSource, handle string, expiresAt time.Time, pollGrace time.Duration) error {
	grace := int(pollGrace / time.Second)
	_, err := ds.ExecContext(ctx, `
		UPDATE oidc_device_flows
		SET expires_at = $2,
		    poll_grace_seconds = $3::integer,
		    poll_deadline = now() + make_interval(secs => $3::integer)
		WHERE handle = $1`, handle, expiresAt, grace)
	return err
}

// deviceFlowLive reports whether the flow is unexpired and still being polled.
func deviceFlowLive(ctx context.Context, ds sqlutil.DataSource, handle string) (bool, error) {
	var n int
	err := ds.GetContext(ctx, &n, `
		SELECT count(*) FROM oidc_device_flows
		WHERE handle = $1 AND expires_at > now() AND poll_deadline > now()`, handle)
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func removeDeviceFlow(ctx context.Context, ds sqlutil.DataSource, handle string) error {
	_, err := ds.ExecContext(ctx, `DELETE FROM oidc_device_flows WHERE handle = $1`, handle)
	return err
}

func recordDeviceFlowResult(ctx context.Context, ds sqlutil.DataSource, handle, sessionID, email, role, errText string) error {
	_, err := ds.ExecContext(ctx, `
		UPDATE oidc_device_flows
		SET done = true, session_id = $2, email = $3, role = $4, err = $5
		WHERE handle = $1`, handle, sessionID, email, role, errText)
	return err
}

// deviceFlowResult is what a poll learns about a flow.
type deviceFlowResult struct {
	known     bool
	terminal  bool
	sessionID string
	email     string
	role      clsessions.UserRole
	flowErr   error
	// startIP is the client address that started the flow.
	startIP string
}

// consumeDeviceFlow deletes a finished flow and returns it. A second caller
// sees unknown. A live row that is not finished yet has its poll deadline
// moved forward and is reported as pending. A row past its expiry or poll
// deadline is unknown.
func consumeDeviceFlow(ctx context.Context, ds sqlutil.DataSource, handle string) (deviceFlowResult, error) {
	var row struct {
		SessionID string `db:"session_id"`
		Email     string
		Role      string
		Err       string
		ClientIP  string `db:"client_ip"`
	}
	err := ds.GetContext(ctx, &row, `
		DELETE FROM oidc_device_flows
		WHERE handle = $1 AND done = true AND expires_at > now() AND poll_deadline > now()
		RETURNING session_id, email, role, err, client_ip`, handle)
	if err == nil {
		res := deviceFlowResult{
			known:     true,
			terminal:  true,
			sessionID: row.SessionID,
			email:     row.Email,
			role:      clsessions.UserRole(row.Role),
			startIP:   row.ClientIP,
		}
		if row.Err != "" {
			res.flowErr = errors.New(row.Err)
		}
		return res, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return deviceFlowResult{}, err
	}
	res, err := ds.ExecContext(ctx, `
		UPDATE oidc_device_flows
		SET poll_deadline = now() + make_interval(secs => poll_grace_seconds)
		WHERE handle = $1 AND expires_at > now() AND poll_deadline > now()`, handle)
	// No done = false filter: if the poller finished between the DELETE and
	// this UPDATE, the row is reported pending and the next poll takes it.
	if err != nil {
		return deviceFlowResult{}, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return deviceFlowResult{}, err
	}
	return deviceFlowResult{known: n == 1}, nil
}
