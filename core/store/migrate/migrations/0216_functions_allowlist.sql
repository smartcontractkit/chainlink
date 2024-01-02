-- +goose Up
-- +goose StatementBegin
CREATE TABLE functions_allowlist(
    id bigint,
    router_contract_address bytea,
    allowed_address bytea CHECK (octet_length(allowed_address) = 20) NOT NULL,
    PRIMARY KEY(router_contract_address, id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS functions_allowlist;
-- +goose StatementEnd
