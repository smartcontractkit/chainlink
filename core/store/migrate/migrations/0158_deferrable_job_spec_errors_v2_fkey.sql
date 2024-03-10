-- BCF-2095
-- +goose Up
ALTER TABLE job_spec_errors
DROP CONSTRAINT job_spec_errors_v2_job_id_fkey,
ADD CONSTRAINT job_spec_errors_v2_job_id_fkey
	FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
	DEFERRABLE INITIALLY IMMEDIATE;

-- +goose Down
ALTER TABLE job_spec_errors
DROP CONSTRAINT job_spec_errors_v2_job_id_fkey,
ADD CONSTRAINT job_spec_errors_v2_job_id_fkey
	FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE;