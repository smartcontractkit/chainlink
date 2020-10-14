import { AppState } from 'reducers'
import build from 'redux-object'

export default ({
  dashboardIndex,
  jobRuns,
}: Pick<AppState, 'dashboardIndex' | 'jobRuns'>) => {
  return (
    dashboardIndex.recentJobRuns &&
    dashboardIndex.recentJobRuns
      .map((id) => build(jobRuns, 'items', id))
      .filter((r) => r)
  )
}
