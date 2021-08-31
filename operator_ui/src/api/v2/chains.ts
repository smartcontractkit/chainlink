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

  private index = this.api.fetchResource<{}, models.Chain[]>(ENDPOINT)
}
