-- +goose Up
-- +goose StatementBegin
ALTER TABLE eth_txes
ADD COLUMN IF NOT EXISTS transmit_checker jsonb DEFAULT NULL;

UPDATE eth_txes
SET transmit_checker = '{"CheckerType": "simulate"}'::jsonb
WHERE simulate;

ALTER TABLE eth_txes DROP COLUMN simulate;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE eth_txes
ADD COLUMN IF NOT EXISTS simulate bool NOT NULL DEFAULT FALSE;

UPDATE eth_txes
SET simulate = true
WHERE transmit_checker::jsonb->>'CheckerType' = 'simulate';

ALTER TABLE eth_txes
DROP COLUMN transmit_checker;
-- +goose StatementEnd
