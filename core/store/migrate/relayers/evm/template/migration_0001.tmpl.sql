-- +goose Up
CREATE TABLE {{.Schema}}.bcf_3266_01 (
    name TEXT PRIMARY KEY,
);
-- +goose Down
DROP TABLE {{.Schema}}.bcf_3266_01;