import {
  RECEIVE_JOB_SPEC_SUCCESS
} from 'actions'

const initialState = {
  items: {}
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case RECEIVE_JOB_SPEC_SUCCESS:
      const newJobRuns = action.item.runs.reduce(
        (acc, r) => { acc[r.id] = r; return acc },
        {}
      )

      return Object.assign(
        {},
        state,
        {items: Object.assign({}, state.items, newJobRuns)}
      )
    default:
      return state
  }
}
