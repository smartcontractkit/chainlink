-- +goose Up
CREATE TABLE IF NOT EXISTS ocr2_blocks (
    id SERIAL PRIMARY KEY,
    config_digest bytea NOT NULL,
    seq_nr numeric(20,0) NOT NULL, -- seq_nr is a uint64 which doesn't fit inside a bigint.
    block bytea NOT NULL
);

ALTER TABLE ocr2_blocks
    ADD CONSTRAINT ocr2_blocks_config_digest_seq_nr_unique
        UNIQUE (config_digest, seq_nr);

-- +goose Down
DROP TABLE ocr2_blocks;
