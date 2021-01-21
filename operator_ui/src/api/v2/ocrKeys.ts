import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'
/**
 * Create adds validates, saves a new OcrKey.
 *
 * @example "POST <application>/keys/ocr"
 */
export const ENDPOINT = '/v2/keys/ocr'

/**
 * Index lists OcrKeys.
 *
 * @example "GET <application>/keys/ocr"
 */
export const INDEX_ENDPOINT = ENDPOINT

/**
 * Destroy deletes a OcrKey.
 *
 * @example "DELETE <application>/keys/ocr/:keyId"
 */
interface DestroyPathParams {
  keyId: string
}
export const DESTROY_ENDPOINT = `${ENDPOINT}/:keyId`

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
  public createOcrKey(): Promise<jsonapi.ApiResponse<models.OcrKey>> {
    return this.create()
  }

  @boundMethod
  public destroyOcrKey(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { keyId: id })
  }

  private create = this.api.createResource<undefined, models.OcrKey>(ENDPOINT)

  private index = this.api.fetchResource<{}, models.OcrKey[]>(INDEX_ENDPOINT)

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}
