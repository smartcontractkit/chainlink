export const DEFAULT_PAGE = 1
export const DEFAULT_PAGE_SIZE = 10
export const MAX_PAGE_SIZE = 100

export interface PaginationParams {
  page?: number
  limit?: number
}

export function parseParams(query: Record<string, string>): PaginationParams {
  const page = parseInt(query.page, 10) || DEFAULT_PAGE
  const limit = Math.min(
    parseInt(query.size, 10) || DEFAULT_PAGE_SIZE,
    MAX_PAGE_SIZE,
  )

  return { page, limit }
}
