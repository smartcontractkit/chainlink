import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/ocr/specs/:jobSpecId/runs'

export class OcrRuns {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getJobSpecRuns({
    jobSpecId,
  }: jsonapi.PaginatedRequestParams & { jobSpecId: string }): Promise<
    jsonapi.PaginatedApiResponse<models.OcrJobRun[]>
  > {
    return this.index(undefined, { jobSpecId })
  }

  private index = this.api.fetchResource<{}, models.OcrJobRun[]>(ENDPOINT)
}
