-- +goose Up
ALTER TABLE direct_request_specs RENAME COLUMN num_confirmations TO min_incoming_confirmations;
-- +goose Down
ALTER TABLE direct_request_specs RENAME COLUMN min_incoming_confirmations TO num_confirmations;
