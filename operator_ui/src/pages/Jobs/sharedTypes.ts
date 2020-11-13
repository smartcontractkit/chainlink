import { ApiResponse, PaginatedApiResponse } from '@chainlink/json-api-client'
import {
  JobSpec,
  OcrJobSpec,
  JobRun,
  OcrJobRun,
  JobSpecError,
  Initiator,
} from 'core/store/models'

export type JobRunsResponse =
  | PaginatedApiResponse<JobRun[]>
  | PaginatedApiResponse<OcrJobRun[]>

export type JobSpecResponse = ApiResponse<JobSpec> | ApiResponse<OcrJobSpec>

export type Job = {
  createdAt: string
  definition: unknown
  errors: JobSpecError[]
  id: string
  initiators: undefined | Initiator[]
  name?: string
  type: 'Off-chain-reporting' | 'Direct request'
}

export type JobData = {
  job?: Job
  jobSpec?: JobSpecResponse['data']
  recentRuns?: JobRunsResponse['data']
  recentRunsCount: number
}
