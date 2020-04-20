import * as jsonapi from '@chainlink/json-api-client'
import { Api } from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'explorer/models'

/**
 * Index lists Heads, one page at a time.
 *
 * @example "<application>/api/v1/admin/nodes?size=1&page=2"
 */
const INDEX_ENDPOINT = '/api/v1/admin/heads'
type IndexRequestParams = jsonapi.PaginatedRequestParams

interface ShowPathParams {
  id: number
}
const SHOW_ENDPOINT = '/api/v1/admin/heads/:id'

export class Heads {
  constructor(private api: Api) {}

  /**
   * getHeads lists Heads, one page at a time.
   * @param page The page number to fetch
   * @param size The maximum number of operators in the page
   */
  @boundMethod
  public getHeads(
    page: number,
    size: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.Head[]>> {
    return this.index({ page, size })
  }

  private index = this.api.fetchResource<IndexRequestParams, models.Head[]>(
    INDEX_ENDPOINT,
  )

  /**
   * getHead gets the full details for an individual Head
   */
  @boundMethod
  public getHead(
    id: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.Head>> {
    return this.show({}, { id })
  }

  private show = this.api.fetchResource<{}, models.Head, ShowPathParams>(
    SHOW_ENDPOINT,
  )
}
