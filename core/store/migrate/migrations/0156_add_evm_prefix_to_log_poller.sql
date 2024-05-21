-- +goose Up

-- alter log_poller_blocks table, constraints and indices
ALTER TABLE log_poller_blocks RENAME TO evm_log_poller_blocks;

ALTER TABLE evm_log_poller_blocks RENAME CONSTRAINT log_poller_blocks_evm_chain_id_fkey TO evm_log_poller_blocks_evm_chain_id_fkey;

-- alter logs table, constraints and indices
ALTER TABLE logs RENAME TO evm_logs;

ALTER INDEX logs_idx RENAME TO evm_logs_idx;

ALTER INDEX logs_idx_data_word_one RENAME TO evm_logs_idx_data_word_one;
ALTER INDEX logs_idx_data_word_two RENAME TO evm_logs_idx_data_word_two;
ALTER INDEX logs_idx_data_word_three RENAME TO evm_logs_idx_data_word_three;

ALTER INDEX logs_idx_topic_two RENAME TO evm_logs_idx_topic_two;
ALTER INDEX logs_idx_topic_three RENAME TO evm_logs_idx_topic_three;
ALTER INDEX logs_idx_topic_four RENAME TO evm_logs_idx_topic_four;

ALTER TABLE evm_logs RENAME CONSTRAINT logs_evm_chain_id_fkey TO evm_logs_evm_chain_id_fkey;

-- +goose Down

-- alter log_poller_blocks table, constraints and indices
ALTER TABLE evm_log_poller_blocks RENAME TO log_poller_blocks;

ALTER TABLE log_poller_blocks RENAME CONSTRAINT evm_log_poller_blocks_evm_chain_id_fkey TO log_poller_blocks_evm_chain_id_fkey;

-- alter logs table, constraints and indices
ALTER TABLE evm_logs RENAME TO logs;

ALTER INDEX evm_logs_idx RENAME TO logs_idx;

ALTER INDEX evm_logs_idx_data_word_one RENAME TO logs_idx_data_word_one;
ALTER INDEX evm_logs_idx_data_word_two RENAME TO logs_idx_data_word_two;
ALTER INDEX evm_logs_idx_data_word_three RENAME TO logs_idx_data_word_three;

ALTER INDEX evm_logs_idx_topic_two RENAME TO logs_idx_topic_two;
ALTER INDEX evm_logs_idx_topic_three RENAME TO logs_idx_topic_three;
ALTER INDEX evm_logs_idx_topic_four RENAME TO logs_idx_topic_four;

ALTER TABLE logs RENAME CONSTRAINT evm_logs_evm_chain_id_fkey TO logs_evm_chain_id_fkey;
