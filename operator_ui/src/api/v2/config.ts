import * as jsonapi from 'utils/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as presenters from 'core/store/presenters'

/**
 * Show returns the whitelist of config variables
 *
 * @example "<application>/config"
 */
const SHOW_ENDPOINT = '/v2/config'

export class Config {
  constructor(private api: jsonapi.Api) {}

  /**
   * Get configuration variables
   */
  @boundMethod
  public getConfiguration(): Promise<
    jsonapi.ApiResponse<presenters.ConfigPrinter>
  > {
    return this.show()
  }

  private show = this.api.fetchResource<{}, presenters.ConfigPrinter, {}>(
    SHOW_ENDPOINT,
  )
}
