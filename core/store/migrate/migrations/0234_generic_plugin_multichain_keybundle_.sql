-- +goose Up
-- +goose StatementBegin
ALTER TABLE ocr2_oracle_specs
    ADD COLUMN onchain_signing_strategy JSONB NOT NULL DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ocr2_oracle_specs
    DROP COLUMN onchain_signing_strategy;
-- +goose StatementEnd