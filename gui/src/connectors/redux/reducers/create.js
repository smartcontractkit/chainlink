import {
  REQUEST_CREATE,
  RECEIVE_CREATE_SUCCESS,
  RECEIVE_CREATE_ERROR,
  MATCH_ROUTE
} from 'actions'

const initialState = {
  fetching: false,
  errors: [],
  successMessage: {},
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case MATCH_ROUTE:
      if (action.match && state.currentUrl !== action.match.url) {
        return Object.assign(
          {},
          initialState
        )
      }
      return state
    case REQUEST_CREATE:
      return Object.assign(
        {},
        state,
        {
          fetching: true,
          networkError: false
        }
      )
    case RECEIVE_CREATE_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          fetching: false,
          errors: action.error || [],
          successMessage: action.response,
          networkError: false
        }
      )
    case RECEIVE_CREATE_ERROR:
      return Object.assign(
        {},
        state,
        {
          fetching: false,
          errors: action.error || [],
          networkError: action.networkError
        }
      )
    default:
      return state
  }
}
