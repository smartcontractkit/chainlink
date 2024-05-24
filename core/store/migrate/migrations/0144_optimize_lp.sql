-- +goose Up
-- Rebuild the indexes without the hex encoding. Its not required postgres can handle bytea comparisons.
DROP INDEX logs_idx_data_word_one, logs_idx_data_word_two, logs_idx_data_word_three, logs_idx_topic_two, logs_idx_topic_three, logs_idx_topic_four;
CREATE INDEX logs_idx_data_word_one ON logs (substring(data from 1 for 32));
CREATE INDEX logs_idx_data_word_two ON logs (substring(data from 33 for 32));
CREATE INDEX logs_idx_data_word_three ON logs (substring(data from 65 for 32));

-- You can only index 3 event arguments. First topic is the event sig which we already have indexed separately.
CREATE INDEX logs_idx_topic_two ON logs ((topics[2]));
CREATE INDEX logs_idx_topic_three ON logs ((topics[3]));
CREATE INDEX logs_idx_topic_four ON logs ((topics[4]));

DROP INDEX logs_idx_evm_id_event_address_block;
DROP INDEX logs_idx_block_number;
ALTER TABLE log_poller_blocks ADD CONSTRAINT block_hash_uniq UNIQUE(evm_chain_id,block_hash);
--
-- +goose Down
DROP INDEX IF EXISTS logs_idx_data_word_one, logs_idx_data_word_two, logs_idx_data_word_three, logs_idx_topic_two, logs_idx_topic_three, logs_idx_topic_four;
CREATE INDEX logs_idx_block_number ON logs using brin(block_number);
CREATE INDEX logs_idx_evm_id_event_address_block ON logs using btree (evm_chain_id,event_sig,address,block_number);
ALTER TABLE log_poller_blocks DROP CONSTRAINT IF EXISTS block_hash_uniq;
