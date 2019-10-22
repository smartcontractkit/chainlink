import { Action, Dispatch } from 'redux'
import { ThunkAction } from 'redux-thunk'
import normalize from 'json-api-normalizer'
import * as api from '../api'
import { JobRunsAction } from '../reducers/jobRuns'
import { Query } from '../reducers/search'
import { State as AppState } from '../reducers'

export function getJobRuns(query: Query, page: number, size: number) {
  return (dispatch: Dispatch) => {
    return api.getJobRuns(query, page, size).then((r: JobRun[]) => {
      const normalizedData = normalize(r, { endpoint: 'jobRuns' })
      const action: JobRunsAction = {
        type: 'UPSERT_JOB_RUNS',
        data: normalizedData,
      }

      dispatch(action)
    })
  }
}

export function getJobRun(
  jobRunId?: string,
): ThunkAction<Promise<void>, AppState, void, Action<string>> {
  return (dispatch: Dispatch) => {
    return api.getJobRun(jobRunId).then((r: JobRun) => {
      const normalizedData = normalize(r, { endpoint: 'jobRun' })
      const action: JobRunsAction = {
        type: 'UPSERT_JOB_RUN',
        data: normalizedData,
      }

      dispatch(action)
    })
  }
}
