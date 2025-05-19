-- +goose Up
-- +goose StatementBegin

ALTER TABLE solana.log_poller_filters ADD COLUMN is_backfilled BOOLEAN;
UPDATE solana.log_poller_filters SET is_backfilled = true;
ALTER TABLE solana.log_poller_filters ALTER COLUMN is_backfilled SET NOT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE solana.log_poller_filters DROP COLUMN is_backfilled;
-- +goose StatementEnd
