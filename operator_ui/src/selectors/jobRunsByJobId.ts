import build from 'redux-object'
import { IJobRun } from '../../@types/operator_ui'
import { IState } from '../connectors/redux/reducers'

export default (state: IState, jobId: string, take: number) => {
  return build(state.jobRuns, 'items')
    .filter((r: IJobRun) => r.jobId === jobId)
    .sort((a: IJobRun, b: IJobRun) => {
      const dateA = new Date(a.createdAt)
      const dateB = new Date(b.createdAt)

      return dateA < dateB ? 1 : -1
    })
    .slice(0, take)
}
