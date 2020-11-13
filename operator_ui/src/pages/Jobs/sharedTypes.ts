import { ApiResponse, PaginatedApiResponse } from '@chainlink/json-api-client'
import {
  Initiator,
  JobRun,
  JobSpec,
  JobSpecError,
  OcrJobRun,
  OcrJobSpec,
  TaskSpec,
} from 'core/store/models'

export type JobRunsResponse =
  | PaginatedApiResponse<JobRun[]>
  | PaginatedApiResponse<OcrJobRun[]>

export type JobSpecResponse = ApiResponse<JobSpec> | ApiResponse<OcrJobSpec>

export type BaseJob = {
  createdAt: string
  definition: unknown
  errors: JobSpecError[]
  id: string
  name?: string
}

export type OffChainReportingJob = BaseJob & {
  type: 'Off-chain reporting'
}

export type DirectRequestJob = BaseJob & {
  earnings: number | null
  initiators: Initiator[]
  minPayment?: string | null
  tasks: TaskSpec[]
  type: 'Direct request'
}

export type JobData = {
  job?: DirectRequestJob | OffChainReportingJob
  jobSpec?: JobSpecResponse['data']
  recentRuns?: JobRunsResponse['data']
  recentRunsCount: number
}
