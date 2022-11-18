-- +goose Up

CREATE TABLE ocr2dr_requests(
    id BIGSERIAL PRIMARY KEY,
    request_id bytea CHECK (octet_length(request_id) = 32) NOT NULL UNIQUE,
    oracle_address bytea CHECK (octet_length(oracle_address) = 20) NOT NULL,
    run_id bigint NOT NULL REFERENCES public.pipeline_runs(id) ON DELETE CASCADE DEFERRABLE,
    received_at timestamp with time zone NOT NULL,
    request_tx_hash bytea CHECK (octet_length(oracle_address) = 32) NOT NULL,
    state INTEGER,
    result_ready_at timestamp with time zone,
    result bytea,
    error_type INTEGER,
    error bytea,
    transmitted_result bytea,
    transmitted_error bytea
);

-- +goose Down

DROP TABLE ocr2dr_requests;
