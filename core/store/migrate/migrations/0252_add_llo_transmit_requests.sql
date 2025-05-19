-- +goose Up

CREATE TABLE llo_mercury_transmit_queue (
  don_id BIGINT NOT NULL,
  server_url TEXT NOT NULL,
  config_digest BYTEA NOT NULL,
  seq_nr BIGINT NOT NULL,
  report BYTEA NOT NULL,
  lifecycle_stage TEXT NOT NULL,
  report_format BIGINT NOT NULL,
  signatures BYTEA[] NOT NULL,
  signers SMALLINT[] NOT NULL,
  transmission_hash BYTEA NOT NULL,
  PRIMARY KEY (transmission_hash)
);

 CREATE INDEX idx_llo_mercury_transmit_queue_don_id_server_url_seq_nr ON llo_mercury_transmit_queue (don_id, server_url, seq_nr DESC);

-- +goose Down

DROP TABLE llo_mercury_transmit_queue;
