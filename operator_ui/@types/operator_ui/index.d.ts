import * as dbTypes from "../db"

interface initiatorParams {
    type: String
    params: Object
}

interface runResult {
    data: { result: string | null }
    error: boolean | null
    jobRunId: string
    taskRunId: string
    status: "In Progress" | "Pending Confirmations" | "Pending Connection" | "Pending Bridge" | "Pending Sleep" | "Errored" | "Completed" | null
}

export interface IBridgeType extends Omit<dbTypes.bridgeType, 'incoming_token_hash' | 'salt'> {
}

//REVIEW what to do with this?
export interface IExternalInitiator extends Omit<dbTypes.externalInitiator, 'salt' | 'hashed_secret' | 'access_key'> {
}

export interface IInitiator extends dbTypes.Initiator {
}

export interface IJobRun extends Omit<dbTypes.jobRun, 'initiator_id' | 'overrides_id' | 'job_spec_id'> {
    initiator: initiatorParams
    jobId: string
    overrides: RunResult
    result: RunResult
    taskRuns: ITaskRuns
}

export interface IJobSpec extends dbTypes.jobSpec {
    initiators: initiatorParams
    tasks: ITaskRuns
}

export interface ITaskRun extends Omit<dbTypes.taskRun, 'task_spec_id' | 'job_run_id' | 'result_id'> {
    result: RunResult
    task: ITaskSpec
    updatedAt: Date
}

export interface ITaskSpec extends Omit<dbTypes.taskSpec, 'job_spec_id'> {
}

//Review Not needed?
export interface ITxAttempt extends dbTypes.txAttempt {
}

export interface ITransaction extends Omit<dbTypes.Txes, 'surrogateId' | 'signedRawTx'> {
    rawHex: string
}

export type IBridgeTypes = Array<IBridgeType>
export type IExternalInitiators = Array<IExternalInitiator>
export type IInitiators = Array<IInitiator>
export type IJobRuns = Array<IJobRun>
export type IJobSpecs = Array<IJobSpec>
export type ITaskRuns = Array<ITaskRun>
export type ITaskSpecs = Array<ITaskSpec>
export type ITxAttempts = Array<ITxAttempt>
export type ITransactions = Array<ITransaction>
