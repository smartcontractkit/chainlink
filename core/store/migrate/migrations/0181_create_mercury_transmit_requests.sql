-- +goose Up

CREATE TABLE mercury_transmit_requests (
  payload BYTEA NOT NULL,
  payload_hash TEXT NOT NULL,
  config_digest BYTEA NOT NULL,
  epoch INT NOT NULL,
  round INT NOT NULL,
  extra_hash BYTEA NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_mercury_transmission_requests_payload_hash ON mercury_transmit_requests (payload_hash);
CREATE INDEX idx_mercury_transmission_requests_created_at ON mercury_transmit_requests (created_at);

-- +goose Down

DROP TABLE mercury_transmit_requests;
