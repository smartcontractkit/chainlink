import * as jsonapi from '@chainlink/json-api-client'
import { Api } from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
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

/**
 * Show returns the details of a JobRun.
 *
 * @example "<application>/api/v1/job_runs/:id"
 */
interface ShowPathParams {
  id: string
}
const SHOW_ENDPOINT = `/api/v1/job_runs/:id`

export class JobRuns {
  constructor(private api: Api) {}

  /**
   * Index lists JobRuns, one page at a time.
   * @param query The token to search
   * @param page The page number to fetch
   * @param size The maximum number of job runs in the page
   */
  @boundMethod
  public getJobRuns(
    query: string | undefined,
    page: number,
    size: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.JobRun[]>> {
    return this.index({ query, page, size })
  }

  /**
   * Get the details of a single JobRun by id
   * @param id The id of the JobRun to obtain
   */
  @boundMethod
  public getJobRun(id: string): Promise<jsonapi.ApiResponse<models.JobRun>> {
    return this.show({}, { id })
  }

  private index = this.api.fetchResource<
    IndexRequestParams,
    models.ChainlinkNode[]
  >(INDEX_ENDPOINT)

  private show = this.api.fetchResource<{}, models.JobRun, ShowPathParams>(
    SHOW_ENDPOINT,
  )
}
