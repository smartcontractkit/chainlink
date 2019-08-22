import * as jsonapi from 'api/transport/json'
import * as models from 'core/store/models'
import * as presenters from 'core/store/presenters'
/**
 * Create adds validates, saves, and starts a new JobSpec.
 *
 * @example "<application>/specs"
 */
const CREATE_ENDPOINT = '/v2/specs'
const create = jsonapi.createResource<
  models.JobSpecRequest,
  presenters.JobSpec
>(CREATE_ENDPOINT)

/**
 * Index lists JobSpecs, one page at a time.
 *
 * @example "<application>/specs?size=1&page=2"
 */
interface IndexParams extends jsonapi.PaginatedRequestParams {
  sort?: '-createdAt'
}
const INDEX_ENDPOINT = '/v2/specs'
const index = jsonapi.fetchResource<IndexParams, models.JobSpec[]>(
  INDEX_ENDPOINT
)

/**
 * Show returns the details of a JobSpec.
 *
 * @example "<application>/specs/:SpecID"
 */
interface ShowPathParams {
  specId: string
}
const SHOW_ENDPOINT = `/v2/specs/:specId`
const show = jsonapi.fetchResource<{}, models.JobSpec, ShowPathParams>(
  SHOW_ENDPOINT
)

/**
 * Destroy soft deletes a job spec.
 *
 * @example "<application>/specs/:SpecID"
 */
interface DestroyPathParams {
  specId: string
}
const DESTROY_ENDPOINT = '/v2/specs/:specId'
const destroy = jsonapi.deleteResource<undefined, null, DestroyPathParams>(
  DESTROY_ENDPOINT
)

/**
 * Index lists JobSpecs, one page at a time.
 * @param page The page number to fetch
 * @param size The maximum number of job specs in the page
 */
export const getJobSpecs = (page: number, size: number) => index({ page, size })

/**
 * Get the most recent n job specs
 * @param n The number of job specs to fetch
 */
export const getRecentJobSpecs = (n: number) => index({ size: n })

/**
 * Get the details of a single JobSpec by id
 * @param id The id of the JobSpec to obtain
 */
export const getJobSpec = (id: string) => show({}, { specId: id })

export const createJobSpec = (jobSpecRequest: models.JobSpecRequest) =>
  create(jobSpecRequest)

export const destroyJobSpec = (id: string) => destroy(undefined, { specId: id })
