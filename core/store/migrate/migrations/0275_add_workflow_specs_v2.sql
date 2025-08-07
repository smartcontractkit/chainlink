-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_specs_v2 (
    id              SERIAL PRIMARY KEY,
    workflow        text NOT NULL,
    config          text  DEFAULT '',
    workflow_id     varchar(64) NOT NULL UNIQUE,
    workflow_owner  varchar(40) NOT NULL,
    workflow_name   varchar(255) NOT NULL,
    status          text NOT NULL DEFAULT '',
    binary_url      text  DEFAULT '',
    config_url      text  DEFAULT '',
    created_at      timestamp with time zone NOT NULL,
    updated_at      timestamp with time zone NOT NULL,
    spec_type       varchar(255) DEFAULT 'yaml'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workflow_specs_v2;
-- +goose StatementEnd