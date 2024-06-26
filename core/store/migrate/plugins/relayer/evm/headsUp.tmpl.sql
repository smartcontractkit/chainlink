-- Do nothing for `evm` schema for backward compatibility
{{ if ne .Schema "evm"}}
CREATE TABLE  {{ .Schema }}.heads (
	id bigserial NOT NULL,
	hash bytea NOT NULL,
	"number" int8 NOT NULL,
	parent_hash bytea NOT NULL,
	created_at timestamptz NOT NULL,
	"timestamp" timestamptz NOT NULL,
	l1_block_number int8 NULL,
--	evm_chain_id numeric(78) NOT NULL,
	base_fee_per_gas numeric(78) NULL,
	CONSTRAINT chk_hash_size CHECK ((octet_length(hash) = 32)),
	CONSTRAINT chk_parent_hash_size CHECK ((octet_length(parent_hash) = 32)),
	CONSTRAINT heads_pkey1 PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_heads_hash ON {{ .Schema }}.heads USING btree (hash);
CREATE INDEX idx_heads_number ON {{ .Schema }}.heads USING btree ("number");

INSERT INTO {{ .Schema }}.heads (hash, "number", parent_hash, created_at, "timestamp", l1_block_number, base_fee_per_gas)
SELECT hash, "number", parent_hash, created_at, "timestamp", l1_block_number,  base_fee_per_gas
FROM evm.heads WHERE evm_chain_id = '{{ .ChainID }}';

DELETE FROM evm.heads WHERE evm_chain_id = '{{ .ChainID }}';
{{ end}}