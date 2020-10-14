import { AppState } from 'reducers'
import build from 'redux-object'

export default ({ jobs }: Pick<AppState, 'jobs'>) => {
  return (
    jobs.recentlyCreated &&
    jobs.recentlyCreated.map((id) => build(jobs, 'items', id)).filter((j) => j)
  )
}
