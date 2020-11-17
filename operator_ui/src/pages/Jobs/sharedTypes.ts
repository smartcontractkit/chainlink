import { ApiResponse, PaginatedApiResponse } from '@chainlink/json-api-client'
import {
  Initiator,
  JobRun,
  JobSpec,
  JobSpecError,
  OcrJobRun,
  OcrJobSpec,
  TaskSpec,
  RunStatus,
} from 'core/store/models'

export type JobRunsResponse =
  | PaginatedApiResponse<JobRun[]>
  | PaginatedApiResponse<OcrJobRun[]>

export type JobSpecResponse = ApiResponse<JobSpec> | ApiResponse<OcrJobSpec>

export type BaseJob = {
  createdAt: string
  errors: JobSpecError[]
  id: string
  name?: string
}

export type BaseJobRun = {
  createdAt: string
  id: string
  status: RunStatus
  jobId: string
}

export type OffChainReportingJob = BaseJob & {
  definition: string
  type: 'Off-chain reporting'
}

export type DirectRequestJob = BaseJob & {
  definition: { [key: string]: any }
  earnings: number | null
  endAt: string | null
  initiators: Initiator[]
  minPayment?: string | null
  startAt: string | null
  tasks: TaskSpec[]
  type: 'Direct request'
}

export type JobData = {
  job?: DirectRequestJob | OffChainReportingJob
  jobSpec?: JobSpecResponse['data']
  recentRuns?: BaseJobRun[]
  recentRunsCount: number
}
