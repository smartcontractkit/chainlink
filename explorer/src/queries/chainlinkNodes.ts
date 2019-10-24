import { Connection, SelectQueryBuilder } from 'typeorm'
import { ChainlinkNode } from '../entity/ChainlinkNode'
import { PaginationParams } from './pagination'

const queryBuilder = (
  db: Connection,
  params: PaginationParams,
): SelectQueryBuilder<ChainlinkNode> => {
  let query = db
    .getRepository(ChainlinkNode)
    .createQueryBuilder('chainlink_node')

  if (params.limit != null) {
    query = query.limit(params.limit)
  }

  if (params.page !== undefined) {
    const offset = (params.page - 1) * params.limit
    query = query.offset(offset)
  }

  return query.orderBy('chainlink_node."createdAt"', 'DESC')
}

export const all = async (
  db: Connection,
  params: PaginationParams,
): Promise<ChainlinkNode[]> => {
  return queryBuilder(db, params).getMany()
}

export const count = async (
  db: Connection,
  params: PaginationParams,
): Promise<number> => {
  return queryBuilder(db, params).getCount()
}
