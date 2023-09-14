-- +goose Up
CREATE TABLE streams (
    id text PRIMARY KEY,
    pipeline_spec_id INT REFERENCES pipeline_specs (id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE, 
    created_at timestamp with time zone NOT NULL
);


-- +goose Down
DROP TABLE streams;
