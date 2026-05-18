-- +goose Up
-- Drop nodes tables
DROP TABLE evm_nodes;
DROP TABLE solana_nodes;
DROP TABLE starknet_nodes;


-- +goose Down
-- evm_nodes definition
CREATE TABLE evm_nodes (
	id serial NOT NULL,
	"name" varchar(255) NOT NULL,
	evm_chain_id numeric(78) NOT NULL,
	ws_url text NULL,
	http_url text NULL,
	send_only bool NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT nodes_http_url_check CHECK ((http_url <> ''::text)),
	CONSTRAINT nodes_name_check CHECK (((name)::text <> ''::text)),
	CONSTRAINT nodes_pkey PRIMARY KEY (id),
	CONSTRAINT nodes_ws_url_check CHECK ((ws_url <> ''::text)),
	CONSTRAINT primary_or_sendonly CHECK (((send_only AND (ws_url IS NULL) AND (http_url IS NOT NULL)) OR ((NOT send_only) AND (ws_url IS NOT NULL))))
);
CREATE INDEX idx_nodes_evm_chain_id ON evm_nodes USING btree (evm_chain_id);
CREATE UNIQUE INDEX idx_nodes_unique_name ON evm_nodes USING btree (lower((name)::text));
CREATE UNIQUE INDEX idx_unique_http_url ON evm_nodes USING btree (http_url);
CREATE UNIQUE INDEX idx_unique_ws_url ON evm_nodes USING btree (ws_url);
-- evm_nodes foreign keys.
ALTER TABLE evm_nodes ADD CONSTRAINT nodes_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE;

-- solana_nodes definition
CREATE TABLE solana_nodes (
	id serial NOT NULL,
	"name" varchar(255) NOT NULL,
	solana_chain_id text NOT NULL,
	solana_url text NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT solana_nodes_name_check CHECK (((name)::text <> ''::text)),
	CONSTRAINT solana_nodes_pkey PRIMARY KEY (id),
	CONSTRAINT solana_nodes_solana_url_check CHECK ((solana_url <> ''::text))
);
CREATE INDEX idx_nodes_solana_chain_id ON solana_nodes USING btree (solana_chain_id);
CREATE UNIQUE INDEX idx_solana_nodes_unique_name ON solana_nodes USING btree (lower((name)::text));
-- solana_nodes foreign keys
ALTER TABLE solana_nodes ADD CONSTRAINT solana_nodes_solana_chain_id_fkey FOREIGN KEY (solana_chain_id) REFERENCES solana_chains(id) ON DELETE CASCADE;

-- starknet_nodes definition
CREATE TABLE starknet_nodes (
	id serial NOT NULL,
	"name" varchar(255) NOT NULL,
	starknet_chain_id text NOT NULL,
	url text NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT starknet_nodes_name_check CHECK (((name)::text <> ''::text)),
	CONSTRAINT starknet_nodes_pkey PRIMARY KEY (id),
	CONSTRAINT starknet_nodes_url_check CHECK ((url <> ''::text))
);
CREATE INDEX idx_starknet_nodes_chain_id ON starknet_nodes USING btree (starknet_chain_id);
CREATE UNIQUE INDEX idx_starknet_nodes_unique_name ON starknet_nodes USING btree (lower((name)::text));
-- starknet_nodes foreign keys
ALTER TABLE starknet_nodes ADD CONSTRAINT starknet_nodes_chain_id_fkey FOREIGN KEY (starknet_chain_id) REFERENCES starknet_chains(id) ON DELETE CASCADE;
