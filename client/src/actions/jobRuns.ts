import { Dispatch } from 'redux'
import * as api from '../api'
import { JobRunsAction } from '../reducers/jobRuns'

const getJobRuns = () => {
  return (dispatch: Dispatch<any>) => {
    api.getJobRuns().then((r: IJobRun[]) => {
      const action = { type: 'UPSERT_JOB_RUNS', items: r } as JobRunsAction
      dispatch(action)
    })
  }
}

export { getJobRuns }
