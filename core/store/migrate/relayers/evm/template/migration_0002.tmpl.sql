-- +goose Up
CREATE TABLE {{.Schema}}.bcf_3266_02 (
    id SERIAL PRIMARY KEY,
);
-- +goose Down
DROP TABLE {{.Schema}}.bcf_3266_02;