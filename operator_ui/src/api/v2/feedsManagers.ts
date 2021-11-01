import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

const ENDPOINT = '/v2/feeds_managers'
const UPDATE_ENDPOINT = `${ENDPOINT}/:id`

export class FeedsManagers {
  constructor(private api: jsonapi.Api) {}

  /**
   * Updates a Feeds Manager
   */
  @boundMethod
  public updateFeedsManager(
    id: string,
    request: models.UpdateFeedsManagerRequest,
  ): Promise<jsonapi.ApiResponse<models.FeedsManager>> {
    return this.update(request, { id })
  }

  private update = this.api.updateResource<
    models.UpdateFeedsManagerRequest,
    models.FeedsManager,
    { id: string }
  >(UPDATE_ENDPOINT)
}
