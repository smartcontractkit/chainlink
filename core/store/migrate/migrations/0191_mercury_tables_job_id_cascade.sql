-- +goose Up
ALTER TABLE mercury_transmit_requests DROP CONSTRAINT mercury_transmit_requests_job_id_fkey;
ALTER TABLE mercury_transmit_requests ADD CONSTRAINT mercury_transmit_requests_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE feed_latest_reports DROP CONSTRAINT feed_latest_reports_job_id_fkey;
ALTER TABLE feed_latest_reports ADD CONSTRAINT feed_latest_reports_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

-- +goose Down
ALTER TABLE mercury_transmit_requests DROP CONSTRAINT mercury_transmit_requests_job_id_fkey;
ALTER TABLE mercury_transmit_requests ADD CONSTRAINT mercury_transmit_requests_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs(id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE feed_latest_reports DROP CONSTRAINT feed_latest_reports_job_id_fkey;
ALTER TABLE feed_latest_reports ADD CONSTRAINT feed_latest_reports_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs(id) DEFERRABLE INITIALLY IMMEDIATE;
