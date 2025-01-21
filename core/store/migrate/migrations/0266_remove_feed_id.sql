-- +goose Up
-- +goose StatementBegin
ALTER TABLE ocr2_oracle_specs DROP COLUMN feed_id;
ALTER TABLE bootstrap_specs DROP COLUMN feed_id;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE ocr2_oracle_specs ADD COLUMN feed_id bytea CHECK (feed_id IS NULL OR octet_length(feed_id) = 32);
ALTER TABLE bootstrap_specs ADD COLUMN feed_id bytea CHECK (feed_id IS NULL OR octet_length(feed_id) = 32);
-- +goose StatementEnd
