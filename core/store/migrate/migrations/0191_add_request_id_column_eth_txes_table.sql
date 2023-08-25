-- +goose Up

ALTER TABLE eth_txes ADD COLUMN request_id bytea UNIQUE;

-- +goose Down

ALTER TABLE eth_txes DROP COLUMN request_id;
