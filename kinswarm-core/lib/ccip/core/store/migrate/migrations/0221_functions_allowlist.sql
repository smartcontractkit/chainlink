-- +goose Up
-- +goose StatementBegin
CREATE TABLE functions_allowlist(
    id BIGSERIAL,
    router_contract_address bytea CHECK (octet_length(router_contract_address) = 20) NOT NULL,
    allowed_address bytea CHECK (octet_length(allowed_address) = 20) NOT NULL,
    PRIMARY KEY(router_contract_address, allowed_address)
);

ALTER TABLE functions_subscriptions
ADD CONSTRAINT router_contract_address_octet_length CHECK (octet_length(router_contract_address) = 20);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE functions_subscriptions
DROP CONSTRAINT router_contract_address_octet_length;

DROP TABLE IF EXISTS functions_allowlist;
-- +goose StatementEnd
