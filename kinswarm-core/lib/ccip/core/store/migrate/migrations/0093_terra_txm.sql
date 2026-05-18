-- +goose Up
-- +goose StatementBegin
ALTER TABLE heads RENAME TO evm_heads;
ALTER TABLE nodes RENAME TO evm_nodes;
CREATE TABLE terra_chains (
    id text PRIMARY KEY,
    cfg jsonb NOT NULL DEFAULT '{}',
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    enabled BOOL DEFAULT TRUE NOT NULL
);
CREATE TABLE terra_nodes (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL CHECK (name != ''),
    terra_chain_id text NOT NULL REFERENCES terra_chains (id),
    tendermint_url text CHECK (tendermint_url != ''),
    fcd_url text CHECK (fcd_url != ''),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE INDEX idx_nodes_terra_chain_id ON terra_nodes (terra_chain_id);
CREATE UNIQUE INDEX idx_terra_nodes_unique_name ON terra_nodes (lower(name));
CREATE FUNCTION notify_terra_msg_insert() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM pg_notify('insert_on_terra_msg'::text, NOW()::text);
    RETURN NULL;
END
$$;
CREATE TABLE terra_msgs (
    id BIGSERIAL PRIMARY KEY,
    terra_chain_id text NOT NULL REFERENCES terra_chains (id),
    contract_id text NOT NULL,
    raw bytea NOT NULL,
    state text NOT NULL,
    tx_hash text,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CHECK (tx_hash<>null OR (state<>'broadcasted' AND state<>'confirmed'))
);
CREATE TRIGGER notify_terra_msg_insertion AFTER INSERT ON terra_msgs FOR EACH STATEMENT EXECUTE PROCEDURE notify_terra_msg_insert();
CREATE INDEX idx_terra_msgs_terra_chain_id_state ON terra_msgs (terra_chain_id, state);

CREATE FUNCTION check_terra_msg_state_transition() RETURNS TRIGGER AS $$
DECLARE
  state_transition_map jsonb := json_build_object(
        'unstarted', json_build_object('errored', true, 'broadcasted', true),
        'broadcasted', json_build_object('errored', true, 'confirmed', true));
BEGIN
    IF NOT state_transition_map ? OLD.state THEN
        RAISE EXCEPTION 'Invalid from state %. Valid from states %', OLD.state, state_transition_map;
    END IF;
    IF NOT state_transition_map->OLD.state ? NEW.state THEN
        RAISE EXCEPTION 'Invalid state transition from % to %. Valid to states %', OLD.state, NEW.state, state_transition_map->OLD.state;
    END IF;
  RETURN NEW;
END
$$ LANGUAGE plpgsql;
CREATE TRIGGER validate_state_update BEFORE UPDATE ON terra_msgs
    FOR EACH ROW EXECUTE PROCEDURE check_terra_msg_state_transition();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE terra_msgs;
DROP FUNCTION notify_terra_msg_insert;
DROP FUNCTION check_terra_msg_state_transition;
DROP TABLE terra_nodes;
DROP TABLE terra_chains;
ALTER TABLE evm_nodes RENAME TO nodes;
ALTER TABLE evm_heads RENAME TO heads;
-- +goose StatementEnd
