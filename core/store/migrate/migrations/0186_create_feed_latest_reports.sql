-- +goose Up

CREATE TABLE feed_latest_reports (
  feed_id BYTEA PRIMARY KEY CHECK (octet_length(feed_id) = 32),
  report BYTEA NOT NULL,
  epoch BIGINT NOT NULL,
  round INT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down

DROP TABLE feed_latest_reports;
