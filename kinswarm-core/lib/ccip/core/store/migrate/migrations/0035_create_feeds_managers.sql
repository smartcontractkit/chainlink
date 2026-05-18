-- +goose Up
CREATE TABLE feeds_managers (
    id BIGSERIAL PRIMARY KEY,
	name VARCHAR (255) NOT NULL,
	uri VARCHAR (255) NOT NULL,
	public_key bytea CHECK (octet_length(public_key) = 32) NOT NULL UNIQUE,
	job_types TEXT [] NOT NULL,
	network VARCHAR (100) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);
-- +goose Down
	DROP TABLE feeds_managers;
