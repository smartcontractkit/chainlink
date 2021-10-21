-- +goose Up
CREATE TYPE log_level AS ENUM (
	'debug',
	'info',
	'warn',
	'error',
	'panic'
);

CREATE TABLE log_configs (
	"id" BIGSERIAL PRIMARY KEY,
	"service_name" text NOT NULL UNIQUE,
	"log_level" log_level NOT NULL,
	"created_at" timestamp with time zone,
	"updated_at" timestamp with time zone
);

-- +goose Down
DROP TABLE IF EXISTS log_configs;
DROP TYPE IF EXISTS log_level;
