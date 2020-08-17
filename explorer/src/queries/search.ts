import { getRepository, SelectQueryBuilder } from 'typeorm'
import { JobRun } from '../entity/JobRun'
import { PaginationParams } from '../utils/pagination'

export interface SearchParams extends PaginationParams {
  searchQuery?: string
}

const normalizeSearchToken = (id: string): string => {
  const MAX_LEN_UNPREFIXED_REQUESTER_HASH = 40
  const MAX_LEN_UNPREFIXED_REQUEST_TX_HASH = 64
  if (
    id &&
    id.substr(0, 2) !== '0x' &&
    (id.length === MAX_LEN_UNPREFIXED_REQUESTER_HASH ||
      id.length === MAX_LEN_UNPREFIXED_REQUEST_TX_HASH)
  )
    return `0x${id}`
  return id
}

const buildSearchTokens = (searchQuery: string): string[] => {
  let searchTokens = searchQuery.split(/\s+/)
  searchTokens = searchTokens.concat(searchTokens.map(normalizeSearchToken))
  return searchTokens
}

const searchBuilder = (searchQuery?: string): SelectQueryBuilder<JobRun> => {
  let query = getRepository(JobRun).createQueryBuilder('job_run')

  if (searchQuery != null) {
    const searchTokens = buildSearchTokens(searchQuery)
    query = query.addSelect('COUNT(*) OVER() AS "totalRecords"').where(
      `
        ARRAY["job_run"."runId", "job_run"."jobId", "job_run"."requestId", "job_run"."requester", "job_run"."txHash"] && ARRAY[:...searchTokens]::citext[]
      `,
      { searchTokens },
    )
  } else {
    query = query.where('true = false')
  }

  return query
}

const pagedSearchBuilder = (
  params: SearchParams,
): SelectQueryBuilder<JobRun> => {
  let query = searchBuilder(params.searchQuery)

  if (params.limit != null) {
    query = query.limit(params.limit)
  }

  if (params.page !== undefined) {
    const offset = (params.page - 1) * params.limit
    query = query.offset(offset)
  }

  return query
}

export class JobRunSearch {
  public results: JobRun[]
  public totalRecords: number
  private query: SelectQueryBuilder<JobRun>

  constructor(params: SearchParams) {
    this.query = pagedSearchBuilder(params)
      .leftJoinAndSelect('job_run.chainlinkNode', 'chainlink_node')
      .orderBy('job_run.createdAt', 'DESC')
  }

  async execute(): Promise<JobRunSearch> {
    const query = await this.query.getRawAndEntities()
    this.results = query.entities
    this.totalRecords = parseInt(query.raw[0]?.totalRecords, 10) || 0
    return this
  }
}

export const search = async (params: SearchParams): Promise<JobRunSearch> => {
  return new JobRunSearch(params).execute()
}
