const initialState = {
  recentJobRuns: null
}

export const UPSERT_RECENT_JOB_RUNS = 'UPSERT_RECENT_JOB_RUNS'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_RECENT_JOB_RUNS: {
      return Object.assign({}, state, {
        recentJobRuns: action.data.meta.recentJobRuns.data.map(r => r.id)
      })
    }
    default:
      return state
  }
}
