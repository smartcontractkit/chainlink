-- +goose Up
ALTER TABLE flux_monitor_specs
ADD min_payment varchar(255);

-- +goose Down
ALTER TABLE flux_monitor_specs
DROP COLUMN min_payment;
