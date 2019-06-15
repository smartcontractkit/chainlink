export interface bridge_types {
    name: string;
    url: string | null;
    confirmations: number | null;
    incoming_token_hash: string | null;
    salt: string | null;
    outgoing_token: string | null;
    minimum_contract_payment: string | null;
}

export interface external_initiators {
    id: number;
    created_at: Date | null;
    updated_at: Date | null;
    deleted_at: Date | null;
    access_key: string | null;
    salt: string | null;
    hashed_secret: string | null;
}

export interface initiators {
    id: number;
    job_spec_id: string | null;
    type: string;
    created_at: Date | null;
    schedule: string | null;
    time: Date | null;
    ran: boolean | null;
    address: string | null;
    requesters: string | null;
    deleted_at: Date | null;
}

export interface job_runs {
    id: string;
    job_spec_id: string;
    result_id: number | null;
    run_request_id: number | null;
    status: string | null;
    created_at: Date | null;
    finished_at: Date | null;
    updated_at: Date | null;
    initiator_id: number | null;
    creation_height: string | null;
    observed_height: string | null;
    overrides_id: number | null;
    deleted_at: Date | null;
}

export interface job_specs {
    id: string;
    created_at: Date | null;
    start_at: Date | null;
    end_at: Date | null;
    deleted_at: Date | null;
}

export interface run_requests {
    id: number;
    request_id: string | null;
    tx_hash: string | null;
    requester: string | null;
    created_at: Date | null;
}

export interface run_results {
    id: number;
    cached_job_run_id: string | null;
    cached_task_run_id: string | null;
    data: string | null;
    status: string | null;
    error_message: string | null;
    amount: string | null;
}

export interface task_runs {
    id: string;
    job_run_id: string;
    result_id: number | null;
    status: string | null;
    task_spec_id: number | null;
    minimum_confirmations: number | null;
    created_at: Date | null;
    confirmations: number;
}

export interface task_specs {
    id: number;
    created_at: Date | null;
    updated_at: Date | null;
    deleted_at: Date | null;
    job_spec_id: string | null;
    type: string;
    confirmations: number | null;
    params: string | null;
}

export interface tx_attempts {
    id: number;
    tx_id: number | null;
    created_at: Date;
    hash: string;
    gas_price: string;
    confirmed: boolean;
    sent_at: number;
    signed_raw_tx: string;
}


export interface txes {
    id: number;
    surrogate_id: string | null;
    from: string;
    to: string;
    data: string;
    nonce: number;
    value: string;
    gas_limit: number;
    hash: string;
    gas_price: string;
    confirmed: boolean;
    sent_at: number;
    signed_raw_tx: string;
}
