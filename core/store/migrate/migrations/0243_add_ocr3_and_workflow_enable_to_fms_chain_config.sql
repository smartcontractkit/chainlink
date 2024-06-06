-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds_manager_chain_configs
ADD COLUMN workflow_config JSONB,
    ADD COLUMN ocr3_capabilities_config JSONB;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds_manager_chain_configs DROP COLUMN workflow,
    DROP COLUMN ocr3_capabilities;

-- +goose StatementEnd
