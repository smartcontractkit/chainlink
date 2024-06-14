-- +goose Up
CREATE TABLE {{.Schema}}.bcf_3266_01 (
    "id" TEXT PRIMARY KEY
);
-- +goose Down
DROP TABLE {{.Schema}}.bcf_3266_01;