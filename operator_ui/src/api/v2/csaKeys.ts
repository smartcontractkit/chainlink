import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/keys/csa'

export class CSAKeys {
  constructor(private api: jsonapi.Api) {}

  /**
   * Get the list of CSA Keys
   */
  @boundMethod
  public getCSAKeys(): Promise<jsonapi.ApiResponse<models.CSAKey[]>> {
    return this.index()
  }

  /**
   * Create a CSA Key
   */
  @boundMethod
  public createCSAKey(): Promise<jsonapi.ApiResponse<models.CSAKey>> {
    return this.create()
  }

  private index = this.api.fetchResource<{}, models.CSAKey[], {}>(ENDPOINT)
  private create = this.api.createResource<undefined, models.CSAKey>(ENDPOINT)
}
