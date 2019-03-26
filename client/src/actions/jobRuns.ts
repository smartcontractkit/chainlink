import { Dispatch } from 'redux'
import * as api from '../api'
import { JobRunsAction } from '../reducers/jobRuns'
import { Query } from '../reducers/search'

const getJobRuns = (query: Query) => {
  return (dispatch: Dispatch<any>) => {
    api.getJobRuns(query).then((r: IJobRun[]) => {
      const action = { type: 'UPSERT_JOB_RUNS', items: r } as JobRunsAction
      dispatch(action)
    })
  }
}

export { getJobRuns }
