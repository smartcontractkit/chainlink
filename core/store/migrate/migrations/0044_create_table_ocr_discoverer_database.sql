-- +goose Up
CREATE TABLE offchainreporting_discoverer_announcements (
	local_peer_id text NOT NULL REFERENCES encrypted_p2p_keys (peer_id) DEFERRABLE INITIALLY IMMEDIATE,
	remote_peer_id text NOT NULL,
	ann bytea NOT NULL,
	created_at timestamptz not null,
	updated_at timestamptz not null,
	PRIMARY KEY(local_peer_id, remote_peer_id)
);
-- +goose Down
DROP TABLE offchainreporting_discoverer_announcements;
