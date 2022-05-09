-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_terra_msgs_terra_chain_id_state;
CREATE INDEX idx_terra_msgs_terra_chain_id_state ON terra_msgs (terra_chain_id, state, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_terra_msgs_terra_chain_id_state;
CREATE INDEX idx_terra_msgs_terra_chain_id_state ON terra_msgs (terra_chain_id, state);
-- +goose StatementEnd
