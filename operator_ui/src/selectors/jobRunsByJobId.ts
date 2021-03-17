import { AppState } from 'reducers'
import { JobRun } from 'operator_ui'
import build from 'redux-object'

export default (state: AppState, jobId: string, take: number) => {
  return build(state.jobRuns, 'items')
    .filter((r: JobRun) => r.jobId === jobId)
    .sort((a: JobRun, b: JobRun) => {
      const dateA = new Date(a.createdAt)
      const dateB = new Date(b.createdAt)

      return dateA < dateB ? 1 : -1
    })
    .slice(0, take)
}
