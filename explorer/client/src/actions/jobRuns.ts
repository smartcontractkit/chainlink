import { Dispatch } from 'redux'
import normalize from 'json-api-normalizer'
import * as api from '../api'
import { JobRunsAction } from '../reducers/jobRuns'
import { Query } from '../reducers/search'

const getJobRuns = (query: Query, page: number, size: number) => {
  return (dispatch: Dispatch<any>) => {
    api.getJobRuns(query, page, size).then((r: IJobRun[]) => {
      const normalizedData = normalize(r, { endpoint: 'jobRuns' })
      const action: JobRunsAction = {
        type: 'UPSERT_JOB_RUNS',
        data: normalizedData,
      }

      dispatch(action)
    })
  }
}

const getJobRun = (jobRunId?: string) => {
  return (dispatch: Dispatch<any>) => {
    api.getJobRun(jobRunId).then((r: IJobRun) => {
      const normalizedData = normalize(r, { endpoint: 'jobRun' })
      const action: JobRunsAction = {
        type: 'UPSERT_JOB_RUN',
        data: normalizedData,
      }

      dispatch(action)
    })
  }
}

export { getJobRuns, getJobRun }
