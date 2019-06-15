import { status, adapterTypes, initiatorTypes } from "../constants"

export interface bridgeType {
    name: string
    url: string | null
    confirmations: number | null
    incomingTokenHash: string | null
    salt: string | null
    outgoingToken: string | null
    minimumContractPayment: string | null
}

export interface externalInitiator {
    id: number
    createdAt: Date | null
    updatedAt: Date | null
    deletedAt: Date | null
    accessKey: string | null
    salt: string | null
    hashedSecret: string | null
}

export interface Initiator {
    id: number
    jobSpecId: string | null
    type: initiatorTypes
    createdAt: Date | null
    schedule: string | null
    time: Date | null
    ran: boolean | null
    address: string | null
    requesters: string | null
    deletedAt: Date | null
}

export interface jobRun {
    id: string
    jobSpecId: string
    resultId: number | null
    runRequestId: number | null
    status: status
    createdAt: Date | null
    finishedAt: Date | null
    updatedAt: Date | null
    initiatorId: number | null
    creationHeight: string | null
    observedHeight: string | null
    overridesId: number | null
    deletedAt: Date | null
}

export interface jobSpec {
    id: string
    createdAt: Date | null
    startAt: Date | null
    endAt: Date | null
    deletedAt: Date | null
}

export interface runRequests {
    id: number
    requestId: string | null
    txHash: string | null
    requester: string | null
    createdAt: Date | null
}

export interface taskRun {
    id: string
    jobRunId: string
    resultId: number | null
    status: status
    taskSpecId: number | null
    minimumConfirmations: number | null
    createdAt: Date | null
    confirmations: number
}

export interface taskSpec {
    id: number
    createdAt: Date | null
    updatedAt: Date | null
    deletedAt: Date | null
    jobSpecId: string | null
    type: adapterTypes
    confirmations: number | null
    params: string | null //REVIEW add conditional params for known adapter types that are listed 2 lines above?
}

export interface txAttempt {
    id: number
    txId: number | null
    createdAt: Date
    hash: string
    gasPrice: string
    confirmed: boolean
    sentAt: number
    signedRawTx: string
}


export interface Tx {
    id: number
    surrogateId: string | null
    from: string
    to: string
    data: string
    nonce: number
    value: string
    gasLimit: number
    hash: string
    gasPrice: string
    confirmed: boolean
    sentAt: number
    signedRawTx: string
}
