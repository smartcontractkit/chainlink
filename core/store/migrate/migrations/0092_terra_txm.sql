-- +goose Up
-- +goose StatementBegin
ALTER TABLE heads RENAME TO evm_heads;
ALTER TABLE nodes RENAME TO evm_nodes;
-- CREATE TABLE terra_chains (
--     id text PRIMARY KEY,
--     cfg jsonb NOT NULL DEFAULT '{}',
--     created_at timestamptz NOT NULL,
--     updated_at timestamptz NOT NULL
-- );
CREATE TABLE terra_nodes (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL CHECK (name != ''),
 	terra_chain_id text NOT NULL, -- REFERENCES terra_chains (id),
	tendermint_url text CHECK (tendermint_url != ''),
	fcd_url text CHECK (fcd_url != ''),
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL
);
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
-- TODO terra_chain_id text NOT NULL REFERENCES terra_chains (id),
	contract_id text NOT NULL,
    msg bytea NOT NULL,
	state text NOT NULL,
	tx_hash text, --TODO: not null for certain states
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL
);
CREATE TRIGGER notify_terra_msg_insertion AFTER INSERT ON terra_msgs FOR EACH STATEMENT EXECUTE PROCEDURE notify_terra_msg_insert();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE terra_msgs;
DROP FUNCTION notify_terra_msg_insert;
DROP TABLE terra_nodes;
ALTER TABLE evm_nodes RENAME TO nodes;
ALTER TABLE evm_heads RENAME TO heads;
-- +goose StatementEnd
