import { RECEIVE_DELETE_SUCCESS } from '../../../actions'
import pickBy from 'lodash/pickBy'

export interface IState {
  items: any
  currentPage?: number
  currentJobRunsCount?: number
}

const initialState = {
  items: {},
  currentPage: undefined,
  currentJobRunsCount: undefined
}

export const UPSERT_JOB_RUNS = 'UPSERT_JOB_RUNS'
export const UPSERT_RECENT_JOB_RUNS = 'UPSERT_RECENT_JOB_RUNS'
export const UPSERT_JOB_RUN = 'UPSERT_JOB_RUN'

export type Action =
  | { type: 'UPSERT_JOB_RUNS'; data: any }
  | { type: 'UPSERT_RECENT_JOB_RUNS'; data: any }
  | { type: 'UPSERT_JOB_RUN'; data: any }
  | { type: 'RECEIVE_DELETE_SUCCESS'; response: any }

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_JOB_RUNS': {
      return Object.assign(
        {},
        state,
        { items: Object.assign({}, state.items, action.data.runs) },
        {
          currentPage: action.data.meta.currentPageJobRuns.data.map(
            (r: any) => r.id
          )
        },
        { currentJobRunsCount: action.data.meta.currentPageJobRuns.meta.count }
      )
    }
    case 'UPSERT_RECENT_JOB_RUNS':
    case 'UPSERT_JOB_RUN': {
      return Object.assign({}, state, {
        items: Object.assign({}, state.items, action.data.runs)
      })
    }
    case 'RECEIVE_DELETE_SUCCESS': {
      const cleanUpRuns = pickBy(
        state.items,
        item => item.attributes.jobId !== action.response
      )
      return Object.assign({}, state, { items: cleanUpRuns })
    }
    default:
      return state
  }
}
