-- +goose Up
ALTER TABLE initiators ADD COLUMN job_id_topic_filter uuid;
-- +goose Down
ALTER TABLE initiators DROP COLUMN job_id_topic_filter;
