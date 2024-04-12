-- +goose Up

drop index if exists evm.evm_logs_idx_data_word_one;
drop index if exists evm.evm_logs_idx_data_word_two;
drop index if exists evm.evm_logs_idx_data_word_three;
drop index if exists evm.evm_logs_idx_data_word_four;
drop index if exists evm.evm_logs_idx_topic_two;
drop index if exists evm.evm_logs_idx_topic_three;
drop index if exists evm.evm_logs_idx_topic_four;

create index evm_logs_idx_data_word_one
    on evm.logs (evm_chain_id, address, event_sig, "substring"(data, 1, 32));

create index evm_logs_idx_data_word_two
    on evm.logs (evm_chain_id, address, event_sig, "substring"(data, 33, 32));

create index evm_logs_idx_data_word_three
    on evm.logs (evm_chain_id, address, event_sig, "substring"(data, 65, 32));

create index evm_logs_idx_data_word_four
    on evm.logs (evm_chain_id, address, event_sig, "substring"(data, 97, 32));

create index evm_logs_idx_data_word_five
    on evm.logs (evm_chain_id, address, event_sig, "substring"(data, 129, 32));

create index evm_logs_idx_topic_two
    on evm.logs (evm_chain_id, address, event_sig, (topics[2]));

create index evm_logs_idx_topic_three
    on evm.logs (evm_chain_id, address, event_sig, (topics[3]));

create index evm_logs_idx_topic_four
    on evm.logs (evm_chain_id, address, event_sig, (topics[4]));

-- +goose Down

drop index if exists evm.evm_logs_idx_data_word_three;
drop index if exists evm.evm_logs_idx_data_word_five;
drop index if exists evm.evm_logs_idx_topic_two;
drop index if exists evm.evm_logs_idx_topic_three;
drop index if exists evm.evm_logs_idx_topic_four;

create index evm_logs_idx_data_word_one
    on evm.logs ("substring"(data, 1, 32));

create index evm_logs_idx_data_word_two
    on evm.logs ("substring"(data, 33, 32));

create index evm_logs_idx_data_word_three
    on evm.logs ("substring"(data, 65, 32));

create index evm_logs_idx_data_word_four
    on evm.logs ("substring"(data, 97, 32));

create index evm_logs_idx_topic_two
    on evm.logs ((topics[2]));

create index evm_logs_idx_topic_three
    on evm.logs ((topics[3]));

create index evm_logs_idx_topic_four
    on evm.logs ((topics[4]));