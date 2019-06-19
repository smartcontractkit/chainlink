export interface IState {
  recentJobRuns?: string
  jobRunsCount?: number
}

const initialState = {
  recentJobRuns: undefined,
  jobRunsCount: undefined
}

export type Action = { type: 'UPSERT_RECENT_JOB_RUNS'; data: any }

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_RECENT_JOB_RUNS': {
      return Object.assign({}, state, {
        recentJobRuns: action.data.meta.recentJobRuns.data.map(r => r.id),
        jobRunsCount: action.data.meta.recentJobRuns.meta.count
      })
    }
    default:
      return state
  }
}
