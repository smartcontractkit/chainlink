import * as dbTypes from '../db'
import { status, adapterTypes, initiatorTypes } from '../constants'

interface RunResult {
  data: { result?: string }
  error?: boolean
  jobRunId: string
  taskRunId: string
  status: status
  amount?: number
}

export interface BridgeType {
  id: string
  name: string
  url: string
  confirmations: number
  outgoingToken: string
  minimumContractPayment?: string
}

export interface ExternalInitiator {
  id: number
  createdAt: string
  updatedAt?: string
  deletedAt?: string
}

export interface Initiator {
  params: object
  id: number
  jobSpecId: string
  type: initiatorTypes
  createdAt: string
  schedule?: string
  time?: string
  ran?: boolean
  address: string
  requesters?: string
  deletedAt?: string
}

export interface JobRun {
  id: string
  initiator: Initiator
  jobId: string
  overrides: RunResult
  result: RunResult
  taskRuns: ITaskRuns
  createdAt: string
  finishedAt: string
  updatedAt: string
  status: status
  creationHeight: string
  observedHeight?: string
  deletedAt?: string
}

export interface JobSpec {
  initiators: IInitiators
  tasks: TaskSpecs
  runs: ITaskRuns
  id: string
  createdAt: string
  startAt?: string
  endAt?: string
  deletedAt?: string
}

export interface TaskRun {
  id: string
  result: RunResult
  status: status
  task: TaskSpec
  type: adapterTypes
  createdAt: string
  updatedAt?: string
  confirmations: number
  minimumConfirmations?: number
}
export interface TaskSpec {
  id: number
  createdAt: string
  updatedAt?: string
  deletedAt?: string
  type: adapterTypes
  confirmations?: number
  params?: object
  jobSpecId: string
}

export interface TxAttempt extends dbTypes.TxAttempt {
  id: number
  txId?: number
  createdAt: string
  hash: string
  gasPrice: string
  confirmed: boolean
  sentAt: number
  signedRawTx: string
}

export interface Transaction {
  rawHex: string
  id: number
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
