import { adapterTypes, status } from './constants'
import * as dbTypes from './db'

interface RunResult {
  data: { result?: string }
  error: boolean | null
  jobRunId: string
  taskRunId: string
  status: status
}

export interface BridgeType
  extends Omit<dbTypes.BridgeType, 'incomingTokenHash' | 'salt'> {
  id: string
}

//REVIEW what to do with this?
export interface ExternalInitiator
  extends Omit<
    dbTypes.ExternalInitiator,
    'salt' | 'hashedSecret' | 'accessKey'
  > {}

export interface Initiator extends dbTypes.Initiator {
  params: object
}

export interface JobRun
  extends Omit<
    dbTypes.JobRun,
    'initiatorId' | 'overridesId' | 'jobSpecId' | 'resultId'
  > {
  id: string
  initiator: Initiator
  jobId: string
  overrides: RunResult
  result: RunResult
  taskRuns: TaskRuns
  createdAt: string
  finishedAt: string
  status: string
  payment: number
}

export interface JobSpec extends dbTypes.JobSpec {
  initiators: Initiators
  tasks: TaskSpecs
  runs: TaskRuns
  errors: JobSpecErrors
}

export interface TaskRun
  extends Omit<dbTypes.TaskRun, 'taskSpecId' | 'jobRunId' | 'resultId'> {
  id: string
  result: RunResult
  task: TaskSpec
  updatedAt: Date
  type: adapterTypes
  minimumConfirmations: number
  status: string
}
export interface TaskSpec {
  id: number
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  type: adapterTypes
  confirmations?: number
  params?: object
}

export interface JobSpecError {
  id: string
  description: string
  occurrences: number
  createdAt: Date
  updatedAt: Date
}

//REVIEW Not needed?
export interface TxAttempt extends dbTypes.TxAttempt {}

export interface Transaction
  extends Omit<dbTypes.Tx, 'surrogateId' | 'signedRawTx'> {
  rawHex: string
}

export type BridgeTypes = Array<BridgeType>
export type ExternalInitiators = Array<ExternalInitiator>
export type Initiators = Array<Initiator>
export type JobRuns = Array<JobRun>
export type JobSpecs = Array<JobSpec>
export type TaskRuns = Array<TaskRun>
export type TaskSpecs = Array<TaskSpec>
export type TxAttempts = Array<TxAttempt>
export type Transactions = Array<Transaction>
export type JobSpecErrors = Array<JobSpecError>
