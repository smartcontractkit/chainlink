import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'
/**
 * Create validates, saves and starts a new off-chain reporting job.
 *
 * @example "POST <application>/ocr/specs"
 */
export const ENDPOINT = '/v2/ocr/specs'

export class OcrSpecs {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public createJobSpec(
    ocrJobSpecRequest: models.OcrJobSpecRequest,
  ): Promise<jsonapi.ApiResponse<models.OcrJobSpec>> {
    return this.create(ocrJobSpecRequest)
  }

  private create = this.api.createResource<
    models.OcrJobSpecRequest,
    models.OcrJobSpec
  >(ENDPOINT)
}
