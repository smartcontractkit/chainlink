import {
  REQUEST_CREATE,
  RECEIVE_CREATE_SUCCESS,
  MATCH_ROUTE
} from 'actions'

const initialState = {
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
        initialState,
        {
          successMessage: action.response,
          networkError: false
        }
      )
    default:
      return initialState
  }
}
