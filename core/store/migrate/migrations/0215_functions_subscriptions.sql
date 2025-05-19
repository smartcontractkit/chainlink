-- +goose Up
-- +goose StatementBegin
CREATE TABLE functions_subscriptions(
    router_contract_address bytea,
    subscription_id bigint,
    owner bytea CHECK (octet_length(owner) = 20) NOT NULL,
    balance bigint,
    blocked_balance bigint,
    proposed_owner bytea,
    consumers bytea[],
    flags bytea,
    PRIMARY KEY(router_contract_address, subscription_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS functions_subscriptions;
-- +goose StatementEnd
