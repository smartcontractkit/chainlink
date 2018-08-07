import {
  REQUEST_CREATE,
  RECEIVE_CREATE_SUCCESS,
  RECEIVE_CREATE_ERROR
} from 'actions'

const initialState = {
  fetching: false,
  errors: [],
  networkError: false
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
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
          errors: action.errors || [],
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
