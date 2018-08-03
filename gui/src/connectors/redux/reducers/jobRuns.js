import {
  RECEIVE_JOB_SPEC_SUCCESS,
  RECEIVE_JOB_SPEC_RUNS_SUCCESS,
  RECEIVE_JOB_SPEC_RUN_SUCCESS
} from 'actions'

const initialState = {
  currentPage: [],
  items: {}
}

export const LATEST_JOB_RUNS_COUNT = 5

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case RECEIVE_JOB_SPEC_RUNS_SUCCESS: {
      const runs = (action.items || [])
      const mapped = runs.reduce(
        (acc, r) => { acc[r.id] = r; return acc },
        {}
      )

      return Object.assign(
        {},
        state,
        {
          currentPage: runs.map(jr => jr.id),
          items: Object.assign({}, state.items, mapped)
        }
      )
    }
    case RECEIVE_JOB_SPEC_RUN_SUCCESS: {
      return Object.assign(
        {},
        state,
        {items: Object.assign({}, state.items, {[action.item.id]: action.item})}
      )
    }
    case RECEIVE_JOB_SPEC_SUCCESS: {
      const runs = action.item.runs || []
      const runsMap = runs.reduce(
        (acc, r) => { acc[r.id] = r; return acc },
        {}
      )

      return Object.assign(
        {},
        state,
        {
          currentPage: runs.map(jr => jr.id).slice(0, LATEST_JOB_RUNS_COUNT),
          items: Object.assign({}, state.items, runsMap)
        }
      )
    }
    default:
      return state
  }
}
