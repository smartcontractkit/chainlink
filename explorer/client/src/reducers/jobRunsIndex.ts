import { Actions } from './actions'
import { Reducer } from 'redux'

export interface State {
  items?: string[]
  count?: number
}

const INITIAL_STATE: State = {}

export const jobRunsIndexReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_JOB_RUNS_SUCCEEDED':
      return {
        ...state,
        items: action.data.meta.currentPageJobRuns.data.map(r => r.id),
        count: action.data.meta.currentPageJobRuns.meta.count,
      }

    case 'FETCH_JOB_RUN_SUCCEEDED':
      return INITIAL_STATE

    default:
      return state
  }
}

export default jobRunsIndexReducer
