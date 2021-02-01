import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/jobs/:jobSpecId/runs'
const SHOW_ENDPOINT = `${ENDPOINT}/:runId`

export class OcrRuns {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getJobSpecRuns({
    jobSpecId,
    page,
    size,
  }: jsonapi.PaginatedRequestParams & { jobSpecId: string }): Promise<
    jsonapi.PaginatedApiResponse<models.OcrJobRun[]>
  > {
    return this.index({ page, size }, { jobSpecId })
  }

  @boundMethod
  public getJobSpecRun({
    jobSpecId,
    runId,
  }: {
    jobSpecId: string
    runId: string
  }): Promise<jsonapi.ApiResponse<models.OcrJobRun>> {
    return this.show({}, { jobSpecId, runId })
  }

  private index = this.api.fetchResource<{}, models.OcrJobRun[]>(ENDPOINT)

  private show = this.api.fetchResource<
    {},
    models.OcrJobRun,
    {
      jobSpecId: string
      runId: string
    }
  >(SHOW_ENDPOINT)
}
