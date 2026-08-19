-- +goose Up
-- +goose StatementBegin

-- The legacy Mercury OCR2 plugin has been removed. Mercury jobs are
-- 'offchainreporting2' jobs whose ocr2_oracle_specs.plugin_type = 'mercury'.
-- Delete the jobs first (this satisfies the jobs.ocr2_oracle_spec_id FK and
-- cascades to mercury_transmit_requests / feed_latest_reports via their
-- job_id ON DELETE CASCADE constraints), then remove the orphaned specs.
DELETE FROM jobs
 WHERE ocr2_oracle_spec_id IN (
   SELECT id FROM ocr2_oracle_specs WHERE plugin_type = 'mercury'
 );

DELETE FROM ocr2_oracle_specs WHERE plugin_type = 'mercury';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Irreversible: deleted Mercury job/spec rows cannot be restored.
-- (No schema change to revert.)

-- +goose StatementEnd
