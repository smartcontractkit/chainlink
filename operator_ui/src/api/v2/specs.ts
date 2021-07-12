import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'
import * as presenters from 'core/store/presenters'
/**
 * Create adds validates, saves, and starts a new JobSpec.
 *
 * @example "<application>/specs"
 */
export const CREATE_ENDPOINT = '/v2/specs'

/**
 * Index lists JobSpecs, one page at a time.
 *
 * @example "<application>/specs?size=1&page=2"
 */
interface IndexParams extends jsonapi.PaginatedRequestParams {
  sort?: '-createdAt'
}
export const INDEX_ENDPOINT = '/v2/specs'

/**
 * Show returns the details of a JobSpec.
 *
 * @example "<application>/specs/:SpecID"
 */
interface ShowPathParams {
  specId: string
}
const SHOW_ENDPOINT = `/v2/specs/:specId`

/**
 * Destroy soft deletes a job spec.
 *
 * @example "<application>/specs/:SpecID"
 */
interface DestroyPathParams {
  specId: string
}
const DESTROY_ENDPOINT = '/v2/specs/:specId'

export class Specs {
  constructor(private api: jsonapi.Api) {}

  /**
   * Index lists JobSpecs, one page at a time.
   * @param page The page number to fetch
   * @param size The maximum number of job specs in the page
   */
  @boundMethod
  public getJobSpecs(
    page: number,
    size: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.JobSpec[]>> {
    return this.index({ page, size })
  }

  /**
   * Get the most recent n job specs
   * @param n The number of job specs to fetch
   */
  @boundMethod
  public getRecentJobSpecs(
    n: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.JobSpec[]>> {
    return this.index({ size: n })
  }

  /**
   * Get the details of a single JobSpec by id
   * @param id The id of the JobSpec to obtain
   */
  @boundMethod
  public getJobSpec(id: string): Promise<jsonapi.ApiResponse<models.JobSpec>> {
    return this.show({}, { specId: id })
  }

  @boundMethod
  public createJobSpec(
    jobSpecRequest: models.JobSpecRequest,
  ): Promise<jsonapi.ApiResponse<presenters.JobSpec>> {
    return this.create(jobSpecRequest)
  }

  @boundMethod
  public destroyJobSpec(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { specId: id })
  }

  private create = this.api.createResource<
    models.JobSpecRequest,
    presenters.JobSpec
  >(CREATE_ENDPOINT)

  private index = this.api.fetchResource<IndexParams, models.JobSpec[]>(
    INDEX_ENDPOINT,
  )

  private show = this.api.fetchResource<{}, models.JobSpec, ShowPathParams>(
    SHOW_ENDPOINT,
  )

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}
