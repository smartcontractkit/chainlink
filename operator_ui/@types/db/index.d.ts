export interface bridgeType {
    name: string;
    url: string | null;
    confirmations: number | null;
    incoming_token_hash: string | null;
    salt: string | null;
    outgoing_token: string | null;
    minimum_contract_payment: string | null;
}

export interface externalInitiator {
    id: number;
    createdAt: Date | null;
    updatedAt: Date | null;
    deletedAt: Date | null;
    access_key: string | null;
    salt: string | null;
    hashed_secret: string | null;
}

export interface Initiator {
    id: number;
    job_spec_id: string | null;
    type: string;
    createdAt: Date | null;
    schedule: string | null;
    time: Date | null;
    ran: boolean | null;
    address: string | null;
    requesters: string | null;
    deletedAt: Date | null;
}

export interface jobRun {
    id: string;
    job_spec_id: string;
    result_id: number | null;
    run_request_id: number | null;
    status: "In Progress" | "Pending Confirmations" | "Pending Connection" | "Pending Bridge" | "Pending Sleep" | "Errored" | "Completed" | null
    createdAt: Date | null;
    finishedAt: Date | null;
    updatedAt: Date | null;
    initiator_id: number | null;
    creationHeight: string | null;
    observedHeight: string | null;
    overrides_id: number | null;
    deletedAt: Date | null;
}

export interface jobSpec {
    id: string;
    createdAt: Date | null;
    startAt: Date | null;
    endAt: Date | null;
    deletedAt: Date | null;
}

export interface runRequests {
    id: number;
    request_id: string | null;
    tx_hash: string | null;
    requester: string | null;
    createdAt: Date | null;
}

export interface taskRun {
    id: string;
    job_run_id: string;
    result_id: number | null;
    status: "In Progress" | "Pending Confirmations" | "Pending Connection" | "Pending Bridge" | "Pending Sleep" | "Errored" | "Completed" | null
    task_spec_id: number | null;
    minimumConfirmations: number | null;
    createdAt: Date | null;
    confirmations: number;
}

export interface taskSpec {
    id: number;
    createdAt: Date | null;
    updatedAt: Date | null;
    deletedAt: Date | null;
    job_spec_id: string | null;
    type: string;
    confirmations: number | null;
    params: string | null;
}

export interface txAttempt {
    id: number;
    tx_id: number | null;
    createdAt: Date;
    hash: string;
    gas_price: string;
    confirmed: boolean;
    sent_at: number;
    signed_raw_tx: string;
}


export interface Tx {
    id: number;
    surrogateId: string | null;
    from: string;
    to: string;
    data: string;
    nonce: number;
    value: string;
    gasLimit: number;
    hash: string;
    gasPrice: string;
    confirmed: boolean;
    sentAt: number;
    signedRawTx: string;
}
