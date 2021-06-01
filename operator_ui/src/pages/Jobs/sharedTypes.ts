import { ApiResponse, PaginatedApiResponse } from 'utils/json-api-client'
import {
  Initiator,
  JobRun,
  JobSpec,
  JobSpecError,
  OcrJobRun,
  JobSpecV2,
  RunResult,
  RunStatus,
  TaskRun,
  TaskSpec,
} from 'core/store/models'
import * as time from 'time'

export type JobRunsResponse =
  | PaginatedApiResponse<JobRun[]>
  | PaginatedApiResponse<OcrJobRun[]>

export type JobSpecResponse = ApiResponse<JobSpec> | ApiResponse<JobSpecV2>

export type BaseJob = {
  createdAt: string
  definition: string
  errors: JobSpecError[]
  id: string
  name: string | null
}

export type JobSpecType =
  | 'directrequest'
  | 'fluxmonitor'
  | 'offchainreporting'
  | 'keeper'
  | 'cron'
  | 'webhook'
  | 'vrf'

export type JobV2 = BaseJob & {
  dotDagSource: string
  type: 'v2'
  specType: JobSpecType
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
}

export type DirectRequestJobRun = BaseJobRun & {
  initiator: Initiator
  overrides: RunResult
  result: RunResult
  taskRuns: TaskRun[]
  payment: string | null
  status: RunStatus
  type: 'Direct request job run'
}

export type PipelineJobRunStatus = 'in_progress' | 'errored' | 'completed'
export type PipelineTaskRunStatus =
  | 'in_progress'
  | 'errored'
  | 'completed'
  | 'not_run'

export type PipelineTaskRun = OcrJobRun['taskRuns'][0] & {
  status: PipelineTaskRunStatus
}

export type PipelineJobRun = BaseJobRun & {
  outputs: null | (string | null)[]
  errors: null | (string | null)[]
  pipelineSpec: {
    dotDagSource: string
  }
  status: PipelineJobRunStatus
  taskRuns: PipelineTaskRun[]
  type: 'Pipeline job run'
}

export type JobData = {
  job?: DirectRequestJob | JobV2
  jobSpec?: JobSpecResponse['data']
  recentRuns?: PipelineJobRun[] | DirectRequestJobRun[]
  recentRunsCount: number
  externalJobID?: string
}
