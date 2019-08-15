import { PaginatedResponse, Response, TODO } from 'api/'
import * as models from 'core/store/models'

/**
 * Index lists JobSpecs, one page at a time.
 * @param page The page number to fetch
 * @param size The maximum number of job specs in the page
 */
export const getJobSpecs = (
  page: number,
  size: number
): PaginatedResponse<models.JobSpec[]> => {
  return TODO(page, size)
}

/**
 * Get the most recent n job specs
 * @param n The number of job specs to fetch
 */
export const getRecentJobSpecs = (
  n: number
): PaginatedResponse<models.JobSpec[]> => {
  return TODO()
}

/**
 * Get the details of a single JobSpec by id
 * @param id The id of the JobSpec to obtain
 */
export const getJobSpec = (id: string): Response<models.JobSpec> => {
  return TODO(id)
}
