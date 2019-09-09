import * as jsonapi from 'api/transport/json'
import * as models from 'core/store/models'
import * as presenters from 'core/store/presenters'
/**
 * Index returns paginated JobRuns for a given JobSpec
 *
 * @example "<application>/runs?jobSpecId=:jobSpecId&size=1&page=2"
 */
export interface IndexParams extends jsonapi.PaginatedRequestParams {
  jobSpecId?: string
  sort?: '-createdAt'
}
const INDEX_ENDPOINT = '/v2/runs'
const index = jsonapi.fetchResource<IndexParams, models.JobRun[]>(
  INDEX_ENDPOINT
)

/**
 * Create starts a new Run for the requested JobSpec.
 *
 * @example "<application>/specs/:SpecID/runs"
 */
interface CreatePathParams {
  specId: string
}
const CREATE_ENDPOINT = '/v2/specs/:specId/runs'
const create = jsonapi.createResource<
  models.RunResult['data'],
  presenters.JobRun,
  CreatePathParams
>(CREATE_ENDPOINT)

/**
 * Show returns the details of a JobRun
 *
 * @example "<application>/runs/:RunID"
 */
interface ShowPathParams {
  runId: string
}
const SHOW_ENDPOINT = '/v2/runs/:runId'
const show = jsonapi.fetchResource<{}, models.JobRun, ShowPathParams>(
  SHOW_ENDPOINT
)

/**
 * Get a paginated response of job spec runs
 *
 * @param params Job spec params
 */
export function getJobSpecRuns(
  params: IndexParams
): Promise<jsonapi.PaginatedApiResponse<models.JobRun[]>> {
  return index({
    sort: '-createdAt',
    ...params
  })
}
/**
 * Get n most recent job runs
 *
 * @param n The number of recent job runs to fetch
 */
export function getRecentJobRuns(
  n: number
): Promise<jsonapi.PaginatedApiResponse<models.JobRun[]>> {
  return index({ size: n, sort: '-createdAt' })
}

/**
 * Get details of a specific job spec run
 *
 * @param id The id of the job run
 */
export function getJobSpecRun(
  id: string
): Promise<jsonapi.ApiResponse<models.JobRun>> {
  return show({}, { runId: id })
}

/**
 * Start a new run for a requested job spec
 *
 * @param id The specification id of the job spec to run
 */
export function createJobSpecRun(
  id: string
): Promise<jsonapi.ApiResponse<presenters.JobRun>> {
  return create(undefined, { specId: id })
}
