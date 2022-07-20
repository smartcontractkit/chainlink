-- +goose Up
-- +goose StatementBegin
-- Needed to speed up FK checks from eth_key_states
CREATE INDEX idx_eth_txes_from_address ON eth_txes (from_address);
-- Since almost all of them are null we can greatly reduce the size of this index by setting the condition
ALTER TABLE eth_txes DROP CONSTRAINT eth_txes_pipeline_task_run_id_key;
CREATE UNIQUE INDEX idx_eth_txes_pipeline_run_task_id ON eth_txes (pipeline_task_run_id) WHERE pipeline_task_run_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_eth_txes_from_address;
ALTER TABLE eth_txes ADD CONSTRAINT eth_txes_pipeline_task_run_id_key UNIQUE (pipeline_task_run_id);
DROP INDEX idx_eth_txes_pipeline_run_task_id;
-- +goose StatementEnd
