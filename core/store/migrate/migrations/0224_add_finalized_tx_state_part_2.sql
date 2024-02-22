-- +goose Up

-- NOTE: 0221 & 0222 supposed to be single migration, but new enum values must be committed before they can be used. See notes section https://www.postgresql.org/docs/16/sql-altertype.html

-- drop constraints that we might have added on migrate down
ALTER TABLE evm.txes DROP CONSTRAINT IF EXISTS eth_txes_state_finalized_removed;
ALTER TABLE evm.tx_attempts DROP CONSTRAINT IF EXISTS eth_tx_attempts_state_finalized_removed;

-- allow tx_attempts to use finalized state
ALTER TABLE evm.tx_attempts DROP CONSTRAINT chk_eth_tx_attempts_fsm;
ALTER TABLE evm.tx_attempts ADD CONSTRAINT chk_eth_tx_attempts_fsm CHECK
    (
        ((state = ANY (ARRAY['in_progress'::public.eth_tx_attempts_state, 'insufficient_eth'::public.eth_tx_attempts_state]))
            AND (broadcast_before_block_num IS NULL))
        OR (state = ANY (ARRAY['broadcast'::public.eth_tx_attempts_state, 'finalized'::public.eth_tx_attempts_state]))
    ) NOT VALID;

-- allow txes to use finalized state
ALTER TABLE evm.txes DROP CONSTRAINT chk_eth_txes_fsm;
ALTER TABLE evm.txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
                state = 'unstarted'::eth_txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
        OR
                state = 'in_progress'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
        OR
                state = 'fatal_error'::eth_txes_state AND nonce IS NULL AND error IS NOT NULL
        OR
                state = 'unconfirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
        OR
                state = 'confirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
        OR
                state = 'confirmed_missing_receipt'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
        OR
                state = 'finalized'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    ) NOT VALID;

-- if the ReaperThreshold config value is high, we might endup with large number of 'confirmed' transactions.
-- Lets mark old transactions as finalized to prevent unnecessary RPC calls
UPDATE evm.txes set state = 'finalized'::eth_txes_state WHERE state = 'confirmed'::eth_txes_state AND created_at < (now() - interval '240 hours');

-- Mark the most recent attempts as finalized for old transactions to ensure that we do not violate an invariant that each finalized tx has a single finalized attempt.
WITH most_recent_attempts as (
    SELECT MAX(evm.tx_attempts.id) as id
    FROM evm.tx_attempts JOIN evm.txes ON evm.tx_attempts.eth_tx_id = evm.txes.id
    WHERE evm.txes.state = 'finalized'::eth_txes_state
    GROUP BY evm.tx_attempts.eth_tx_id
)
UPDATE evm.tx_attempts set state = 'finalized'::eth_tx_attempts_state FROM most_recent_attempts  WHERE evm.tx_attempts.id = most_recent_attempts.id;

-- +goose Down

-- it's not possible to remove label from the enum. The only option is to restrict it's usage;
UPDATE evm.txes set state = 'confirmed'::eth_txes_state WHERE state = 'finalized'::eth_txes_state;
ALTER TABLE evm.txes ADD CONSTRAINT eth_txes_state_finalized_removed CHECK (state <> 'finalized'::eth_txes_state);

UPDATE evm.tx_attempts set state = 'broadcast'::eth_tx_attempts_state WHERE state = 'finalized'::eth_tx_attempts_state;
ALTER TABLE evm.tx_attempts ADD CONSTRAINT eth_tx_attempts_state_finalized_removed CHECK (state <> 'finalized'::eth_tx_attempts_state);

-- rollback changes to the constraints
ALTER TABLE evm.tx_attempts DROP CONSTRAINT chk_eth_tx_attempts_fsm;
ALTER TABLE evm.tx_attempts ADD CONSTRAINT chk_eth_tx_attempts_fsm CHECK
    ((
        ((state = ANY (ARRAY['in_progress'::public.eth_tx_attempts_state, 'insufficient_eth'::public.eth_tx_attempts_state]))
            AND (broadcast_before_block_num IS NULL))
        OR (state = 'broadcast'::public.eth_tx_attempts_state)
    )) NOT VALID;

ALTER TABLE evm.txes DROP CONSTRAINT chk_eth_txes_fsm;
ALTER TABLE evm.txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
                state = 'unstarted'::eth_txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
        OR
                state = 'in_progress'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
        OR
                state = 'fatal_error'::eth_txes_state AND nonce IS NULL AND error IS NOT NULL
        OR
                state = 'unconfirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
        OR
                state = 'confirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
        OR
                state = 'confirmed_missing_receipt'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    ) NOT VALID;