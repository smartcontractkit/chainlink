import * as sessionStorage from 'utils/sessionStorage'
import {
  REQUEST_SIGNIN,
  RECEIVE_SIGNIN_SUCCESS,
  RECEIVE_SIGNIN_FAIL,
  RECEIVE_SIGNIN_ERROR,
  REQUEST_SIGNOUT,
  RECEIVE_SIGNOUT_SUCCESS,
  RECEIVE_SIGNOUT_ERROR
} from 'actions'

const defaultState = {
  fetching: false,
  authenticated: false,
  errors: [],
  networkError: false
}

const initialState = Object.assign(
  {},
  defaultState,
  sessionStorage.get()
)

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_SIGNOUT:
    case REQUEST_SIGNIN:
      return Object.assign(
        {},
        state,
        {
          fetching: true,
          networkError: false
        }
      )
    case RECEIVE_SIGNOUT_SUCCESS:
    case RECEIVE_SIGNIN_SUCCESS: {
      const auth = {authenticated: action.authenticated}
      sessionStorage.set(auth)
      return Object.assign(
        {},
        state,
        auth,
        {
          fetching: false,
          errors: action.errors || [],
          networkError: false
        }
      )
    }
    case RECEIVE_SIGNIN_FAIL: {
      const auth = {authenticated: false}
      sessionStorage.set(auth)
      return Object.assign(
        {},
        state,
        auth,
        {
          fetching: false,
          errors: []
        }
      )
    }
    case RECEIVE_SIGNIN_ERROR:
    case RECEIVE_SIGNOUT_ERROR: {
      const auth = {authenticated: false}
      sessionStorage.set(auth)
      return Object.assign(
        {},
        state,
        auth,
        {
          fetching: false,
          errors: action.errors || [],
          networkError: action.networkError
        }
      )
    }
    default:
      return state
  }
}
