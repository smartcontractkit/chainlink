import { ApiResponse, PaginatedApiResponse } from '@chainlink/json-api-client'
import { JobSpec, JobRun } from 'core/store/models'

export type JobData = {
  jobSpec?: ApiResponse<JobSpec>['data']
  recentRuns?: PaginatedApiResponse<JobRun[]>['data']
  recentRunsCount: number
}
