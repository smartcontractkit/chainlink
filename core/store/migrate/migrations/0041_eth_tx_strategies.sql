-- +goose Up
ALTER TABLE eth_txes ADD COLUMN subject uuid;
CREATE INDEX idx_eth_txes_unstarted_subject_id ON eth_txes (subject, id) WHERE subject IS NOT NULL AND state = 'unstarted';
-- +goose Down
ALTER TABLE eth_txes DROP COLUMN subject;
