-- +goose Up
UPDATE eth_tx_attempts SET state = 'broadcast', broadcast_before_block_num = eth_receipts.block_number
FROM eth_receipts
WHERE eth_tx_attempts.state = 'in_progress' AND eth_tx_attempts.hash = eth_receipts.tx_hash;
