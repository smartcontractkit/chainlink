import * as jsonapi from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as presenters from 'core/store/presenters'

/**
 * AccountBalances returns the account balances of ETH & LINK.
 *
 * @example "<application>/user/balances"
 */
export const ACCOUNT_BALANCES_ENDPOINT = '/v2/user/balances'

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
