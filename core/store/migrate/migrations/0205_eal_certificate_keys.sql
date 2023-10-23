-- +goose Up
ALTER TABLE eal_specs ADD ca_certificate text;

-- +goose Down
ALTER TABLE eal_specs DROP ca_certificate;