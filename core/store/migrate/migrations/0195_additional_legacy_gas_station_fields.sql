-- +goose Up
ALTER TABLE legacy_gas_station_sidecar_specs ADD client_certificate text;
ALTER TABLE legacy_gas_station_sidecar_specs ADD client_key text;

-- +goose Down
ALTER TABLE legacy_gas_station_sidecar_specs DROP client_certificate;
ALTER TABLE legacy_gas_station_sidecar_specs DROP client_key;