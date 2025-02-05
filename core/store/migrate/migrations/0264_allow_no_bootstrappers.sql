-- +goose Up
-- +goose StatementBegin
ALTER TABLE ocr2_oracle_specs
    ADD COLUMN allow_no_bootstrappers BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ocr2_oracle_specs
    DROP COLUMN allow_no_bootstrappers;
-- +goose StatementEnd