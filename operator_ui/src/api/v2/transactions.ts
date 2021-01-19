import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as presenters from 'core/store/presenters'

/**
 * Index returns paginated transaction attempts
 */
type IndexParams = jsonapi.PaginatedRequestParams
const INDEX_ENDPOINT = '/v2/transactions'

/**
 * Show returns the details of a Ethereum Transasction details.
 *
 * @example "<application>/transactions/:TxHash"
 */
interface ShowPathParams {
  txHash: string
}
const SHOW_ENDPOINT = '/v2/transactions/:txHash'

export class Transactions {
  constructor(private api: jsonapi.Api) {}

  @boundMethod
  public getTransactions(
    page: number,
    size: number,
  ): Promise<jsonapi.PaginatedApiResponse<presenters.Tx[]>> {
    return this.index({ page, size })
  }

  @boundMethod
  public getTransaction(
    txHash: string,
  ): Promise<jsonapi.ApiResponse<presenters.Tx>> {
    return this.show(undefined, { txHash })
  }

  private index = this.api.fetchResource<IndexParams, presenters.Tx[]>(
    INDEX_ENDPOINT,
  )

  private show = this.api.fetchResource<
    undefined,
    presenters.Tx,
    ShowPathParams
  >(SHOW_ENDPOINT)
}
