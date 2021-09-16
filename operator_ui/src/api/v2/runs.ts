import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/jobs/:jobId/runs'
const SHOW_ENDPOINT = `${ENDPOINT}/:runId`
const ALL_RUNS_ENDPOINT = '/v2/pipeline/runs'

export class Runs {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getAllJobRuns({
    page,
    size,
  }: jsonapi.PaginatedRequestParams): Promise<
    jsonapi.PaginatedApiResponse<models.JobRunV2[]>
  > {
    return this.allRuns({ page, size })
  }

  @boundMethod
  public getJobRuns({
    jobId,
    page,
    size,
  }: jsonapi.PaginatedRequestParams & { jobId: string }): Promise<
    jsonapi.PaginatedApiResponse<models.JobRunV2[]>
  > {
    return this.index({ page, size }, { jobId })
  }

  @boundMethod
  public getJobRun({
    jobId,
    runId,
  }: {
    jobId: string
    runId: string
  }): Promise<jsonapi.ApiResponse<models.JobRunV2>> {
    return this.show({}, { jobId, runId })
  }

  private index = this.api.fetchResource<{}, models.JobRunV2[]>(ENDPOINT)

  private show = this.api.fetchResource<
    {},
    models.JobRunV2,
    {
      jobId: string
      runId: string
    }
  >(SHOW_ENDPOINT)

  private allRuns = this.api.fetchResource<{}, models.JobRunV2[]>(
    ALL_RUNS_ENDPOINT,
  )
}
