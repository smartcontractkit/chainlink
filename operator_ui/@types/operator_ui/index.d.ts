import * as dbTypes from '../db'
import { status, adapterTypes, initiatorTypes } from '../constants'

interface RunResult {
  data: { result: string | null }
  error: boolean | null
  jobRunId: string
  taskRunId: string
  status: status
  amount: number | null
}

export interface IBridgeType {
  id: string
  name: string
  url: string
  confirmations: number
  outgoingToken: string
  minimumContractPayment: string | null
}

export interface IExternalInitiator {
  id: number
  createdAt: string
  updatedAt: string | null
  deletedAt: string | null
}

export interface IInitiator {
  params: object
  id: number
  jobSpecId: string
  type: initiatorTypes
  createdAt: string
  schedule: string | null
  time: string | null
  ran: boolean | null
  address: string
  requesters: string | null
  deletedAt: string | null
}

export interface IJobRun {
  id: string
  initiator: IInitiator
  jobId: string
  overrides: RunResult
  result: RunResult
  taskRuns: ITaskRuns
  createdAt: string
  finishedAt: string
  updatedAt: string
  status: string
  creationHeight: string
  observedHeight: string | null
  deletedAt: string | null
}

export interface IJobSpec {
  initiators: IInitiators
  tasks: ITaskSpecs
  runs: ITaskRuns
  id: string
  createdAt: string
  startAt: string | null
  endAt: string | null
  deletedAt: string | null
}

export interface ITaskRun {
  id: string
  result: RunResult
  status: status
  task: ITaskSpec
  type: adapterTypes
  createdAt: string
  updatedAt: string | null
  confirmations: number
  minimumConfirmations: number | null
}
export interface ITaskSpec {
  id: number
  createdAt: string
  updatedAt: string | null
  deletedAt: string | null
  type: adapterTypes
  confirmations: number | null
  params: object | null
  jobSpecId: string
}

export interface ITxAttempt extends dbTypes.TxAttempt {
  id: number
  txId: number | null
  createdAt: string
  hash: string
  gasPrice: string
  confirmed: boolean
  sentAt: number
  signedRawTx: string
}

export interface ITransaction {
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

export type IBridgeTypes = Array<IBridgeType>
export type IExternalInitiators = Array<IExternalInitiator>
export type IInitiators = Array<IInitiator>
export type IJobRuns = Array<IJobRun>
export type IJobSpecs = Array<IJobSpec>
export type ITaskRuns = Array<ITaskRun>
export type ITaskSpecs = Array<ITaskSpec>
export type ITxAttempts = Array<ITxAttempt>
export type ITransactions = Array<ITransaction>
