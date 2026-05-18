-- +goose Up

create index evm_logs_idx_data_word_five
    on evm.logs (address, event_sig, evm_chain_id, "substring"(data, 129, 32));

-- +goose Down

drop index if exists evm.evm_logs_idx_data_word_five;
