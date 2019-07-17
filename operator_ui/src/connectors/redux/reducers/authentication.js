import * as authenticationStorage from 'utils/authenticationStorage'
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
  allowed: false,
  errors: [],
  networkError: false
}

const initialState = Object.assign(
  {},
  defaultState,
  authenticationStorage.get()
)

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REQUEST_SIGNOUT:
    case REQUEST_SIGNIN:
      return Object.assign({}, state, { networkError: false })
    case RECEIVE_SIGNOUT_SUCCESS:
    case RECEIVE_SIGNIN_SUCCESS: {
      const allowed = { allowed: action.authenticated }
      authenticationStorage.set(allowed)
      return Object.assign({}, state, allowed, {
        errors: [],
        networkError: false
      })
    }
    case RECEIVE_SIGNIN_FAIL: {
      const allowed = { allowed: false }
      authenticationStorage.set(allowed)
      return Object.assign({}, state, allowed, { errors: [] })
    }
    case RECEIVE_SIGNIN_ERROR:
    case RECEIVE_SIGNOUT_ERROR: {
      const allowed = { allowed: false }
      authenticationStorage.set(allowed)
      return Object.assign({}, state, allowed, {
        errors: action.errors || [],
        networkError: action.networkError
      })
    }
    default:
      return state
  }
}
