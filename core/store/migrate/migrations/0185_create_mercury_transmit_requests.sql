-- +goose Up

CREATE TABLE mercury_transmit_requests (
  payload_hash BYTEA PRIMARY KEY,
  payload BYTEA NOT NULL,
  config_digest BYTEA NOT NULL,
  epoch INT NOT NULL,
  round INT NOT NULL,
  extra_hash BYTEA NOT NULL
);

CREATE INDEX idx_mercury_transmission_requests_epoch_round ON mercury_transmit_requests (epoch DESC, round DESC);

-- +goose Down

DROP TABLE mercury_transmit_requests;
