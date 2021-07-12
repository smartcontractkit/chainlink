import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'core/store/models'

/**
 * Delete removes all runs given a query
 *
 * @example "<application>/bulk_delete_runs"
 */
const DELETE_ENDPOINT = '/v2/bulk_delete_runs'

export class BulkDeleteRuns {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public bulkDeleteJobRuns(
    bulkDeleteRunRequest: models.BulkDeleteRunRequest,
  ): Promise<jsonapi.ApiResponse<null>> {
    return this.destroy(bulkDeleteRunRequest)
  }

  private destroy = this.api.deleteResource<models.BulkDeleteRunRequest, null>(
    DELETE_ENDPOINT,
  )
}
