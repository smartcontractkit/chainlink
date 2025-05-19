-- +goose Up

-- see 0154_ocr2dr_requests_table.sql for initial definition
ALTER TABLE ocr2dr_requests RENAME TO functions_requests;
ALTER INDEX idx_ocr2dr_requests RENAME TO idx_functions_requests;

ALTER TABLE functions_requests DROP COLUMN run_id;

ALTER TABLE functions_requests
  ADD COLUMN flags bytea,
  ADD COLUMN aggregation_method INTEGER,
  ADD COLUMN callback_gas_limit INTEGER,
  ADD COLUMN coordinator_contract_address bytea CHECK (octet_length(coordinator_contract_address) = 20),
  ADD COLUMN onchain_metadata bytea,
  ADD COLUMN processing_metadata bytea;

-- +goose Down

ALTER TABLE functions_requests
  DROP COLUMN flags,
  DROP COLUMN aggregation_method,
  DROP COLUMN callback_gas_limit,
  DROP COLUMN coordinator_contract_address,
  DROP COLUMN onchain_metadata,
  DROP COLUMN processing_metadata;

ALTER TABLE functions_requests ADD COLUMN run_id bigint;

ALTER INDEX idx_functions_requests RENAME TO idx_ocr2dr_requests;
ALTER TABLE functions_requests RENAME TO ocr2dr_requests;
