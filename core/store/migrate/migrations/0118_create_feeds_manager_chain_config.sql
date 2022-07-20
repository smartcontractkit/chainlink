-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds_manager_chain_configs (
    id SERIAL PRIMARY KEY,
    chain_id VARCHAR NOT NULL,
    chain_type VARCHAR NOT NULL,
    account_address VARCHAR NOT NULL,
    admin_address VARCHAR NOT NULL,
    feeds_manager_id INTEGER REFERENCES feeds_managers ON DELETE CASCADE,
    flux_monitor_config JSONB,
    ocr1_config JSONB,
    ocr2_config JSONB,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE INDEX idx_feeds_manager_chain_configs_chain_id_chain_type ON feeds_manager_chain_configs(chain_id, chain_type);
CREATE UNIQUE INDEX idx_feeds_manager_chain_configs_chain_id_chain_type_feeds_manager_id ON feeds_manager_chain_configs(chain_id, chain_type, feeds_manager_id);

-- Remove the old configuration columns
ALTER TABLE feeds_managers
DROP CONSTRAINT chk_ocr_bootstrap_peer_multiaddr,
DROP COLUMN job_types,
DROP COLUMN is_ocr_bootstrap_peer,
DROP COLUMN ocr_bootstrap_peer_multiaddr;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds_managers
ADD COLUMN job_types TEXT[],
ADD COLUMN is_ocr_bootstrap_peer boolean NOT NULL DEFAULT false,
ADD COLUMN ocr_bootstrap_peer_multiaddr VARCHAR,
ADD CONSTRAINT chk_ocr_bootstrap_peer_multiaddr CHECK ( NOT (
	is_ocr_bootstrap_peer AND
	(
		ocr_bootstrap_peer_multiaddr IS NULL OR
		ocr_bootstrap_peer_multiaddr = ''
	)
));

DROP TABLE feeds_manager_chain_configs;
-- +goose StatementEnd
