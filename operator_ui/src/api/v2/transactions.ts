import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as presenters from 'core/store/presenters'

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
  public getTransaction(
    txHash: string,
  ): Promise<jsonapi.ApiResponse<presenters.Tx>> {
    return this.show(undefined, { txHash })
  }

  private show = this.api.fetchResource<
    undefined,
    presenters.Tx,
    ShowPathParams
  >(SHOW_ENDPOINT)
}
