import { EntityRepository, EntityManager } from 'typeorm'
import { Head } from '../entity/Head'
import { PaginationParams } from '../utils/pagination'

@EntityRepository()
export class EthereumHeadRepository {
  constructor(private manager: EntityManager) {}

  /**
   * Get a page of EthereumNode's sorted by their index in ascending order
   */
  public all(params: PaginationParams): Promise<Head[]> {
    let query = this.manager
      .createQueryBuilder(Head, 'ethereum_head')
      .orderBy('ethereum_head.id', 'ASC')

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
   * Return the total count of EthereumNode's
   */
  public count(): Promise<number> {
    return this.manager.createQueryBuilder(Head, 'ethereum_head').getCount()
  }
}
