-- +goose Up
ALTER TABLE eth_txes DROP COLUMN access_list;


-- +goose Down
ALTER TABLE eth_txes ADD COLUMN access_list jsonb;
