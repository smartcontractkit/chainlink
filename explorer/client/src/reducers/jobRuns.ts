import { Actions } from './actions'
import { Reducer } from 'redux'
import { JobRun } from 'explorer/models'

export interface State {
  items?: Record<string, JobRun>
}

const INITIAL_STATE: State = {}

export const jobRunsReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_JOB_RUNS_SUCCEEDED':
      return { items: { ...action.data.jobRuns } }
    case 'FETCH_JOB_RUN_SUCCEEDED':
      return { items: { ...action.data.jobRuns } }
    default:
      return state
  }
}

export default jobRunsReducer
