-- +goose Up
UPDATE keys SET next_nonce = 0 WHERE next_nonce IS NULL;
ALTER TABLE keys ALTER COLUMN next_nonce SET NOT NULL, ALTER COLUMN next_nonce SET DEFAULT 0;

-- +goose Down
ALTER TABLE keys ALTER COLUMN next_nonce SET DEFAULT NULL;
