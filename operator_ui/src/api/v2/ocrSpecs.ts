import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/jobs'
const SHOW_ENDPOINT = `${ENDPOINT}/:specId`
const DESTROY_ENDPOINT = `${ENDPOINT}/:specId`

export class OcrSpecs {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getJobSpecs(): Promise<jsonapi.ApiResponse<models.OcrJobSpec[]>> {
    return this.index()
  }

  @boundMethod
  public getJobSpec(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.OcrJobSpec>> {
    return this.show({}, { specId: id })
  }

  @boundMethod
  public createJobSpec(
    ocrJobSpecRequest: models.OcrJobSpecRequest,
  ): Promise<jsonapi.ApiResponse<models.OcrJobSpec>> {
    return this.create(ocrJobSpecRequest)
  }

  @boundMethod
  public destroyJobSpec(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { specId: id })
  }

  private index = this.api.fetchResource<{}, models.OcrJobSpec[]>(ENDPOINT)

  private create = this.api.createResource<
    models.OcrJobSpecRequest,
    models.OcrJobSpec
  >(ENDPOINT)

  private show = this.api.fetchResource<
    {},
    models.OcrJobSpec,
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
}
