import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/jobs/:jobId/runs'
const SHOW_ENDPOINT = `${ENDPOINT}/:runId`

export class Runs {
  constructor(private api: jsonapi.Api) {}

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

  private show = this.api.fetchResource<
    {},
    models.JobRunV2,
    {
      jobId: string
      runId: string
    }
  >(SHOW_ENDPOINT)
}
