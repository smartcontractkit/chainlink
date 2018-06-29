import {
  REQUEST_BRIDGES,
  RECEIVE_BRIDGES_SUCCESS,
  RECEIVE_BRIDGES_ERROR
} from 'actions'

const initialState = {
  items: [],
  currentPage: [],
  count: 0,
  fetching: false,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_BRIDGES:
      return Object.assign(
        {},
        state,
        {
          fetching: true,
          networkError: false
        }
      )
    case RECEIVE_BRIDGES_SUCCESS: {
      return Object.assign(
        {},
        state,
        {
          items: Object.assign([], state.items, action.items),
          currentPage: action.items.map(b => b.name),
          count: action.count,
          fetching: false,
          networkError: false
        }
      )
    }
    case RECEIVE_BRIDGES_ERROR:
      return Object.assign(
        {},
        state,
        {
          fetching: false,
          networkError: !!action.networkError
        }
      )
    default:
      return state
  }
}
