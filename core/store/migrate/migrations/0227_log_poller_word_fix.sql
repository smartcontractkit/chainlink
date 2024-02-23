-- +goose Up
DROP INDEX evm.evm_logs_idx_data_word_four;
CREATE INDEX evm_logs_idx_data_word_five ON evm.logs (substring(data from 129 for 32));


-- +goose Down
DROP INDEX IF EXISTS evm.evm_logs_idx_data_word_five;
CREATE INDEX evm_logs_idx_data_word_four ON evm.logs (substring(data from 97 for 32));
