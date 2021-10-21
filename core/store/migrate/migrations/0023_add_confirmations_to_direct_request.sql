-- +goose Up
ALTER TABLE direct_request_specs ADD COLUMN num_confirmations bigint DEFAULT NULL;
-- +goose Down
ALTER TABLE direct_request_specs DROP COLUMN num_confirmations;
