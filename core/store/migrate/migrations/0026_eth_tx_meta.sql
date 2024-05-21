-- +goose Up
ALTER TABLE eth_txes ADD COLUMN meta jsonb;
-- +goose Down
ALTER TABLE eth_txes DROP COLUMN meta;
