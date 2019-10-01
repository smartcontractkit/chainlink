import { PostgresErrorCode } from './constants'

export function isError(e: unknown): e is Error {
  return !!(e instanceof Error || ((e as Error).message && (e as Error).name))
}

interface PostgresError extends Error {
  code: PostgresErrorCode
}

export function isPostgresError(e: unknown): e is PostgresError {
  return !!(e as PostgresError).code
}
