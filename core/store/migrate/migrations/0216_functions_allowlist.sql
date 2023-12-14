-- +goose Up
-- +goose StatementBegin
CREATE TABLE functions_allowlist(
    allowed_address bytea CHECK (octet_length(allowed_address) = 20) PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS functions_allowlist;
-- +goose StatementEnd
