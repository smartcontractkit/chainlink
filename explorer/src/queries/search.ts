import { Connection, SelectQueryBuilder } from 'typeorm'
import { JobRun } from '../entity/JobRun'

export interface ISearchParams {
  searchQuery?: string
  page?: number
  limit?: number
}

const searchBuilder = (
  db: Connection,
  params: ISearchParams
): SelectQueryBuilder<JobRun> => {
  let query = db.getRepository(JobRun).createQueryBuilder('job_run')

  if (params.searchQuery != null) {
    const searchTokens = params.searchQuery.split(/\s+/)
    query = query
      .where('job_run.runId IN(:...searchTokens)', { searchTokens })
      .orWhere('job_run.jobId IN(:...searchTokens)', { searchTokens })
      .orWhere('job_run.requester IN(:...searchTokens)', { searchTokens })
      .orWhere('job_run.requestId IN(:...searchTokens)', { searchTokens })
      .orWhere('job_run.txHash IN(:...searchTokens)', { searchTokens })
  }

  if (params.limit != null) {
    query = query.limit(params.limit)
  }

  if (params.page !== undefined) {
    const offset = (params.page - 1) * params.limit
    query = query.offset(offset)
  }

  return query
    .leftJoinAndSelect('job_run.chainlinkNode', 'chainlink_node')
    .orderBy('job_run.createdAt', 'DESC')
}

export const search = async (
  db: Connection,
  params: ISearchParams
): Promise<JobRun[]> => {
  return searchBuilder(db, params).getMany()
}

export const count = async (
  db: Connection,
  params: ISearchParams
): Promise<number> => {
  return searchBuilder(db, params).getCount()
}
