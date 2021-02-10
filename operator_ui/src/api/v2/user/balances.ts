import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as presenters from 'core/store/presenters'

/**
 * AccountBalances returns the account balances of ETH & LINK.
 *
 * @example "<application>/keys/eth"
 */
export const ACCOUNT_BALANCES_ENDPOINT = '/v2/keys/eth'

export class Balances {
  constructor(private api: jsonapi.Api) {}

  /**
   * Get account balances in ETH and LINK
   */
  @boundMethod
  public getAccountBalances(): Promise<
    jsonapi.ApiResponse<presenters.AccountBalance[]>
  > {
    return this.accountBalances()
  }

  private accountBalances = this.api.fetchResource<
    {},
    presenters.AccountBalance[],
    {}
  >(ACCOUNT_BALANCES_ENDPOINT)
}
