-- +goose Up
ALTER TABLE eal_specs ADD client_certificate text;
ALTER TABLE eal_specs ADD client_key text;

-- +goose Down
ALTER TABLE eal_specs DROP client_certificate;
ALTER TABLE eal_specs DROP client_key;