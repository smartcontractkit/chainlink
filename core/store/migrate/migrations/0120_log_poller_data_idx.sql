-- +goose Up
-- We only index the first 3 words for the log. Revisit should we actually have products that need more.
-- The word value range is only helpful for integer based event arguments.
CREATE INDEX logs_idx_data_word_one ON logs (encode(substring(data from 1 for 32), 'hex'));
CREATE INDEX logs_idx_data_word_two ON logs (encode(substring(data from 33 for 32), 'hex'));
CREATE INDEX logs_idx_data_word_three ON logs (encode(substring(data from 65 for 32), 'hex'));
-- You can only index 3 event arguments. First topic is the event sig which we already have indexed separately.
CREATE INDEX logs_idx_topic_two ON logs (encode(topics[2], 'hex'));
CREATE INDEX logs_idx_topic_three ON logs (encode(topics[3], 'hex'));
CREATE INDEX logs_idx_topic_four ON logs (encode(topics[4], 'hex'));

-- +goose Down
DROP INDEX IF EXISTS logs_idx_data_word_one, logs_idx_data_word_two, logs_idx_data_word_three;
DROP INDEX IF EXISTS logs_idx_topic_two, logs_idx_topic_three, logs_idx_topic_four;
