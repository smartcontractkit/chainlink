import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/nodes'

export class Nodes {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getNodes(): Promise<jsonapi.ApiResponse<models.Node[]>> {
    return this.index()
  }

  private index = this.api.fetchResource<{}, models.Node[]>(ENDPOINT)
}
