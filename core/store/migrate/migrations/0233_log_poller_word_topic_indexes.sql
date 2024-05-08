-- +goose Up

drop index if exists evm.evm_logs_idx_data_word_one;
drop index if exists evm.evm_logs_idx_data_word_two;
drop index if exists evm.evm_logs_idx_data_word_three;
drop index if exists evm.evm_logs_idx_data_word_four;
drop index if exists evm.evm_logs_idx_topic_two;
drop index if exists evm.evm_logs_idx_topic_three;
drop index if exists evm.evm_logs_idx_topic_four;

create index evm_logs_idx_data_word_one
    on evm.logs (address, event_sig, evm_chain_id, "substring"(data, 1, 32));

create index evm_logs_idx_data_word_two
    on evm.logs (address, event_sig, evm_chain_id, "substring"(data, 33, 32));

create index evm_logs_idx_data_word_three
    on evm.logs (address, event_sig, evm_chain_id, "substring"(data, 65, 32));

create index evm_logs_idx_data_word_four
    on evm.logs (address, event_sig, evm_chain_id, "substring"(data, 97, 32));

create index evm_logs_idx_data_word_five
    on evm.logs (address, event_sig, evm_chain_id, "substring"(data, 129, 32));

create index evm_logs_idx_topic_two
    on evm.logs (address, event_sig, evm_chain_id, (topics[2]));

create index evm_logs_idx_topic_three
    on evm.logs (address, event_sig, evm_chain_id, (topics[3]));

create index evm_logs_idx_topic_four
    on evm.logs (address, event_sig, evm_chain_id, (topics[4]));

-- +goose Down

drop index if exists evm.evm_logs_idx_data_word_one;
drop index if exists evm.evm_logs_idx_data_word_two;
drop index if exists evm.evm_logs_idx_data_word_three;
drop index if exists evm.evm_logs_idx_data_word_four;
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