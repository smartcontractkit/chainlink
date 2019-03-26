import 'whatwg-fetch'
import { Query } from './reducers/search'

const base = () => `${document.location.protocol}//${document.location.host}`

const getJobRuns = async (query: Query): Promise<IJobRun[]> => {
  const url = new URL('/api/v1/job_runs', base())
  if (query) { url.searchParams.set('query', query) }

  const r: Response = await fetch(url.toString())
  return r.json()
}

export { getJobRuns }
