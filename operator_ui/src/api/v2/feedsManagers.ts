import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

export const ENDPOINT = '/v2/feeds_managers'

export class FeedsManagers {
  constructor(private api: jsonapi.Api) {}

  /**
   * Get the list of Feeds Managers
   */
  @boundMethod
  public getFeedsManagers(): Promise<
    jsonapi.ApiResponse<models.FeedsManager[]>
  > {
    return this.index()
  }

  /**
   * Creates a Feeds Manager
   */
  @boundMethod
  public createFeedsManager(
    request: models.CreateFeedsManagerRequest,
  ): Promise<jsonapi.ApiResponse<models.FeedsManager>> {
    return this.create(request)
  }

  private index = this.api.fetchResource<{}, models.FeedsManager[], {}>(
    ENDPOINT,
  )
  private create = this.api.createResource<
    models.CreateFeedsManagerRequest,
    models.FeedsManager
  >(ENDPOINT)
}
