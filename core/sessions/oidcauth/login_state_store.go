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

func putPendingAuth(ctx context.Context, ds sqlutil.DataSource, state, verifier, nonce string, expiresAt time.Time) error {
	if _, err := ds.ExecContext(ctx, `DELETE FROM oidc_pending_auth WHERE expires_at <= now()`); err != nil {
		return err
	}
	_, err := ds.ExecContext(ctx, `
		INSERT INTO oidc_pending_auth (state, verifier, nonce, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (state) DO UPDATE
		SET verifier = EXCLUDED.verifier,
		    nonce = EXCLUDED.nonce,
		    expires_at = EXCLUDED.expires_at`,
		state, verifier, nonce, expiresAt)
	return err
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

func reserveDeviceFlow(ctx context.Context, ds sqlutil.DataSource, handle, clientIP string, expiresAt time.Time) error {
	return sqlutil.TransactDataSource(ctx, ds, nil, func(tx sqlutil.DataSource) error {
		if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1::bigint)`, int64(oidcDeviceFlowLock)); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM oidc_device_flows WHERE expires_at <= now()`); err != nil {
			return err
		}
		var total int
		if err := tx.GetContext(ctx, &total, `SELECT count(*) FROM oidc_device_flows`); err != nil {
			return err
		}
		if total >= maxConcurrentDeviceFlows {
			return errTooManyDeviceFlows
		}
		if clientIP != "" {
			var perIP int
			if err := tx.GetContext(ctx, &perIP, `SELECT count(*) FROM oidc_device_flows WHERE client_ip = $1`, clientIP); err != nil {
				return err
			}
			if perIP >= maxDeviceFlowsPerIP {
				return errTooManyDeviceFlowsPerIP
			}
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO oidc_device_flows (handle, client_ip, expires_at)
			VALUES ($1, $2, $3)`, handle, clientIP, expiresAt)
		return err
	})
}

func deviceFlowLive(ctx context.Context, ds sqlutil.DataSource, handle string) (bool, error) {
	var n int
	err := ds.GetContext(ctx, &n, `
		SELECT count(*) FROM oidc_device_flows
		WHERE handle = $1 AND expires_at > now()`, handle)
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

// consumeDeviceFlow deletes a finished flow and returns it. A second caller
// sees unknown. A row that is not finished yet is left in place and reported
// as pending.
func consumeDeviceFlow(ctx context.Context, ds sqlutil.DataSource, handle string) (sessionID, email string, role clsessions.UserRole, terminal, known bool, flowErr error, err error) {
	var row struct {
		SessionID string `db:"session_id"`
		Email     string
		Role      string
		Err       string
		ExpiresAt time.Time `db:"expires_at"`
	}
	err = ds.GetContext(ctx, &row, `
		DELETE FROM oidc_device_flows
		WHERE handle = $1 AND done = true
		RETURNING session_id, email, role, err, expires_at`, handle)
	if err == nil {
		if time.Now().After(row.ExpiresAt) {
			return "", "", "", false, false, nil, nil
		}
		var storedErr error
		if row.Err != "" {
			storedErr = errors.New(row.Err)
		}
		return row.SessionID, row.Email, clsessions.UserRole(row.Role), true, true, storedErr, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", "", "", false, false, nil, err
	}
	var pending int
	if err = ds.GetContext(ctx, &pending, `
		SELECT count(*) FROM oidc_device_flows
		WHERE handle = $1 AND expires_at > now()`, handle); err != nil {
		return "", "", "", false, false, nil, err
	}
	if pending == 1 {
		return "", "", "", false, true, nil, nil
	}
	return "", "", "", false, false, nil, nil
}
