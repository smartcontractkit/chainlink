import * as jsonapi from '@chainlink/json-api-client'
import { Api } from '@chainlink/json-api-client'
import { boundMethod } from 'autobind-decorator'
import * as models from 'explorer/models'

/**
 * Index lists Operators, one page at a time.
 *
 * @example "<application>/api/v1/admin/nodes?size=1&page=2"
 */
const INDEX_ENDPOINT = '/api/v1/admin/nodes'
type IndexRequestParams = jsonapi.PaginatedRequestParams

interface ShowPathParams {
  id: string
}
const SHOW_ENDPOINT = '/api/v1/admin/nodes/:id'

export class Operators {
  constructor(private api: Api) {}

  /**
   * Index lists Operators, one page at a time.
   * @param page The page number to fetch
   * @param size The maximum number of operators in the page
   */
  @boundMethod
  public getOperators(
    page: number,
    size: number,
  ): Promise<jsonapi.PaginatedApiResponse<models.ChainlinkNode[]>> {
    return this.index({ page, size })
  }

  @boundMethod
  public getOperator(
    id: string,
  ): Promise<jsonapi.ApiResponse<models.ChainlinkNode>> {
    return this.show({}, { id })
  }

  private index = this.api.fetchResource<
    IndexRequestParams,
    models.ChainlinkNode[]
  >(INDEX_ENDPOINT)

  private show = this.api.fetchResource<
    {},
    models.ChainlinkNode,
    ShowPathParams
  >(SHOW_ENDPOINT)
}
