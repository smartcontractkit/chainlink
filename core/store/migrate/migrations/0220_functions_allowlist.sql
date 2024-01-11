-- +goose Up
-- +goose StatementBegin
CREATE TABLE functions_allowlist(
    id BIGSERIAL,
    router_contract_address bytea,
    allowed_address bytea CHECK (octet_length(allowed_address) = 20) NOT NULL,
    PRIMARY KEY(router_contract_address, allowed_address)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS functions_allowlist;
-- +goose StatementEnd
