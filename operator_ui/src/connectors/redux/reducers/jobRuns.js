import { RECEIVE_DELETE_SUCCESS } from '../../../actions'
import pickBy from 'lodash/pickBy'

const initialState = {
  items: {},
  currentPage: null,
  currentJobRunsCount: null
}

export const UPSERT_JOB_RUNS = 'UPSERT_JOB_RUNS'
export const UPSERT_RECENT_JOB_RUNS = 'UPSERT_RECENT_JOB_RUNS'
export const UPSERT_JOB_RUN = 'UPSERT_JOB_RUN'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_JOB_RUNS: {
      return Object.assign(
        {},
        state,
        { items: Object.assign({}, state.items, action.data.runs) },
        {
          currentPage: action.data.meta.currentPageJobRuns.data.map(r => r.id)
        },
        { currentJobRunsCount: action.data.meta.currentPageJobRuns.meta.count }
      )
    }
    case UPSERT_RECENT_JOB_RUNS:
    case UPSERT_JOB_RUN: {
      return Object.assign({}, state, {
        items: Object.assign({}, state.items, action.data.runs)
      })
    }
    case RECEIVE_DELETE_SUCCESS: {
      const cleanUpRuns = pickBy(state.items, item => item.attributes.jobId !== action.response)
      return Object.assign({}, state, { items: cleanUpRuns })
    }
    default:
      return state
  }
}
