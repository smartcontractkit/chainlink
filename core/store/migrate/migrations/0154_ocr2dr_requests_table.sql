-- +goose Up

CREATE TABLE ocr2dr_requests(
    request_id bytea CHECK (octet_length(request_id) = 32) PRIMARY KEY,
    contract_address bytea CHECK (octet_length(contract_address) = 20) NOT NULL,
    run_id bigint, --  NOT NULL REFERENCES public.pipeline_runs(id) ON DELETE CASCADE DEFERRABLE
    received_at timestamp with time zone NOT NULL,
    request_tx_hash bytea CHECK (octet_length(request_tx_hash) = 32) NOT NULL,
    state INTEGER,
    result_ready_at timestamp with time zone,
    result bytea,
    error_type INTEGER,
    error bytea,
    transmitted_result bytea,
    transmitted_error bytea
);

CREATE INDEX idx_ocr2dr_requests ON ocr2dr_requests (contract_address);

-- +goose Down

DROP INDEX IF EXISTS idx_ocr2dr_requests;
DROP TABLE ocr2dr_requests;
