import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const INDEX_ENDPOINT = '/v2/features'

export class Features {
  constructor(private api: jsonapi.Api) {}

  /**
   * Get the list of CSA Keys
   */
  @boundMethod
  public getFeatureFlags(): Promise<jsonapi.ApiResponse<models.FeatureFlag[]>> {
    return this.index()
  }

  private index = this.api.fetchResource<{}, models.FeatureFlag[], {}>(
    INDEX_ENDPOINT,
  )
}
