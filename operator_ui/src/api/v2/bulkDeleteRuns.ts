import * as jsonapi from 'api/transport/json'
import * as models from 'core/store/models'

/**
 * Delete removes all runs given a query
 *
 * @example "<application>/bulk_delete_runs"
 */
const DELETE_ENDPOINT = '/v2/bulk_delete_runs'
const destroy = jsonapi.deleteResource<models.BulkDeleteRunRequest, null>(
  DELETE_ENDPOINT
)

export const bulkDeleteJobRuns = (
  bulkDeleteRunRequest: models.BulkDeleteRunRequest
) => destroy(bulkDeleteRunRequest)
