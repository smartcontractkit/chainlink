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
  createdAt: time.Time
  finishedAt: time.Time | null
  id: string
  jobId: string
  status: RunStatus
}

export type DirectRequestJobRun = BaseJobRun & {
  initiator: Initiator
  overrides: RunResult
  result: RunResult
  taskRuns: TaskRun[]
  payment: string | null
  type: 'Direct request job run'
}

export type OcrTaskRun = OcrJobRun['taskRuns'][0] & {
  status: 'in_progress' | 'errored' | 'completed' | 'aborted'
}

export type OffChainReportingJobRun = BaseJobRun & {
  outputs: (string | null)[]
  errors: (string | null)[]
  pipelineSpec: {
    DotDagSource: string
  }
  taskRuns: OcrTaskRun[]
  type: 'Off-chain reporting job run'
}

export type JobData = {
  job?: DirectRequestJob | OffChainReportingJob
  jobSpec?: JobSpecResponse['data']
  recentRuns?: BaseJobRun[]
  recentRunsCount: number
}
