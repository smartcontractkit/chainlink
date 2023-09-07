-- +goose Up

ALTER TABLE eth_txes ADD COLUMN idempotency_key varchar(2000) UNIQUE;

-- +goose Down

ALTER TABLE eth_txes DROP COLUMN idempotency_key;
