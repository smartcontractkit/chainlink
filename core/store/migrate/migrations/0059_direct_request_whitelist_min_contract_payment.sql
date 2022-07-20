-- +goose Up
ALTER TABLE direct_request_specs ADD COLUMN min_contract_payment numeric(78,0); 
-- +goose Down
ALTER TABLE direct_request_specs DROP COLUMN min_contract_payment; 
