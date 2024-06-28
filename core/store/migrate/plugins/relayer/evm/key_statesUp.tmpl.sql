CREATE TABLE {{ .Schema }}.key_states (
	id serial4 NOT NULL,
	address bytea NOT NULL,
	disabled bool DEFAULT false NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
--	evm_chain_id numeric(78) NOT NULL,
	CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20)),
	CONSTRAINT eth_key_states_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_evm_key_states_address ON {{ .Schema }}.key_states USING btree (address);
--CREATE UNIQUE INDEX idx_evm_key_states_evm_chain_id_address ON {{ .Schema }}.key_states USING btree (evm_chain_id, address);

INSERT INTO {{ .Schema }}.key_states (address, disabled, created_at, updated_at)
SELECT address, disabled, created_at, updated_at FROM evm.key_states WHERE evm_chain_id = '{{ .ChainID }}';

DELETE FROM evm.key_states WHERE evm_chain_id = '{{ .ChainID }}';
