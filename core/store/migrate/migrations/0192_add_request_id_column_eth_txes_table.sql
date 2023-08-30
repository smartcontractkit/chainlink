-- +goose Up

ALTER TABLE eth_txes ADD COLUMN request_id varchar(2000) UNIQUE;

-- +goose Down

ALTER TABLE eth_txes DROP COLUMN request_id;
