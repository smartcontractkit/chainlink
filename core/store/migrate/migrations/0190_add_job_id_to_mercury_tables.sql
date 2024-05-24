-- +goose Up

DELETE FROM feed_latest_reports;
ALTER TABLE feed_latest_reports ADD COLUMN job_id INTEGER NOT NULL REFERENCES jobs(id) DEFERRABLE INITIALLY IMMEDIATE;
DELETE FROM mercury_transmit_requests;
ALTER TABLE mercury_transmit_requests ADD COLUMN job_id INTEGER NOT NULL REFERENCES jobs(id) DEFERRABLE INITIALLY IMMEDIATE;;

-- +goose Down

ALTER TABLE feed_latest_reports DROP COLUMN job_id;
ALTER TABLE mercury_transmit_requests DROP COLUMN job_id;
