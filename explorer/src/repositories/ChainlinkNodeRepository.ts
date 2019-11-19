import { EntityRepository, EntityManager } from 'typeorm'
import { ChainlinkNode } from '../entity/ChainlinkNode'
import { PaginationParams } from '../utils/pagination'

@EntityRepository()
export class ChainlinkNodeRepository {
  constructor(private manager: EntityManager) {}

  /**
   * Get a page of ChainlinkNode's sorted by their index in ascending order
   */
  public all(params: PaginationParams): Promise<ChainlinkNode[]> {
    let query = this.manager
      .createQueryBuilder(ChainlinkNode, 'chainlinkNode')
      .orderBy('chainlinkNode.createdAt', 'ASC')

    if (params.limit != null) {
      query = query.limit(params.limit)
    }

    if (params.page !== undefined) {
      const offset = (params.page - 1) * params.limit
      query = query.offset(offset)
    }

    return query.getMany()
  }

  /**
   *
   * Return the total count of ChainlinkNode's
   */
  public count(): Promise<number> {
    return this.manager
      .createQueryBuilder(ChainlinkNode, 'chainlinkNode')
      .getCount()
  }
}
