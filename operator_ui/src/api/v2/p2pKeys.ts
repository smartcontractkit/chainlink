import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'
/**
 * Create adds validates, saves a new P2P key.
 *
 * @example "POST <application>/keys/p2p"
 */
export const ENDPOINT = '/v2/keys/p2p'

/**
 * Index lists P2P Keys.
 *
 * @example "GET <application>/keys/p2p"
 */
export const INDEX_ENDPOINT = ENDPOINT

/**
 * Destroy deletes a P2P Keys.
 *
 * @example "DELETE <application>/keys/p2p/:keyId"
 */
interface DestroyPathParams {
  keyId: string
}
export const DESTROY_ENDPOINT = `${ENDPOINT}/:keyId`

export class P2PKeys {
  constructor(private api: jsonapi.Api) {}

  /**
   * Index lists P2PKeys
   */
  @boundMethod
  public getP2PKeys(): Promise<jsonapi.ApiResponse<models.P2PKey[]>> {
    return this.index()
  }

  @boundMethod
  public createP2PKey(): Promise<jsonapi.ApiResponse<models.P2PKey>> {
    return this.create()
  }

  @boundMethod
  public destroyP2PKey(id: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { keyId: id })
  }

  private create = this.api.createResource<undefined, models.P2PKey>(ENDPOINT)

  private index = this.api.fetchResource<{}, models.P2PKey[]>(INDEX_ENDPOINT)

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}
