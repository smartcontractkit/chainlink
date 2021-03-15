import * as jsonapi from 'utils/json-api-client'
import { Api } from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

// Create adds the BridgeType to the given context.
const CREATE_ENDPOINT = '/v2/bridge_types'

// Index lists Bridges, one page at a time.
const INDEX_ENDPOINT = '/v2/bridge_types'
type IndexRequestParams = jsonapi.PaginatedRequestParams

// Show returns the details of a specific Bridge.
const SHOW_ENDPOINT = '/v2/bridge_types/:bridgeName'
interface ShowPathParams {
  bridgeName: string
}

// Update can change the restricted attributes for a bridge
interface UpdatePathParams {
  bridgeName: string
}
const UPDATE_ENDPOINT = '/v2/bridge_types/:bridgeName'

// Destroy deletes a bridge.
interface DestroyPathParams {
  bridgeName: string
}
const DESTROY_ENDPOINT = '/v2/bridge_types/:bridgeName'

export class BridgeTypes {
  constructor(private api: Api) {}

  @boundMethod
  public getBridges(
    page: number,
    size: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.BridgeType[]>> {
    return this.index({ page, size })
  }

  /**
   * Get a bridge spec
   *
   * @param name The name of the bridge spec to fetch
   */
  @boundMethod
  public getBridgeSpec(
    name: string,
  ): Promise<jsonapi.ApiResponse<models.BridgeType>> {
    return this.show({}, { bridgeName: name })
  }

  /**
   * Create a bridge type from a bridge type request
   *
   * @param bridgeTypeRequest The request object to create a bridge type from
   */
  @boundMethod
  public createBridge(
    bridgeTypeRequest: models.BridgeTypeRequest,
  ): Promise<jsonapi.ApiResponse<models.BridgeTypeAuthentication>> {
    return this.create(bridgeTypeRequest)
  }

  @boundMethod
  public updateBridge(
    bridgeTypeRequest: models.BridgeTypeRequest,
  ): Promise<jsonapi.ApiResponse<models.BridgeType>> {
    return this.update(bridgeTypeRequest, {
      bridgeName: bridgeTypeRequest.name,
    })
  }

  @boundMethod
  public destroyBridge(name: string): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(undefined, { bridgeName: name })
  }

  private create = this.api.createResource<
    models.BridgeTypeRequest,
    models.BridgeTypeAuthentication
  >(CREATE_ENDPOINT)

  private index = this.api.fetchResource<
    IndexRequestParams,
    models.BridgeType[]
  >(INDEX_ENDPOINT)

  private show = this.api.fetchResource<{}, models.BridgeType, ShowPathParams>(
    SHOW_ENDPOINT,
  )

  private update = this.api.updateResource<
    models.BridgeTypeRequest,
    models.BridgeType,
    UpdatePathParams
  >(UPDATE_ENDPOINT)

  private destroy = this.api.deleteResource<undefined, null, DestroyPathParams>(
    DESTROY_ENDPOINT,
  )
}
