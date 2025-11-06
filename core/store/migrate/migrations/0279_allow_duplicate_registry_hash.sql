-- +goose Up
ALTER TABLE registry_syncer_states DROP CONSTRAINT registry_syncer_states_data_hash_key;

-- +goose Down
-- +goose StatementBegin
ALTER TABLE registry_syncer_states ADD CONSTRAINT registry_syncer_states_data_hash_key UNIQUE (data_hash);
-- +goose StatementEnd