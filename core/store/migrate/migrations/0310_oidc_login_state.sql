-- +goose Up
-- Short-lived login state shared by every replica. The provider device code
-- is not stored here. Only the opaque handle and the browser state key are.
-- client_key and ip_bucket are the rate-limit bucket for the caller (an IPv4
-- address or an IPv6 /64), used by the per-client caps.
CREATE TABLE IF NOT EXISTS oidc_pending_auth (
    state text PRIMARY KEY,
    verifier text NOT NULL,
    nonce text NOT NULL,
    client_key text NOT NULL DEFAULT '',
    expires_at timestamp WITH time zone NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_oidc_pending_auth_client_key
    ON oidc_pending_auth (client_key);

-- poll_deadline moves forward each time the CLI polls. A flow nobody polls
-- for poll_grace_seconds stops counting against the caps and is dropped.
CREATE TABLE IF NOT EXISTS oidc_device_flows (
    handle text PRIMARY KEY,
    client_ip text NOT NULL DEFAULT '',
    ip_bucket text NOT NULL DEFAULT '',
    expires_at timestamp WITH time zone NOT NULL,
    poll_deadline timestamp WITH time zone NOT NULL,
    poll_grace_seconds integer NOT NULL,
    done boolean NOT NULL DEFAULT false,
    session_id text NOT NULL DEFAULT '',
    email text NOT NULL DEFAULT '',
    role text NOT NULL DEFAULT '',
    err text NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_oidc_device_flows_ip_bucket
    ON oidc_device_flows (ip_bucket);

-- +goose Down
DROP TABLE IF EXISTS oidc_device_flows;
DROP TABLE IF EXISTS oidc_pending_auth;
