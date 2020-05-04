import { Actions } from './actions'
import { Reducer } from 'redux'
import { JobRun } from 'explorer/models'

export interface State {
  loading: boolean
  error: boolean
  items?: Record<JobRun['id'], JobRun>
}

const INITIAL_STATE: State = {
  loading: false,
  error: false,
}

export const jobRunsReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_JOB_RUNS_BEGIN':
      return {
        ...state,
        loading: true,
        error: false,
      }

    case 'FETCH_JOB_RUNS_SUCCEEDED':
      return {
        ...state,
        loading: false,
        items: { ...action.data.jobRuns },
      }

    case 'FETCH_JOB_RUNS_ERROR':
      return {
        ...state,
        loading: false,
        error: true,
      }

    case 'FETCH_JOB_RUN_BEGIN':
      return {
        ...state,
        loading: true,
      }

    case 'FETCH_JOB_RUN_SUCCEEDED':
      return {
        ...state,
        loading: false,
        items: { ...action.data.jobRuns },
      }

    case 'FETCH_JOB_RUN_ERROR':
      return {
        ...state,
        loading: false,
        error: true,
      }

    default:
      return state
  }
}

export default jobRunsReducer
