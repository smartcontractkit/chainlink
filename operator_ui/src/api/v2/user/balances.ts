import * as jsonapi from 'api/transport/json'
import * as presenters from 'core/store/presenters'

/**
 * AccountBalances returns the account balances of ETH & LINK.
 *
 * @example "<application>/user/balances"
 */
const ACCOUNT_BALANCES_ENDPOINT = '/v2/user/balances'
const accountBalances = jsonapi.fetchResource<
  {},
  presenters.AccountBalance[],
  {}
>(ACCOUNT_BALANCES_ENDPOINT)

/**
 * Get account balances in ETH and LINK
 */
export function getAccountBalances(): jsonapi.ApiResponse<
  presenters.AccountBalance[]
> {
  return accountBalances()
}
