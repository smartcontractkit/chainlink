-- +goose Up
-- +goose StatementBegin
CREATE TABLE functions_subscriptions(
    subscription_id bigint PRIMARY KEY,
    owner bytea CHECK (octet_length(owner) = 20) NOT NULL,
    balance bigint,
    blocked_balance bigint,
    proposed_owner bytea,
    consumers bytea[],
    flags bytea
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS functions_subscriptions;
-- +goose StatementEnd
