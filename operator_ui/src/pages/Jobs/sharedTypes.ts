import { ApiResponse, PaginatedApiResponse } from '@chainlink/json-api-client'
import {
  Initiator,
  JobRun,
  JobSpec,
  JobSpecError,
  OcrJobRun,
  OcrJobSpec,
  RunResult,
  RunStatus,
  TaskRun,
  TaskSpec,
} from 'core/store/models'
import * as time from 'time'

export type JobRunsResponse =
  | PaginatedApiResponse<JobRun[]>
  | PaginatedApiResponse<OcrJobRun[]>

export type JobSpecResponse = ApiResponse<JobSpec> | ApiResponse<OcrJobSpec>

export type BaseJob = {
  createdAt: string
  definition: string
  errors: JobSpecError[]
  id: string
  name?: string
}

export type OffChainReportingJob = BaseJob & {
  dotDagSource: string
  type: 'Off-chain reporting'
}

export type DirectRequestJob = BaseJob & {
  earnings: number | null
  endAt: string | null
  initiators: Initiator[]
  minPayment?: string | null
  startAt: string | null
  tasks: TaskSpec[]
  type: 'Direct request'
}

export type BaseJobRun = {
  createdAt: string
  id: string
  status: RunStatus
  jobId: string
}

export type DirectRequestJobRun = BaseJobRun & {
  initiator: Initiator
  overrides: RunResult
  result: RunResult
  taskRuns: TaskRun[]
  finishedAt: time.Time | null
  payment: string | null
}

export type JobData = {
  job?: DirectRequestJob | OffChainReportingJob
  jobSpec?: JobSpecResponse['data']
  recentRuns?: BaseJobRun[]
  recentRunsCount: number
}
