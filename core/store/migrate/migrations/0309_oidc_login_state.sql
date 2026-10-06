-- +goose Up
-- Short-lived login state shared by every replica. The provider device code
-- is not stored here. Only the opaque handle and the browser state key are.
CREATE TABLE IF NOT EXISTS oidc_pending_auth (
    state text PRIMARY KEY,
    verifier text NOT NULL,
    nonce text NOT NULL,
    expires_at timestamp WITH time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS oidc_device_flows (
    handle text PRIMARY KEY,
    client_ip text NOT NULL DEFAULT '',
    expires_at timestamp WITH time zone NOT NULL,
    done boolean NOT NULL DEFAULT false,
    session_id text NOT NULL DEFAULT '',
    email text NOT NULL DEFAULT '',
    role text NOT NULL DEFAULT '',
    err text NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_oidc_device_flows_client_ip
    ON oidc_device_flows (client_ip);

-- +goose Down
DROP TABLE IF EXISTS oidc_device_flows;
DROP TABLE IF EXISTS oidc_pending_auth;
