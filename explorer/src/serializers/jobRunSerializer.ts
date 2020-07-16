import {
  Serializer as JSONAPISerializer,
  SerializerOptions,
} from 'jsonapi-serializer'
import { JobRun } from '../entity/JobRun'
import { Config } from '../config'

export const BASE_ATTRIBUTES = [
  'chainlinkNode',
  'runId',
  'jobId',
  'status',
  'type',
  'requestId',
  'txHash',
  'requester',
  'error',
  'createdAt',
  'finishedAt',
]

export const chainlinkNode = {
  ref: 'id',
  attributes: ['name', 'url'],
}

export const taskRuns = {
  ref: 'id',
  attributes: [
    'jobRunId',
    'index',
    'type',
    'status',
    'confirmations',
    'minimumConfirmations',
    'error',
    'transactionHash',
    'transactionStatus',
  ],
}

const ETHERSCAN_HOST = Config.etherscanHost()

const jobRunSerializer = (run: JobRun) => {
  const opts = {
    attributes: BASE_ATTRIBUTES.concat(['taskRuns']),
    chainlinkNode,
    keyForAttribute: 'camelCase',
    taskRuns,
    meta: {
      etherscanHost: ETHERSCAN_HOST,
    },
  } as SerializerOptions

  return new JSONAPISerializer('job_runs', opts).serialize(run)
}

export default jobRunSerializer
