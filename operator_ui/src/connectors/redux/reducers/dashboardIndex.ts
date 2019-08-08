export interface IState {
  recentJobRuns?: string
  jobRunsCount?: number
}

const initialState = {
  recentJobRuns: undefined,
  jobRunsCount: undefined
}

interface IRecentJobRun {
  id: string
}

interface IRecentJobRunsMeta {
  count: number
}

interface IRecentJobRuns {
  data: IRecentJobRun[]
  meta: IRecentJobRunsMeta
}

interface IMeta {
  recentJobRuns: IRecentJobRuns
}

interface IData {
  meta: IMeta
}

export type Action = { type: 'UPSERT_RECENT_JOB_RUNS'; data: IData }

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_RECENT_JOB_RUNS': {
      return {
        ...state,
        recentJobRuns: action.data.meta.recentJobRuns.data.map(r => r.id),
        jobRunsCount: action.data.meta.recentJobRuns.meta.count
      }
    }
    default:
      return state
  }
}
