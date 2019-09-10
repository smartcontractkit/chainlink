export interface State {
  recentJobRuns?: string
  jobRunsCount?: number
}

const initialState = {
  recentJobRuns: undefined,
  jobRunsCount: undefined,
}

interface RecentJobRun {
  id: string
}

interface RecentJobRunsMeta {
  count: number
}

interface RecentJobRuns {
  data: RecentJobRun[]
  meta: RecentJobRunsMeta
}

interface Meta {
  recentJobRuns: RecentJobRuns
}

interface Data {
  meta: Meta
}

export type Action = { type: 'UPSERT_RECENT_JOB_RUNS'; data: Data }

export default (state: State = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_RECENT_JOB_RUNS': {
      return Object.assign({}, state, {
        recentJobRuns: action.data.meta.recentJobRuns.data.map(r => r.id),
        jobRunsCount: action.data.meta.recentJobRuns.meta.count,
      })
    }
    default:
      return state
  }
}
