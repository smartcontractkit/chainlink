-- +goose Up
ALTER TABLE external_initiators DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS token_secret;
