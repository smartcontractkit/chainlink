import * as jsonapi from 'api/transport/json'
import * as presenters from 'core/store/presenters'
/**
 * Show returns the whitelist of config variables
 *
 * @example "<application>/config"
 */
const SHOW_ENDPOINT = '/v2/config'
const show = jsonapi.fetchResource<{}, presenters.ConfigWhitelist, {}>(
  SHOW_ENDPOINT,
)

/**
 * Get configuration variables
 */
export function getConfiguration(): Promise<
  jsonapi.ApiResponse<presenters.ConfigWhitelist>
> {
  return show()
}
