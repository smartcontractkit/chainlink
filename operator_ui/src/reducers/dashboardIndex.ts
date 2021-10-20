import { Reducer } from 'redux'
import { Actions, ResourceActionType } from './actions'

export interface State {
  recentJobRuns?: string[]
  jobRunsCount?: number
}

const INITIAL_STATE: State = {
  recentJobRuns: undefined,
  jobRunsCount: undefined,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case ResourceActionType.UPSERT_RECENT_JOB_RUNS: {
      const recent = action.data.meta.recentJobRuns
      const recentJobRuns = recent.data.map((r) => r.id)
      const jobRunsCount = recent.meta.count

      return {
        ...state,
        recentJobRuns,
        jobRunsCount,
      }
    }
    default:
      return state
  }
}

export default reducer
