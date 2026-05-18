-- +goose Up
ALTER TABLE mercury_transmit_requests ADD COLUMN feed_id BYTEA CHECK (feed_id IS NULL OR octet_length(feed_id) = 32);
DROP INDEX idx_mercury_transmission_requests_epoch_round;
CREATE INDEX idx_mercury_transmission_requests_job_id_epoch_round ON mercury_transmit_requests (job_id, epoch DESC, round DESC);
CREATE INDEX idx_mercury_transmit_requests_job_id ON mercury_transmit_requests (job_id);
CREATE INDEX idx_mercury_transmit_requests_feed_id ON mercury_transmit_requests (feed_id);
CREATE INDEX idx_mercury_feed_latest_reports_job_id ON feed_latest_reports (job_id);

-- +goose Down
ALTER TABLE mercury_transmit_requests DROP COLUMN feed_id;
DROP INDEX idx_mercury_transmit_requests_job_id;
DROP INDEX idx_mercury_feed_latest_reports_job_id;
CREATE INDEX idx_mercury_transmission_requests_epoch_round ON mercury_transmit_requests (epoch DESC, round DESC);
DROP INDEX idx_mercury_transmission_requests_job_id_epoch_round;
