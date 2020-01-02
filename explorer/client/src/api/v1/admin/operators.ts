import * as jsonapi from '@chainlink/json-api-client'
import * as models from 'explorer/models'

/**
 * Index lists Operators, one page at a time.
 *
 * @example "<application>/api/v1/admin/nodes?size=1&page=2"
 */
const INDEX_ENDPOINT = '/api/v1/admin/nodes'
type IndexRequestParams = jsonapi.PaginatedRequestParams
const index = jsonapi.fetchResource<IndexRequestParams, models.ChainlinkNode[]>(
  INDEX_ENDPOINT,
)

/**
 * Index lists Operators, one page at a time.
 * @param page The page number to fetch
 * @param size The maximum number of operators in the page
 */
export function getOperators(
  page: number,
  size: number,
): Promise<jsonapi.PaginatedApiResponse<models.ChainlinkNode[]>> {
  return index({ page, size })
}
