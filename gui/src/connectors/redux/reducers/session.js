import {
  RECEIVE_SESSION_SUCCESS,
  RECEIVE_SESSION_ERROR,
  RECEIVE_SIGNOUT_SUCCESS,
  RECEIVE_SIGNOUT_ERROR
} from 'actions'

const initialState = {
  authenticated: false,
  errors: [],
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case RECEIVE_SIGNOUT_SUCCESS:
    case RECEIVE_SESSION_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          authenticated: action.authenticated,
          errors: action.errors,
          networkError: false
        }
      )
    case RECEIVE_SESSION_ERROR:
    case RECEIVE_SIGNOUT_ERROR:
      return Object.assign(
        {},
        state,
        {
          authenticated: false,
          errors: action.errors || [],
          networkError: action.networkError
        }
      )
    default:
      return state
  }
}
