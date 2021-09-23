-- +goose Up

ALTER TABLE p2p_peers DROP CONSTRAINT p2p_peers_peer_id_fkey;

-- +goose Down

ALTER TABLE p2p_peers ADD CONSTRAINT p2p_peers_peer_id_fkey FOREIGN KEY (peer_id) REFERENCES encrypted_p2p_keys (peer_id);
