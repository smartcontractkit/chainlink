-- +goose Up
-- Records which metadata source (e.g. ContractWorkflowSource) produced each
-- persisted workflow spec. 
-- Pre-migration rows default to '' and are filled opportunistically on
-- the next event that touches them.
ALTER TABLE workflow_specs_v2 ADD COLUMN source text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE workflow_specs_v2 DROP COLUMN source;
