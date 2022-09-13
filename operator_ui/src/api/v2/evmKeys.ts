import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/keys/evm/chain'

export class EVMKeys {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public chain(
    request: models.EVMKeysChainRequest,
  ): Promise<jsonapi.ApiResponse<models.EVMKey>> {
    const query = new URLSearchParams()

    query.append('address', request.address)
    query.append('evmChainID', request.evmChainID)
    if (request.nextNonce !== null) {
      query.append('nextNonce', request.nextNonce)
    }
    if (request.abandon !== null) {
      query.append('abandon', String(request.abandon))
    }
    if (request.enabled !== null) {
      query.append('enabled', String(request.enabled))
    }

    const endpoint = ENDPOINT + '?' + query.toString()

    return this.api.createResource<models.EVMKeysChainRequest, models.EVMKey>(
      endpoint,
    )()
  }
}
