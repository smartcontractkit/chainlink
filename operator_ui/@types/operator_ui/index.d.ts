import * as dbTypes from '../db'
import { status, adapterTypes } from '../constants'

interface RunResult {
    data: { result: string | null }
    error: boolean | null
    jobRunId: string
    taskRunId: string
    status: status
    amount?: number
}

export interface IBridgeType
  extends Omit<dbTypes.BridgeType, 'incomingTokenHash' | 'salt'> {
  id: string
}

//REVIEW what to do with this?
export interface IExternalInitiator
  extends Omit<
    dbTypes.ExternalInitiator,
    'salt' | 'hashedSecret' | 'accessKey'
  > {}

export interface IInitiator extends dbTypes.Initiator {
  params: object
}

export interface IJobRun
  extends Omit<
    dbTypes.JobRun,
    'initiatorId' | 'overridesId' | 'jobSpecId' | 'resultId'
  > {
  id: string
  initiator: IInitiator
  jobId: string
  overrides: RunResult
  result: RunResult
  taskRuns: ITaskRuns
  createdAt: string
  finishedAt: string
  status: string
}

export interface IJobSpec extends dbTypes.JobSpec {
  initiators: IInitiators
  tasks: ITaskSpecs
  runs: ITaskRuns
}

export interface ITaskRun
  extends Omit<dbTypes.TaskRun, 'taskSpecId' | 'jobRunId' | 'resultId'> {
  id: string
  result: RunResult
  task: ITaskSpec
  updatedAt: Date
  type: adapterTypes
  status: string
  minimumConfirmations: number
}

// export interface ITaskSpec extends Omit<dbTypes.TaskSpec, 'jobSpecId'> {}
export interface ITaskSpec {
  id: number
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  type: adapterTypes
  confirmations?: number
  params?: object
}

//REVIEW Not needed?
export interface ITxAttempt extends dbTypes.TxAttempt {}

export interface ITransaction
  extends Omit<dbTypes.Tx, 'surrogateId' | 'signedRawTx'> {
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
