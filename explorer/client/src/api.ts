import 'whatwg-fetch'
import { Query } from './reducers/search'

const base = () => `${document.location.protocol}//${document.location.host}`

const getJobRuns = async (
  query: Query,
  page: number,
  size: number
): Promise<IJobRun[]> => {
  const url = new URL('/api/v1/job_runs', base())
  url.searchParams.set('page', page.toString())
  url.searchParams.set('size', size.toString())
  if (query) {
    url.searchParams.set('query', query)
  }

  const r: Response = await fetch(url.toString())
  return r.json()
}

const getJobRun = async (jobRunId?: string): Promise<IJobRun> => {
  const url = new URL(`/api/v1/job_runs/${jobRunId}`, base())
  const r: Response = await fetch(url.toString())
  return r.json()
}

export { getJobRuns, getJobRun }
