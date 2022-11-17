-- +goose Up

ALTER TABLE bridge_last_value ALTER finished_at TYPE timestamptz;


-- +goose Down
ALTER TABLE bridge_last_value ALTER finished_at TYPE timestamp;
