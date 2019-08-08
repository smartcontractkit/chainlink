import pickBy from 'lodash/pickBy'

export interface IState {
  items: { [k: string]: any }
  currentPage?: number
  currentJobRunsCount?: number
}

const initialState: IState = {
  items: {},
  currentPage: undefined,
  currentJobRunsCount: undefined
}

export type Action =
  | { type: 'UPSERT_JOB_RUNS'; data: any }
  | { type: 'UPSERT_RECENT_JOB_RUNS'; data: any }
  | { type: 'UPSERT_JOB_RUN'; data: any }
  | { type: 'RECEIVE_DELETE_SUCCESS'; response: any }

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_JOB_RUNS': {
      return {
        ...state,

        items: {
          ...state.items,
          ...action.data.runs
        },

        currentPage: action.data.meta.currentPageJobRuns.data.map(
          (r: { id: string }) => r.id
        ),

        currentJobRunsCount: action.data.meta.currentPageJobRuns.meta.count
      }
    }
    case 'UPSERT_RECENT_JOB_RUNS':
    case 'UPSERT_JOB_RUN': {
      return {
        ...state,

        items: {
          ...state.items,
          ...action.data.runs
        }
      }
    }
    case 'RECEIVE_DELETE_SUCCESS': {
      const cleanUpRuns = pickBy(
        state.items,
        item => item.attributes.jobId !== action.response
      )
      return {
        ...state,
        items: cleanUpRuns
      }
    }
    default:
      return state
  }
}
