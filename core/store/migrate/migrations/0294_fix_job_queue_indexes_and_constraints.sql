-- +goose Up

-- Drop task_job_id column from ccv_storage_writer_jobs and its archive table.
-- The column was never used and removing it makes the schema identical to ccv_task_verifier_jobs.
ALTER TABLE ccv_storage_writer_jobs DROP COLUMN IF EXISTS task_job_id;
ALTER TABLE ccv_storage_writer_jobs_archive DROP COLUMN IF EXISTS task_job_id;


-- Narrow the consume index to status = 'pending' only.
-- Failed jobs are never consumed (they go straight to archive via Fail/Retry), so
-- including 'failed' in the index only adds write overhead with no read benefit.

DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_consume;
CREATE INDEX idx_ccv_task_verifier_jobs_consume
    ON ccv_task_verifier_jobs (owner_id, available_at ASC, id ASC)
    WHERE status = 'pending';

DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_consume;
CREATE INDEX idx_ccv_storage_writer_jobs_consume
    ON ccv_storage_writer_jobs (owner_id, available_at ASC, id ASC)
    WHERE status = 'pending';


-- Add UNIQUE constraint to prevent duplicate jobs.
-- When the verifier restarts, it may try to publish the same job again;
-- ON CONFLICT DO NOTHING in Publish relies on this constraint.

ALTER TABLE ccv_task_verifier_jobs
    ADD CONSTRAINT ccv_task_verifier_jobs_unique_job
    UNIQUE (owner_id, chain_selector, message_id);

ALTER TABLE ccv_storage_writer_jobs
    ADD CONSTRAINT ccv_storage_writer_jobs_unique_job
    UNIQUE (owner_id, chain_selector, message_id);


-- Drop indexes that are not used by any query in postgres_queue.go.
-- Keeping unused indexes wastes storage and slows down every INSERT/UPDATE/DELETE.

-- (chain_selector, status) – no query filters active jobs by chain_selector
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_chain_status;
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_chain_status;

-- (chain_selector, message_id) – superseded by the unique constraint on
-- (owner_id, chain_selector, message_id); no query uses this two-column variant
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_chain_message;
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_chain_message;

-- (created_at DESC) – no query orders or filters by created_at
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_created;
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_created;

-- (chain_selector, completed_at DESC) on archive tables – no query uses this
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_archive_chain;
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_archive_chain;


-- Add message_id index on active and archive tables for debugging.
-- Allows looking up a job by message_id alone without knowing owner_id or chain_selector.
CREATE INDEX idx_ccv_task_verifier_jobs_message_id
    ON ccv_task_verifier_jobs (message_id);
CREATE INDEX idx_ccv_storage_writer_jobs_message_id
    ON ccv_storage_writer_jobs (message_id);
CREATE INDEX idx_ccv_task_verifier_jobs_archive_message_id
    ON ccv_task_verifier_jobs_archive (message_id);
CREATE INDEX idx_ccv_storage_writer_jobs_archive_message_id
    ON ccv_storage_writer_jobs_archive (message_id);


-- Replace archive cleanup index (completed_at DESC) with (owner_id, completed_at DESC).
-- The Cleanup query always filters by both owner_id (equality) and completed_at (range),
-- so leading with owner_id lets Postgres skip to the right owner's rows before the range scan.

DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_archive_completed;
CREATE INDEX idx_ccv_task_verifier_jobs_archive_completed
    ON ccv_task_verifier_jobs_archive (owner_id, completed_at DESC);

DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_archive_completed;
CREATE INDEX idx_ccv_storage_writer_jobs_archive_completed
    ON ccv_storage_writer_jobs_archive (owner_id, completed_at DESC);


-- +goose Down

-- Restore task_job_id column
ALTER TABLE ccv_storage_writer_jobs ADD COLUMN IF NOT EXISTS task_job_id UUID;
ALTER TABLE ccv_storage_writer_jobs_archive ADD COLUMN IF NOT EXISTS task_job_id UUID;

-- Remove unique constraints
ALTER TABLE ccv_task_verifier_jobs
    DROP CONSTRAINT IF EXISTS ccv_task_verifier_jobs_unique_job;

ALTER TABLE ccv_storage_writer_jobs
    DROP CONSTRAINT IF EXISTS ccv_storage_writer_jobs_unique_job;

-- Restore consume indexes with 'failed' status included
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_consume;
CREATE INDEX idx_ccv_task_verifier_jobs_consume
    ON ccv_task_verifier_jobs (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_consume;
CREATE INDEX idx_ccv_storage_writer_jobs_consume
    ON ccv_storage_writer_jobs (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

-- Drop message_id indexes
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_message_id;
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_message_id;
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_archive_message_id;
DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_archive_message_id;

-- Restore dropped indexes

CREATE INDEX idx_ccv_task_verifier_jobs_chain_status
    ON ccv_task_verifier_jobs (chain_selector, status);
CREATE INDEX idx_ccv_storage_writer_jobs_chain_status
    ON ccv_storage_writer_jobs (chain_selector, status);

CREATE INDEX idx_ccv_task_verifier_jobs_chain_message
    ON ccv_task_verifier_jobs (chain_selector, message_id);
CREATE INDEX idx_ccv_storage_writer_jobs_chain_message
    ON ccv_storage_writer_jobs (chain_selector, message_id);

CREATE INDEX idx_ccv_task_verifier_jobs_created
    ON ccv_task_verifier_jobs (created_at DESC);
CREATE INDEX idx_ccv_storage_writer_jobs_created
    ON ccv_storage_writer_jobs (created_at DESC);

CREATE INDEX idx_ccv_task_verifier_jobs_archive_chain
    ON ccv_task_verifier_jobs_archive (chain_selector, completed_at DESC);
CREATE INDEX idx_ccv_storage_writer_jobs_archive_chain
    ON ccv_storage_writer_jobs_archive (chain_selector, completed_at DESC);

-- Restore original archive cleanup index (completed_at only)
DROP INDEX IF EXISTS idx_ccv_task_verifier_jobs_archive_completed;
CREATE INDEX idx_ccv_task_verifier_jobs_archive_completed
    ON ccv_task_verifier_jobs_archive (completed_at DESC);

DROP INDEX IF EXISTS idx_ccv_storage_writer_jobs_archive_completed;
CREATE INDEX idx_ccv_storage_writer_jobs_archive_completed
    ON ccv_storage_writer_jobs_archive (completed_at DESC);
