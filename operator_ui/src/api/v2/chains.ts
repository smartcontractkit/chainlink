import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/chains/evm'

export class Chains {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getChains(): Promise<jsonapi.ApiResponse<models.Chain[]>> {
    return this.index()
  }

  @boundMethod
  public createChain(
    request: models.CreateChainRequest,
  ): Promise<jsonapi.ApiResponse<models.Chain>> {
    return this.create(request)
  }

  private index = this.api.fetchResource<{}, models.Chain[]>(ENDPOINT)

  private create = this.api.createResource<
    models.CreateChainRequest,
    models.Chain
  >(ENDPOINT)
}
