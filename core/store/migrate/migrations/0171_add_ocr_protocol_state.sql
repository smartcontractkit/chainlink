-- +goose Up
CREATE TABLE IF NOT EXISTS ocr_mercury_protocol_states (
	config_digest bytea NOT NULL CHECK (octet_length(config_digest) = 32),
	key text NOT NULL CHECK (key != ''),
	value bytea NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ocr_mercury_protocol_states ON ocr_mercury_protocol_states (config_digest, key);


-- +goose Down
DROP TABLE ocr_mercury_protocol_states;
