-- +goose Up

ALTER TABLE eth_txes ADD COLUMN idempotency_key uuid UNIQUE;

-- +goose Down

ALTER TABLE eth_txes DROP COLUMN idempotency_key;
