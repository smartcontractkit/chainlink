-- +goose Up
-- Create ccv_task_verifier_jobs queue table
-- This table stores tasks that need to be verified by the TaskVerifierProcessor

CREATE TABLE IF NOT EXISTS ccv_task_verifier_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL,

    -- Owner identification (e.g. "CCTPVerifier", "LombardVerifier")
    -- Multiple verifiers share the same table but only consume their own jobs
    owner_id TEXT NOT NULL,

    -- Chain and message identification
    chain_selector NUMERIC(20,0) NOT NULL,
    message_id BYTEA NOT NULL,

    -- Job payload stored as JSONB for flexibility
    -- Contains serialized VerificationTask struct
    task_data JSONB NOT NULL,

    -- Job lifecycle state
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,

    -- Retry handling
    attempt_count INT NOT NULL DEFAULT 0,
    retry_deadline TIMESTAMPTZ NOT NULL,
    last_error TEXT,

    -- Constraints
    CONSTRAINT ccv_task_verifier_jobs_status_check
        CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

-- Index for efficient job consumption
-- Using partial index to only index jobs that can be consumed
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_consume
    ON ccv_task_verifier_jobs (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

-- Index for efficient stale job reclamation
-- Covers stale 'processing' jobs where started_at has expired
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_stale
    ON ccv_task_verifier_jobs (owner_id, started_at ASC, id ASC)
    WHERE status = 'processing' AND started_at IS NOT NULL;

-- Index for efficient stats queries
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_status
    ON ccv_task_verifier_jobs (owner_id, status);

-- Index for chain-specific queries and monitoring
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_chain_status
    ON ccv_task_verifier_jobs (chain_selector, status);

-- Index for deduplication and message tracking
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_chain_message
    ON ccv_task_verifier_jobs (chain_selector, message_id);

-- Index for time-based queries
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_created
    ON ccv_task_verifier_jobs (created_at DESC);

-- Archive table for completed verification tasks
CREATE TABLE IF NOT EXISTS ccv_task_verifier_jobs_archive (
    id BIGINT PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL,
    owner_id TEXT NOT NULL,
    chain_selector NUMERIC(20,0) NOT NULL,
    message_id BYTEA NOT NULL,
    task_data JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    available_at TIMESTAMPTZ NOT NULL,
    started_at TIMESTAMPTZ,
    attempt_count INT NOT NULL,
    retry_deadline TIMESTAMPTZ NOT NULL,
    last_error TEXT,
    completed_at TIMESTAMPTZ NOT NULL
);

-- Index for archive cleanup
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_archive_completed
    ON ccv_task_verifier_jobs_archive (completed_at DESC);

-- Index for archive queries by chain
CREATE INDEX IF NOT EXISTS idx_ccv_task_verifier_jobs_archive_chain
    ON ccv_task_verifier_jobs_archive (chain_selector, completed_at DESC);

-- Create ccv_storage_writer_jobs queue table
-- This table stores verification results that need to be written to storage

CREATE TABLE IF NOT EXISTS ccv_storage_writer_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL,

    -- Owner identification (e.g. "CCTPVerifier", "LombardVerifier")
    owner_id TEXT NOT NULL,

    -- Chain and message identification
    chain_selector NUMERIC(20,0) NOT NULL,
    message_id BYTEA NOT NULL,

    -- Job payload stored as JSONB
    -- Contains serialized VerifierNodeResult struct
    task_data JSONB NOT NULL,

    -- Job lifecycle state
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,

    -- Retry handling
    attempt_count INT NOT NULL DEFAULT 0,
    retry_deadline TIMESTAMPTZ NOT NULL,
    last_error TEXT,

    -- Link to source task for traceability (soft reference, no FK constraint)
    -- Tasks and results have independent lifecycles
    task_job_id UUID,

    -- Constraints
    CONSTRAINT ccv_storage_writer_jobs_status_check
        CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

-- Index for efficient job consumption
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_consume
    ON ccv_storage_writer_jobs (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

-- Index for efficient stale job reclamation
-- Covers stale 'processing' jobs where started_at has expired
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_stale
    ON ccv_storage_writer_jobs (owner_id, started_at ASC, id ASC)
    WHERE status = 'processing' AND started_at IS NOT NULL;

-- Index for efficient stats queries
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_status
    ON ccv_storage_writer_jobs (owner_id, status);

-- Index for chain-specific queries
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_chain_status
    ON ccv_storage_writer_jobs (chain_selector, status);

-- Index for message tracking
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_chain_message
    ON ccv_storage_writer_jobs (chain_selector, message_id);

-- Index for time-based queries
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_created
    ON ccv_storage_writer_jobs (created_at DESC);

-- Archive table for completed verification results
CREATE TABLE IF NOT EXISTS ccv_storage_writer_jobs_archive (
    id BIGINT PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL,
    owner_id TEXT NOT NULL,
    chain_selector NUMERIC(20,0) NOT NULL,
    message_id BYTEA NOT NULL,
    task_data JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    available_at TIMESTAMPTZ NOT NULL,
    started_at TIMESTAMPTZ,
    attempt_count INT NOT NULL,
    retry_deadline TIMESTAMPTZ NOT NULL,
    last_error TEXT,
    task_job_id UUID,
    completed_at TIMESTAMPTZ NOT NULL
);

-- Index for archive cleanup
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_archive_completed
    ON ccv_storage_writer_jobs_archive (completed_at DESC);

-- Index for archive queries by chain
CREATE INDEX IF NOT EXISTS idx_ccv_storage_writer_jobs_archive_chain
    ON ccv_storage_writer_jobs_archive (chain_selector, completed_at DESC);


-- +goose Down

-- Drop archive tables
DROP TABLE IF EXISTS ccv_storage_writer_jobs_archive;
DROP TABLE IF EXISTS ccv_task_verifier_jobs_archive;

-- Drop main tables
DROP TABLE IF EXISTS ccv_storage_writer_jobs;
DROP TABLE IF EXISTS ccv_task_verifier_jobs;


