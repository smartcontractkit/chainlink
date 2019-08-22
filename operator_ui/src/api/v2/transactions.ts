import * as jsonapi from 'api/transport/json'
import * as presenters from 'core/store/presenters'

/**
 * Index returns paginated transaction attempts
 */
type IndexParams = jsonapi.PaginatedRequestParams
const INDEX_ENDPOINT = '/v2/transactions'
const index = jsonapi.fetchResource<IndexParams, presenters.Tx[]>(
  INDEX_ENDPOINT
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
  SHOW_ENDPOINT
)

export const getTransactions = (page: number, size: number) =>
  index({ page, size })

export const getTransaction = (txHash: string) => show(undefined, { txHash })
