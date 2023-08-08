-- +goose Up

CREATE TABLE feed_latest_reports (
  feed_id BYTEA PRIMARY KEY CHECK (octet_length(feed_id) = 32),
  report BYTEA NOT NULL,
  updated_at TIMESTAMPTZ
);

-- +goose Down

DROP TABLE feed_latest_reports;
