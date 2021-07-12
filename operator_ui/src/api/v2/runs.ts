import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
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

/**
 * Create starts a new Run for the requested JobSpec.
 *
 * @example "<application>/specs/:SpecID/runs"
 */
interface CreatePathParams {
  specId: string
}
const CREATE_ENDPOINT = '/v2/specs/:specId/runs'

/**
 * Show returns the details of a JobRun
 *
 * @example "<application>/runs/:RunID"
 */
interface ShowPathParams {
  runId: string
}
const SHOW_ENDPOINT = '/v2/runs/:runId'

export class Runs {
  constructor(private api: jsonapi.Api) {}

  /**
   * Get a paginated response of job spec runs
   *
   * @param params Job spec params
   */
  @boundMethod
  public getJobSpecRuns(
    params: IndexParams,
  ): Promise<jsonapi.PaginatedApiResponse<models.JobRun[]>> {
    return this.index({
      sort: '-createdAt',
      ...params,
    })
  }
  /**
   * Get n most recent job runs
   *
   * @param n The number of recent job runs to fetch
   */
  @boundMethod
  public getRecentJobRuns(
    n: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.JobRun[]>> {
    return this.index({ size: n, sort: '-createdAt' })
  }

  /**
   * Get details of a specific job spec run
   *
   * @param id The id of the job run
   */
  @boundMethod
  public getJobSpecRun(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.JobRun>> {
    return this.show({}, { runId: id })
  }

  /**
   * Start a new run for a requested job spec
   *
   * @param id The specification id of the job spec to run
   */
  @boundMethod
  public createJobSpecRun(
    id: string,
  ): Promise<jsonapi.ApiResponse<presenters.JobRun>> {
    return this.create(undefined, { specId: id })
  }

  private index = this.api.fetchResource<IndexParams, models.JobRun[]>(
    INDEX_ENDPOINT,
  )

  private create = this.api.createResource<
    models.RunResult['data'],
    presenters.JobRun,
    CreatePathParams
  >(CREATE_ENDPOINT)

  private show = this.api.fetchResource<{}, models.JobRun, ShowPathParams>(
    SHOW_ENDPOINT,
  )
}
