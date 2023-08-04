-- +goose Up

CREATE TABLE mercury_latest_reports (
  feed_id BYTEA PRIMARY KEY CHECK (octet_length(feed_id) = 32),
  report BYTEA NOT NULL,
  updated_at TIMESTAMPTZ
);

-- +goose Down

DROP TABLE mercury_latest_reports;
