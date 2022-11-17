-- +goose Up

CREATE TYPE ocr2dr_request_state AS ENUM (
    'in_progress',
    'result_ready',
    'transmitted',
    'confirmed'
);

CREATE TYPE ocr2dr_error_type AS ENUM {
    'none',
    'node_exception',
    'sandbox_timeout',
    'user_exception'
};

CREATE TABLE ocr2dr_requests(
    id BIGSERIAL PRIMARY KEY,
    request_id bytea CHECK (octet_length(request_id) = 32) NOT NULL UNIQUE,
    oracle_address bytea CHECK (octet_length(oracle_address) = 20) NOT NULL,
    run_id bigint NOT NULL REFERENCES public.pipeline_runs(id) ON DELETE CASCADE DEFERRABLE,
    received_at timestamp with time zone NOT NULL,
    request_tx_hash bytea CHECK (octet_length(oracle_address) = 32) NOT NULL,
    state ocr2dr_request_state,
    result_ready_at timestamp with time zone,
    result bytea,
    error_type ocr2dr_error_type,
    error bytea,
    is_ocr_participant boolean,
    transmitted_result bytea,
    transmitted_error bytea,
    on_chain_result bytea,
    on_chain_error bytea
);

-- +goose Down

DROP TABLE ocr2dr_requests;
DROP TYPE ocr2dr_error_type;
DROP TYPE ocr2dr_request_state;
