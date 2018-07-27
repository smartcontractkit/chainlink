import {
  REQUEST_SESSION,
  REQUEST_SIGNOUT,
  RECEIVE_SESSION_SUCCESS,
  RECEIVE_SESSION_ERROR,
  RECEIVE_SIGNOUT_SUCCESS,
  RECEIVE_SIGNOUT_ERROR
} from 'actions'

const initialState = {
  fetching: false,
  authenticated: false,
  errors: [],
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_SIGNOUT:
    case REQUEST_SESSION:
      return Object.assign(
        {},
        state,
        {
          fetching: true,
          networkError: false
        }
      )
    case RECEIVE_SIGNOUT_SUCCESS:
    case RECEIVE_SESSION_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          fetching: false,
          authenticated: action.authenticated,
          errors: action.errors || [],
          networkError: false
        }
      )
    case RECEIVE_SESSION_ERROR:
    case RECEIVE_SIGNOUT_ERROR:
      return Object.assign(
        {},
        state,
        {
          fetching: false,
          authenticated: false,
          errors: action.errors || [],
          networkError: action.networkError
        }
      )
    default:
      return state
  }
}
