import { Dispatch } from 'redux'
import { normalize, schema } from 'normalizr'
import * as api from '../api'
import { JobRunsAction } from '../reducers/jobRuns'
import { Query } from '../reducers/search'
import { JobRun } from '../entities'

const getJobRuns = (query: Query) => {
  return (dispatch: Dispatch<any>) => {
    api.getJobRuns(query).then((r: IJobRun[]) => {
      const normalizedData = normalize(r, [JobRun])
      const action = {
        type: 'UPSERT_JOB_RUNS',
        data: normalizedData
      } as JobRunsAction

      dispatch(action)
    })
  }
}

const getJobRun = (jobRunId?: string) => {
  return (dispatch: Dispatch<any>) => {
    api.getJobRun(jobRunId).then((r: IJobRun) => {
      const normalizedData = normalize(r, JobRun)
      const action = {
        type: 'UPSERT_JOB_RUN',
        data: normalizedData
      } as JobRunsAction

      dispatch(action)
    })
  }
}

export { getJobRuns, getJobRun }
