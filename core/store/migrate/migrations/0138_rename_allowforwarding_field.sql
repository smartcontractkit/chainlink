-- +goose Up
ALTER TABLE jobs RENAME COLUMN allow_forwarding TO forwarding_allowed;
-- +goose Down
ALTER TABLE jobs RENAME COLUMN forwarding_allowed TO allow_forwarding;