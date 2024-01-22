-- +goose Up

WITH variables AS (
    SELECT
        evm_chain_id,
        CASE
            WHEN evm_chain_id = 43113 then 1 -- Avax Fuji
            WHEN evm_chain_id = 43114 then 1 -- Avax Mainnet
            WHEN evm_chain_id = 84531 THEN 200 -- Base Goerli
            WHEN evm_chain_id = 8453 THEN 200 -- Base Mainnet
            WHEN evm_chain_id = 42220 THEN 1 -- Celo Mainnet
            WHEN evm_chain_id = 44787 THEN 1 -- Celo Testnet
            WHEN evm_chain_id = 8217 THEN 1 -- Klaytn Mainnet
            WHEN evm_chain_id = 1001 THEN 1 -- Klaytn Mainnet
            WHEN evm_chain_id = 1088 THEN 1 -- Metis Mainnet
            WHEN evm_chain_id = 588 THEN 1 -- Metis Rinkeby
            WHEN evm_chain_id = 420 THEN 200 -- Optimism Goerli
            WHEN evm_chain_id = 10 THEN 200 -- Optimism Mainnet
            WHEN evm_chain_id = 137 THEN 500 -- Polygon Mainnet
            WHEN evm_chain_id = 80001 THEN 500 -- Polygon Mumbai
            WHEN evm_chain_id = 534352 THEN 1 -- Scroll Mainnet
            WHEN evm_chain_id = 534351 THEN 1 -- Scroll Sepolia
            ELSE 50 -- all other chains
            END AS finality_depth
    FROM evm.log_poller_blocks
    GROUP BY evm_chain_id
)

UPDATE evm.log_poller_blocks AS lpb
SET finalized_block_number = greatest(lpb.block_number - v.finality_depth, 0)
FROM variables v
WHERE lpb.evm_chain_id = v.evm_chain_id
  AND lpb.finalized_block_number = 0;