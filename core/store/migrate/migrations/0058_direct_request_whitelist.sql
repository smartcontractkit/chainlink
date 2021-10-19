-- +goose Up
ALTER TABLE direct_request_specs ADD COLUMN requesters TEXT; 
-- +goose Down
ALTER TABLE direct_request_specs DROP COLUMN requesters; 
