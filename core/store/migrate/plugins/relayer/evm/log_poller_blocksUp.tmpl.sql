CREATE TABLE {{ .Schema }}.log_poller_blocks (
--	evm_chain_id numeric(78) NOT NULL,
	block_hash bytea NOT NULL,
	block_number int8 NOT NULL,
	created_at timestamptz NOT NULL,
	block_timestamp timestamptz NOT NULL,
	finalized_block_number int8 DEFAULT 0 NOT NULL,
	CONSTRAINT block_hash_uniq UNIQUE (block_hash),
	CONSTRAINT log_poller_blocks_block_number_check CHECK ((block_number > 0)),
	CONSTRAINT log_poller_blocks_finalized_block_number_check CHECK ((finalized_block_number >= 0)),
	CONSTRAINT log_poller_blocks_pkey PRIMARY KEY (block_number)
);
CREATE INDEX idx_evm_log_poller_blocks_order_by_block ON {{ .Schema }}.log_poller_blocks USING btree (block_number DESC);

INSERT INTO {{ .Schema }}.log_poller_blocks (block_hash, block_number, created_at, block_timestamp, finalized_block_number)
SELECT block_hash, block_number, created_at, block_timestamp, finalized_block_number FROM evm.log_poller_blocks WHERE evm_chain_id = '{{ .ChainID}}';

DELETE FROM evm.log_poller_blocks WHERE evm_chain_id = '{{ .ChainID}}';
