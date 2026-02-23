-- +goose Up
-- Create verification_tasks queue table
-- This table stores tasks that need to be verified by the TaskVerifierProcessor

CREATE TABLE IF NOT EXISTS verification_tasks (
    id BIGSERIAL PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    -- Owner identification (e.g. "CCTPVerifier", "LombardVerifier")
    -- Multiple verifiers share the same table but only consume their own jobs
    owner_id TEXT NOT NULL,

    -- Chain and message identification
    chain_selector TEXT NOT NULL,
    message_id TEXT NOT NULL,

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
    CONSTRAINT verification_tasks_status_check
        CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

-- Index for efficient job consumption
-- Using partial index to only index jobs that can be consumed
CREATE INDEX IF NOT EXISTS idx_verification_tasks_consume
    ON verification_tasks (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

-- Index for efficient stale job reclamation
-- Covers stale 'processing' jobs where started_at has expired
CREATE INDEX IF NOT EXISTS idx_verification_tasks_stale
    ON verification_tasks (owner_id, started_at ASC, id ASC)
    WHERE status = 'processing' AND started_at IS NOT NULL;

-- Index for efficient stats queries
CREATE INDEX IF NOT EXISTS idx_verification_tasks_status
    ON verification_tasks (owner_id, status);

-- Index for chain-specific queries and monitoring
CREATE INDEX IF NOT EXISTS idx_verification_tasks_chain_status
    ON verification_tasks (chain_selector, status);

-- Index for deduplication and message tracking
CREATE INDEX IF NOT EXISTS idx_verification_tasks_chain_message
    ON verification_tasks (chain_selector, message_id);

-- Index for time-based queries
CREATE INDEX IF NOT EXISTS idx_verification_tasks_created
    ON verification_tasks (created_at DESC);

-- Archive table for completed verification tasks
CREATE TABLE IF NOT EXISTS verification_tasks_archive (
    id BIGINT PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL,
    owner_id TEXT NOT NULL,
    chain_selector TEXT NOT NULL,
    message_id TEXT NOT NULL,
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
CREATE INDEX IF NOT EXISTS idx_verification_tasks_archive_completed
    ON verification_tasks_archive (completed_at DESC);

-- Index for archive queries by chain
CREATE INDEX IF NOT EXISTS idx_verification_tasks_archive_chain
    ON verification_tasks_archive (chain_selector, completed_at DESC);

-- Create verification_results queue table
-- This table stores verification results that need to be written to storage

CREATE TABLE IF NOT EXISTS verification_results (
    id BIGSERIAL PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    -- Owner identification (e.g. "CCTPVerifier", "LombardVerifier")
    owner_id TEXT NOT NULL,

    -- Chain and message identification
    chain_selector TEXT NOT NULL,
    message_id TEXT NOT NULL,

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
    CONSTRAINT verification_results_status_check
        CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

-- Index for efficient job consumption
CREATE INDEX IF NOT EXISTS idx_verification_results_consume
    ON verification_results (owner_id, available_at ASC, id ASC)
    WHERE status IN ('pending', 'failed');

-- Index for efficient stale job reclamation
-- Covers stale 'processing' jobs where started_at has expired
CREATE INDEX IF NOT EXISTS idx_verification_results_stale
    ON verification_results (owner_id, started_at ASC, id ASC)
    WHERE status = 'processing' AND started_at IS NOT NULL;

-- Index for efficient stats queries
CREATE INDEX IF NOT EXISTS idx_verification_results_status
    ON verification_results (owner_id, status);

-- Index for chain-specific queries
CREATE INDEX IF NOT EXISTS idx_verification_results_chain_status
    ON verification_results (chain_selector, status);

-- Index for message tracking
CREATE INDEX IF NOT EXISTS idx_verification_results_chain_message
    ON verification_results (chain_selector, message_id);

-- Index for time-based queries
CREATE INDEX IF NOT EXISTS idx_verification_results_created
    ON verification_results (created_at DESC);

-- Archive table for completed verification results
CREATE TABLE IF NOT EXISTS verification_results_archive (
    id BIGINT PRIMARY KEY,
    job_id UUID UNIQUE NOT NULL,
    owner_id TEXT NOT NULL,
    chain_selector TEXT NOT NULL,
    message_id TEXT NOT NULL,
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
CREATE INDEX IF NOT EXISTS idx_verification_results_archive_completed
    ON verification_results_archive (completed_at DESC);

-- Index for archive queries by chain
CREATE INDEX IF NOT EXISTS idx_verification_results_archive_chain
    ON verification_results_archive (chain_selector, completed_at DESC);


-- +goose Down

-- Drop archive tables
DROP TABLE IF EXISTS verification_results_archive;
DROP TABLE IF EXISTS verification_tasks_archive;

-- Drop main tables (results first due to foreign key)
DROP TABLE IF EXISTS verification_results;
DROP TABLE IF EXISTS verification_tasks;


