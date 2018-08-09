import {
  REQUEST_CREATE,
  RECEIVE_CREATE_SUCCESS,
  RECEIVE_CREATE_ERROR,
  MATCH_ROUTE
} from 'actions'

const initialState = {
  errors: [],
  successMessage: {},
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case MATCH_ROUTE:
      return Object.assign(
        {},
        initialState
      )
    case REQUEST_CREATE:
      return Object.assign(
        {},
        state,
        { networkError: false }
      )
    case RECEIVE_CREATE_SUCCESS:
      return Object.assign(
        {},
        state,
        {
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
          errors: action.error || [],
          networkError: action.networkError
        }
      )
    default:
      return state
  }
}
