-- +goose Up
-- +goose StatementBegin
ALTER TABLE llo_mercury_transmit_queue ADD COLUMN inserted_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
CREATE INDEX idx_inserted_at ON llo_mercury_transmit_queue USING BRIN (inserted_at);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE llo_mercury_transmit_queue DROP COLUMN inserted_at;
-- +goose StatementEnd
