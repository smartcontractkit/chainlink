import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/jobs'
const SHOW_ENDPOINT = `${ENDPOINT}/:specId`
const DESTROY_ENDPOINT = `${ENDPOINT}/:specId`
const RUN_JOB_ENDPOINT = `${ENDPOINT}/:specId/runs`

// Jobs represents the v2 jobs
export class Jobs {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getJobSpecs(): Promise<jsonapi.ApiResponse<models.JobSpecV2[]>> {
    return this.index()
  }

  @boundMethod
  public getJobSpec(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.JobSpecV2>> {
    return this.show({}, { specId: id })
  }

  @boundMethod
  public createJobSpec(
    request: models.JobSpecV2Request,
  ): Promise<jsonapi.ApiResponse<models.JobSpecV2>> {
    return this.create(request)
  }

  @boundMethod
  public destroyJobSpec(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { specId: id })
  }

  @boundMethod
  public createJobRunV2(
    id: string,
    pipelineInput: string,
  ): Promise<jsonapi.ApiResponse<null>> {
    return this.post(pipelineInput, { specId: id })
  }

  private index = this.api.fetchResource<{}, models.JobSpecV2[]>(ENDPOINT)

  private create = this.api.createResource<
    models.JobSpecV2Request,
    models.JobSpecV2
  >(ENDPOINT)

  private show = this.api.fetchResource<
    {},
    models.JobSpecV2,
    {
      specId: string
    }
  >(SHOW_ENDPOINT)

  private destroy = this.api.deleteResource<
    undefined,
    null,
    {
      specId: string
    }
  >(DESTROY_ENDPOINT)

  private post = this.api.createResource<
    string,
    null,
    {
      specId: string
    }
  >(RUN_JOB_ENDPOINT, true)
}
