import * as jsonapi from 'api/transport/json'
import * as models from 'core/store/models'

// Create adds the BridgeType to the given context.
const CREATE_ENDPOINT = '/v2/bridge_types'
const create = jsonapi.createResource<
  models.BridgeTypeRequest,
  models.BridgeTypeAuthentication
>(CREATE_ENDPOINT)

// Index lists Bridges, one page at a time.
const INDEX_ENDPOINT = '/v2/bridge_types'
type IndexRequestParams = jsonapi.PaginatedRequestParams
const index = jsonapi.fetchResource<IndexRequestParams, models.BridgeType[]>(
  INDEX_ENDPOINT
)

// Show returns the details of a specific Bridge.
const SHOW_ENDPOINT = '/v2/bridge_types/:bridgeName'
interface ShowPathParams {
  bridgeName: string
}
const show = jsonapi.fetchResource<{}, models.BridgeType, ShowPathParams>(
  SHOW_ENDPOINT
)

// Update can change the restricted attributes for a bridge
interface UpdatePathParams {
  bridgeName: string
}
const UPDATE_ENDPOINT = '/v2/bridge_types/:bridgeName'
const update = jsonapi.updateResource<
  models.BridgeTypeRequest,
  models.BridgeType,
  UpdatePathParams
>(UPDATE_ENDPOINT)

export function getBridges(
  page: number,
  size: number
): Promise<jsonapi.PaginatedApiResponse<models.BridgeType[]>> {
  return index({ page, size })
}

/**
 * Get a bridge spec
 *
 * @param name The name of the bridge spec to fetch
 */
export function getBridgeSpec(
  name: string
): Promise<jsonapi.ApiResponse<models.BridgeType>> {
  return show({}, { bridgeName: name })
}

/**
 * Create a bridge type from a bridge type request
 *
 * @param bridgeTypeRequest The request object to create a bridge type from
 */
export function createBridge(
  bridgeTypeRequest: models.BridgeTypeRequest
): Promise<jsonapi.ApiResponse<models.BridgeTypeAuthentication>> {
  return create(bridgeTypeRequest)
}

export function updateBridge(
  bridgeTypeRequest: models.BridgeTypeRequest
): Promise<jsonapi.ApiResponse<models.BridgeType>> {
  return update(bridgeTypeRequest, { bridgeName: bridgeTypeRequest.name })
}
