-- +goose Up
-- +goose StatementBegin
CREATE TABLE chain_reader_spec (
                                id              SERIAL PRIMARY KEY,
                                spec     text NOT NULL,
                                created_at      timestamp with time zone NOT NULL,
                                updated_at      timestamp with time zone NOT NULL
);

ALTER TABLE jobs
    ADD COLUMN chain_reader_spec_id INT REFERENCES chain_reader_spec (id);

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE jobs
    DROP COLUMN chain_reader_spec_id;

DROP TABLE chain_reader_spec;
-- +goose StatementEnd
