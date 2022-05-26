-- +goose Up

-- +goose StatementBegin

--
-- Name: eth_txes_state; Type: TYPE; Schema: public; Owner: postgres
--
CREATE TYPE eth_tx_attempt AS (
		 id BIGINT,
		 eth_tx_id BIGINT,
		 gas_price NUMERIC(78,0),
		 signed_raw_tx BYTEA,
		 hash BYTEA,
		 broadcast_before_block_num BIGINT,
		 state eth_tx_attempts_state,
		 created_at TIMESTAMP,
		 chain_specific_gas_limit BIGINT,
		 tx_type SMALLINT,
		 gas_tip_cap NUMERIC(78,0),
		 gas_fee_cap NUMERIC(78,0)
)
-- +goose StatementEnd

-- +goose Down
DROP TYPE IF EXISTS eth_tx_attempt;
