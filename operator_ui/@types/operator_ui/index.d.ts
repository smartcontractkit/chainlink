import * as dbTypes from "../db"
import { status, adapterTypes } from "../constants"

interface runResult {
    data: { result: string | null }
    error: boolean | null
    jobRunId: string
    taskRunId: string
    status: status
}

export interface IBridgeType extends Omit<dbTypes.bridgeType, 'incomingTokenHash' | 'salt'> { }

//REVIEW what to do with this?
export interface IExternalInitiator extends Omit<dbTypes.externalInitiator, 'salt' | 'hashedSecret' | 'accessKey'> { }

export interface IInitiator extends dbTypes.Initiator { }

export interface IJobRun extends Omit<dbTypes.jobRun, 'initiatorId' | 'overridesId' | 'jobSpecId'> {
    initiator: dbTypes.Initiator
    jobId: string
    overrides: runResult
    result: runResult
    taskRuns: ITaskRuns
}

export interface IJobSpec extends dbTypes.jobSpec {
    initiators: dbTypes.Initiator
    tasks: ITaskSpecs
    runs: ITaskRuns
}

export interface ITaskRun extends Omit<dbTypes.taskRun, 'taskSpecId' | 'jobRunId' | 'resultId'> {
    result: runResult
    task: ITaskSpec
    updatedAt: Date
    type: adapterTypes
}

export interface ITaskSpec extends Omit<dbTypes.taskSpec, 'jobSpecId'> { }

//REVIEW Not needed?
export interface ITxAttempt extends dbTypes.txAttempt { }

export interface ITransaction extends Omit<dbTypes.Tx, 'surrogateId' | 'signedRawTx'> {
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
