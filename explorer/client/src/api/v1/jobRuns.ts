import * as jsonapi from '@chainlink/json-api-client'
import * as models from 'explorer/models'

/**
 * Index lists JobRuns, one page at a time.
 *
 * @example "<application>/api/v1/job_runs?size=1&page=2"
 */
const INDEX_ENDPOINT = '/api/v1/job_runs'
interface IndexRequestParams extends jsonapi.PaginatedRequestParams {
  query: string | undefined
}
const index = jsonapi.fetchResource<IndexRequestParams, models.ChainlinkNode[]>(
  INDEX_ENDPOINT,
)

/**
 * Show returns the details of a JobRun.
 *
 * @example "<application>/api/v1/job_runs/:id"
 */
interface ShowPathParams {
  id: string
}
const SHOW_ENDPOINT = `/api/v1/job_runs/:id`
const show = jsonapi.fetchResource<{}, models.JobRun, ShowPathParams>(
  SHOW_ENDPOINT,
)

/**
 * Index lists JobRuns, one page at a time.
 * @param query The token to search
 * @param page The page number to fetch
 * @param size The maximum number of job runs in the page
 */
export function getJobRuns(
  query: string | undefined,
  page: number,
  size: number,
): Promise<jsonapi.PaginatedApiResponse<models.JobRun[]>> {
  return index({ query, page, size })
}

/**
 * Get the details of a single JobRun by id
 * @param id The id of the JobRun to obtain
 */
export function getJobRun(
  id: string,
): Promise<jsonapi.ApiResponse<models.JobRun>> {
  return show({}, { id })
}
