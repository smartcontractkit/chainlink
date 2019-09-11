import * as jsonapi from 'api/transport/json'
import * as presenters from 'core/store/presenters'

/**
 * Index returns paginated transaction attempts
 */
type IndexParams = jsonapi.PaginatedRequestParams
const INDEX_ENDPOINT = '/v2/transactions'
const index = jsonapi.fetchResource<IndexParams, presenters.Tx[]>(
  INDEX_ENDPOINT,
)

/**
 * Show returns the details of a Ethereum Transasction details.
 *
 * @example "<application>/transactions/:TxHash"
 */
interface ShowPathParams {
  txHash: string
}
const SHOW_ENDPOINT = '/v2/transactions/:txHash'
const show = jsonapi.fetchResource<undefined, presenters.Tx, ShowPathParams>(
  SHOW_ENDPOINT,
)

export function getTransactions(
  page: number,
  size: number,
): Promise<jsonapi.PaginatedApiResponse<presenters.Tx[]>> {
  return index({ page, size })
}

export function getTransaction(
  txHash: string,
): Promise<jsonapi.ApiResponse<presenters.Tx>> {
  return show(undefined, { txHash })
}
