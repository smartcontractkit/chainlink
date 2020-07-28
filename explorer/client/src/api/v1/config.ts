import * as jsonapi from '@chainlink/json-api-client'
import { Api } from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'explorer/models'

const CONFIG_ENDPOINT = '/api/v1/config'

export class Config {
  constructor(private api: Api) {}

  /**
   * getConfig returns app config.
   */
  @boundMethod
  public getConfig(): Promise<jsonapi.ApiResponse<models.Config>> {
    return this.api.GET(CONFIG_ENDPOINT)()
  }
}
