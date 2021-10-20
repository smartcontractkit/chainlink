import { AppState } from 'reducers'
import build from 'redux-object'

export default ({ jobRuns }: Pick<AppState, 'jobRuns'>) => {
  return (
    jobRuns.currentPage &&
    jobRuns.currentPage
      .map((id) => build(jobRuns, 'items', id))
      .filter((r) => r)
  )
}
