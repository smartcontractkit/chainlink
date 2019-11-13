import 'whatwg-fetch'
import { Query } from './reducers/search'
import { ADMIN_USERNAME_HEADER, ADMIN_PASSWORD_HEADER } from './constants'

const base = () => `${document.location.protocol}//${document.location.host}`

export async function getJobRuns(
  query: Query,
  page: number,
  size: number,
): Promise<JobRun[]> {
  const url = new URL('/api/v1/job_runs', base())
  url.searchParams.set('page', page.toString())
  url.searchParams.set('size', size.toString())
  if (query) {
    url.searchParams.set('query', query)
  }

  const r: Response = await fetch(url.toString())

  return r.json()
}

export async function getJobRun(jobRunId?: string): Promise<JobRun> {
  const url = new URL(`/api/v1/job_runs/${jobRunId}`, base())
  const r: Response = await fetch(url.toString())

  return r.json()
}

export async function getOperators(): Promise<number> {
  const url = new URL('/api/v1/admin/operators', base())
  const r: Response = await fetch(url.toString())

  return r.status
}

export async function signIn(
  username: string,
  password: string,
): Promise<number> {
  const url = new URL('/api/v1/admin/login', base())
  const r: Response = await fetch(url.toString(), {
    method: 'POST',
    headers: {
      [ADMIN_USERNAME_HEADER]: username,
      [ADMIN_PASSWORD_HEADER]: password,
    },
  })

  return r.status
}

export async function signOut(): Promise<number> {
  const url = new URL('/api/v1/admin/logout', base())
  const r: Response = await fetch(url.toString(), { method: 'DELETE' })

  return r.status
}
