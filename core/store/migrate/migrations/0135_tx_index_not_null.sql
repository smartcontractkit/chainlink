-- +goose Up
-- +goose StatementBegin
UPDATE log_broadcasts SET tx_index=-1 WHERE tx_index IS NULL;
ALTER TABLE log_broadcasts ALTER COLUMN tx_index SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE log_broadcasts ALTER COLUMN tx_index DROP NOT NULL;
-- +goose StatementEnd
