-- +goose Up

-- Failed jobs should not be consumed, so they don't need to be in the index
-- This improves query performance by reducing the index size

-- Drop and recreate index for ccv_task_verifier_jobs without 'failed' status
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_consume;
CREATE INDEX idx_ccv_task_verifier_jobs_consume
    ON ccv_task_verifier_jobs (owner_id, available_at ASC, id ASC)
    WHERE status = 'pending';

-- Drop and recreate index for ccv_storage_writer_jobs without 'failed' status
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_consume;
CREATE INDEX idx_ccv_storage_writer_jobs_consume
    ON ccv_storage_writer_jobs (owner_id, available_at ASC, id ASC)
    WHERE status = 'pending';


-- Add UNIQUE constraint to prevent duplicate jobs
-- When the verifier restarts, it may try to publish the same job again
-- This constraint ensures duplicates are rejected at the database level

-- Add unique constraint for ccv_task_verifier_jobs
ALTER TABLE ccv_task_verifier_jobs
    ADD CONSTRAINT ccv_task_verifier_jobs_unique_job
    UNIQUE (owner_id, chain_selector, message_id);

-- Add unique constraint for ccv_storage_writer_jobs
ALTER TABLE ccv_storage_writer_jobs
    ADD CONSTRAINT ccv_storage_writer_jobs_unique_job
    UNIQUE (owner_id, chain_selector, message_id);


-- +goose Down

-- Remove unique constraints
ALTER TABLE ccv_storage_writer_jobs
    DROP CONSTRAINT IF EXISTS ccv_storage_writer_jobs_unique_job;

ALTER TABLE ccv_task_verifier_jobs
    DROP CONSTRAINT IF EXISTS ccv_task_verifier_jobs_unique_job;

-- Restore original indexes with 'failed' status
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_consume;
CREATE INDEX idx_ccv_storage_writer_jobs_consume
    ON ccv_storage_writer_jobs (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_consume;
CREATE INDEX idx_ccv_task_verifier_jobs_consume
    ON ccv_task_verifier_jobs (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

