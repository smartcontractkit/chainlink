-- +goose Up
-- +goose StatementBegin
ALTER TABLE terra_nodes
    DROP COLUMN fcd_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE terra_nodes
    ADD COLUMN fcd_url text CHECK (fcd_url != '');
-- +goose StatementEnd