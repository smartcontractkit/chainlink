-- +goose Up
-- Add a `status` column to the `workflow_specs` table.
ALTER TABLE workflow_specs
ADD COLUMN status TEXT DEFAULT '' NOT NULL;

-- +goose Down
-- Remove the `status` column from the `workflow_specs` table.
ALTER TABLE workflow_specs
DROP COLUMN status;