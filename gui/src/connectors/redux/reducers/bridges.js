import {
  REQUEST_BRIDGES,
  RECEIVE_BRIDGES_SUCCESS,
  RECEIVE_BRIDGES_ERROR
} from 'actions'

const initialState = {
  items: [],
  currentPage: [],
  count: 0,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_BRIDGES:
      return Object.assign(
        {},
        state,
        {networkError: false}
      )
    case RECEIVE_BRIDGES_SUCCESS: {
      return Object.assign(
        {},
        state,
        {
          items: Object.assign([], state.items, action.items),
          currentPage: action.items.map(b => b.name),
          count: action.count,
          networkError: false
        }
      )
    }
    case RECEIVE_BRIDGES_ERROR:
      return Object.assign(
        {},
        state,
        {networkError: !!action.networkError}
      )
    default:
      return state
  }
}
