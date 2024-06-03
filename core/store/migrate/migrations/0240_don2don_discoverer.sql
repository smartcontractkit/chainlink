-- +goose Up
-- this migration is for the don2don_discoverer_announcements table
-- it is essentially the same as ocr_discoverer_announcements but scoped to the don2don use case
CREATE TABLE don2don_discoverer_announcements (
	local_peer_id text NOT NULL REFERENCES encrypted_p2p_keys (peer_id) DEFERRABLE INITIALLY IMMEDIATE,
	remote_peer_id text NOT NULL,
	ann bytea NOT NULL,
	created_at timestamptz not null,
	updated_at timestamptz not null,
	PRIMARY KEY(local_peer_id, remote_peer_id)
);
-- +goose Down
DROP TABLE don2don_discoverer_announcements;
