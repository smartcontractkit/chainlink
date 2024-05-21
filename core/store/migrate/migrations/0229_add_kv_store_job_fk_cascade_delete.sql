-- +goose Up

BEGIN;

ALTER TABLE job_kv_store DROP CONSTRAINT job_kv_store_job_id_fkey;
ALTER TABLE job_kv_store
    ADD CONSTRAINT job_kv_store_job_id_fkey
        FOREIGN KEY (job_id)
            REFERENCES jobs(id)
            ON DELETE CASCADE;

COMMIT;

-- +goose Down
BEGIN;

ALTER TABLE job_kv_store DROP CONSTRAINT job_kv_store_job_id_fkey;
ALTER TABLE job_kv_store
    ADD CONSTRAINT job_kv_store_job_id_fkey
        FOREIGN KEY (job_id)
            REFERENCES jobs(id);

COMMIT;
