import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'
import * as presenters from 'core/store/presenters'
/**
 * Create adds validates, saves a new OcrKey.
 *
 * @example "<application>/off_chain_reporting_keys"
 */
const CREATE_ENDPOINT = '/v2/off_chain_reporting_keys'

/**
 * Index lists OcrKeys.
 *
 * @example "<application>/off_chain_reporting_keys"
 */
const INDEX_ENDPOINT = '/v2/off_chain_reporting_keys'

/**
 * Destroy deletes a OcrKey.
 *
 * @example "<application>/off_chain_reporting_keys/:keyId"
 */
interface DestroyPathParams {
  keyId: string
}
const DESTROY_ENDPOINT = '/v2/off_chain_reporting_keys/:keyId'

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
  ): Promise<jsonapi.ApiResponse<presenters.OcrKey>> {
    return this.create(OcrKeyRequest)
  }

  @boundMethod
  public destroyOcrKey(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { keyId: id })
  }

  private create = this.api.createResource<
    models.OcrKeyRequest,
    presenters.OcrKey
  >(CREATE_ENDPOINT)

  private index = this.api.fetchResource<{}, models.OcrKey[]>(INDEX_ENDPOINT)

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}
