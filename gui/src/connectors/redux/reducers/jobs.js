import {
  REQUEST_JOBS,
  RECEIVE_JOBS_SUCCESS,
  RECEIVE_JOBS_ERROR,
  RECEIVE_JOB_SPEC_SUCCESS,
  RECEIVE_JOB_SPEC_RUNS_SUCCESS
} from 'actions'

const initialState = {
  items: {},
  currentPage: [],
  count: 0,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_JOBS:
      return Object.assign(
        {},
        state,
        {networkError: false}
      )
    case RECEIVE_JOB_SPEC_RUNS_SUCCESS: {
      const runs = (action.items || [])
      if (runs.length <= 0) {
        return state
      }
      const jobId = runs[0].jobId
      const items = Object.assign(
        {},
        state.items,
        {[jobId]: { runsCount: action.runsCount }}
      )

      return Object.assign(
        {},
        state,
        {items: items}
      )
    }
    case RECEIVE_JOB_SPEC_SUCCESS: {
      const runs = (action.item.runs || [])
      const jobSpec = Object.assign(
        {},
        action.item,
        {runsCount: runs.length}
      )

      return Object.assign(
        {},
        state,
        {items: Object.assign({}, state.items, {[jobSpec.id]: jobSpec})}
      )
    }
    case RECEIVE_JOBS_SUCCESS: {
      const newJobs = action.items.reduce(
        (acc, job) => { acc[job.id] = job; return acc },
        {}
      )

      return Object.assign(
        {},
        state,
        {
          items: Object.assign({}, state.items, newJobs),
          currentPage: action.items.map(j => j.id),
          count: action.count,
          networkError: false
        }
      )
    }
    case RECEIVE_JOBS_ERROR:
      return Object.assign(
        {},
        state,
        {networkError: !!action.networkError}
      )
    default:
      return state
  }
}
