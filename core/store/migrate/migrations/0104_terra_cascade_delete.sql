-- +goose Up
ALTER TABLE terra_nodes
DROP CONSTRAINT terra_nodes_terra_chain_id_fkey;

ALTER TABLE terra_nodes
    ADD FOREIGN KEY (terra_chain_id) REFERENCES terra_chains(id) ON DELETE CASCADE;

ALTER TABLE terra_msgs
DROP CONSTRAINT terra_msgs_terra_chain_id_fkey;

ALTER TABLE terra_msgs
    ADD FOREIGN KEY (terra_chain_id) REFERENCES terra_chains(id) ON DELETE CASCADE;

--+goose Down
ALTER TABLE terra_nodes
DROP CONSTRAINT terra_nodes_terra_chain_id_fkey;

ALTER TABLE terra_nodes
    ADD FOREIGN KEY (terra_chain_id) REFERENCES terra_chains(id);

ALTER TABLE terra_msgs
DROP CONSTRAINT terra_msgs_terra_chain_id_fkey;

ALTER TABLE terra_msgs
    ADD FOREIGN KEY (terra_chain_id) REFERENCES terra_chains(id);
