import {
  RECEIVE_SESSION_SUCCESS,
  RECEIVE_SESSION_ERROR
} from 'actions'

const initialState = {
  authenticated: false,
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case RECEIVE_SESSION_SUCCESS:
      return Object.assign(
        {},
        state,
        {
          authenticated: true,
          networkError: false
        }
      )
    case RECEIVE_SESSION_ERROR:
      return Object.assign(
        {},
        state,
        {
          networkError: action.networkError
        }
      )
    default:
      return state
  }
}
