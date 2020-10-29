import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'
/**
 * Create adds validates, saves a new OcrKey.
 *
 * @example "POST <application>/off_chain_reporting_keys"
 */
const ENDPOINT = '/v2/off_chain_reporting_keys'

/**
 * Index lists OcrKeys.
 *
 * @example "GET <application>/off_chain_reporting_keys"
 */
const INDEX_ENDPOINT = ENDPOINT

/**
 * Destroy deletes a OcrKey.
 *
 * @example "DELETE <application>/off_chain_reporting_keys/:keyId"
 */
interface DestroyPathParams {
  keyId: string
}
const DESTROY_ENDPOINT = `${ENDPOINT}/:keyId`

export class OcrKeys {
  constructor(private api: jsonapi.Api) {}

  /**
   * Index lists OcrKeys
   */
  @boundMethod
  public getOcrKeys(): Promise<jsonapi.ApiResponse<models.OcrKey[]>> {
    return this.index()
  }

  @boundMethod
  public createOcrKey(
    OcrKeyRequest: models.OcrKeyRequest,
  ): Promise<jsonapi.ApiResponse<models.OcrKey>> {
    return this.create(OcrKeyRequest)
  }

  @boundMethod
  public destroyOcrKey(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { keyId: id })
  }

  private create = this.api.createResource<models.OcrKeyRequest, models.OcrKey>(
    ENDPOINT,
  )

  private index = this.api.fetchResource<{}, models.OcrKey[]>(INDEX_ENDPOINT)

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}
