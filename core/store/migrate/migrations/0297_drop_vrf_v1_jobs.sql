-- +goose Up
-- +goose StatementBegin

-- Remove legacy VRF v1 jobs (pipeline uses [type=vrf], not vrfv2/vrfv2plus).
DELETE FROM jobs
WHERE type = 'vrf'
  AND id IN (
    SELECT j.id
    FROM jobs j
    INNER JOIN job_pipeline_specs jps ON j.id = jps.job_id AND jps.is_primary
    INNER JOIN pipeline_specs ps ON jps.pipeline_spec_id = ps.id
    WHERE ps.dot_dag_source ~ '\[type=vrf\]'
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Deleted VRF v1 jobs cannot be restored.
SELECT 1;

-- +goose StatementEnd
