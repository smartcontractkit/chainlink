import { ApiResponse, PaginatedApiResponse } from '@chainlink/json-api-client'
import { JobSpec, OcrJobSpec, JobRun, OcrJobRun } from 'core/store/models'

export type JobRunsResponse =
  | PaginatedApiResponse<JobRun[]>
  | PaginatedApiResponse<OcrJobRun[]>

export type JobSpecResponse = ApiResponse<JobSpec> | ApiResponse<OcrJobSpec>

export type JobData = {
  jobSpec?: JobSpecResponse['data']
  recentRuns?: JobRunsResponse['data']
  recentRunsCount: number
}
