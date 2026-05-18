-- +goose Up
-- +goose StatementBegin
-- Simplify into one index to improve write performance.
DROP INDEX IF EXISTS idx_terra_msgs_terra_chain_id_state;
DROP INDEX IF EXISTS idx_terra_msgs_terra_chain_id_contract_id_state;
-- We order by state first, then contract_id, to permit efficient queries when grouping unstarted txes
-- across contracts.
CREATE INDEX idx_terra_msgs_terra_chain_id_state_contract_id ON terra_msgs (terra_chain_id, state, contract_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_terra_msgs_terra_chain_id_state_contract_id;
CREATE INDEX idx_terra_msgs_terra_chain_id_state ON terra_msgs (terra_chain_id, state, created_at);
CREATE INDEX  idx_terra_msgs_terra_chain_id_contract_id_state ON terra_msgs(terra_chain_id, contract_id, state);
-- +goose StatementEnd
