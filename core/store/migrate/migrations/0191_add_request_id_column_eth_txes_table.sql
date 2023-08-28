-- +goose Up

ALTER TABLE eth_txes ADD COLUMN request_id bytea UNIQUE CHECK (octet_length(request_id) <= 2000);

-- +goose Down

ALTER TABLE eth_txes DROP COLUMN request_id;
