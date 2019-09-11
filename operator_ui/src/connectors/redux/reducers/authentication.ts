import * as authenticationStorage from 'utils/authenticationStorage'

const defaultState = {
  allowed: false,
  errors: [],
  networkError: false,
}

const initialState = Object.assign(
  {},
  defaultState,
  authenticationStorage.get(),
)

export type AuthenticationAction =
  | { type: AuthenticationActionType.REQUEST_SIGNIN }
  | {
      type: AuthenticationActionType.RECEIVE_SIGNIN_SUCCESS
      authenticated: boolean
    }
  | { type: AuthenticationActionType.REQUEST_SIGNOUT }
  | {
      type: AuthenticationActionType.RECEIVE_SIGNIN_ERROR
      errors: object[] // CHECKME
      networkError: boolean
    }
  | { type: AuthenticationActionType.RECEIVE_SIGNIN_FAIL }
  | {
      type: AuthenticationActionType.RECEIVE_SIGNOUT_SUCCESS
      authenticated: boolean
    }
  | {
      type: AuthenticationActionType.RECEIVE_SIGNOUT_ERROR
      errors: object[] // CHECKME
      networkError: boolean
    }

enum AuthenticationActionType {
  REQUEST_SIGNOUT = 'REQUEST_SIGNOUT',
  REQUEST_SIGNIN = 'REQUEST_SIGNIN',
  RECEIVE_SIGNOUT_SUCCESS = 'RECEIVE_SIGNOUT_SUCCESS',
  RECEIVE_SIGNIN_SUCCESS = 'RECEIVE_SIGNIN_SUCCESS',
  RECEIVE_SIGNIN_FAIL = 'RECEIVE_SIGNIN_FAIL',
  RECEIVE_SIGNIN_ERROR = 'RECEIVE_SIGNIN_ERROR',
  RECEIVE_SIGNOUT_ERROR = 'RECEIVE_SIGNOUT_ERROR',
}

export default (state = initialState, action: AuthenticationAction) => {
  switch (action.type) {
    case AuthenticationActionType.REQUEST_SIGNOUT:
    case AuthenticationActionType.REQUEST_SIGNIN:
      return Object.assign({}, state, { networkError: false })
    case AuthenticationActionType.RECEIVE_SIGNOUT_SUCCESS:
    case AuthenticationActionType.RECEIVE_SIGNIN_SUCCESS: {
      const allowed = { allowed: action.authenticated }
      authenticationStorage.set(allowed)
      return Object.assign({}, state, allowed, {
        errors: [],
        networkError: false,
      })
    }
    case AuthenticationActionType.RECEIVE_SIGNIN_FAIL: {
      const allowed = { allowed: false }
      authenticationStorage.set(allowed)
      return Object.assign({}, state, allowed, { errors: [] })
    }
    case AuthenticationActionType.RECEIVE_SIGNIN_ERROR:
    case AuthenticationActionType.RECEIVE_SIGNOUT_ERROR: {
      const allowed = { allowed: false }
      authenticationStorage.set(allowed)
      return Object.assign({}, state, allowed, {
        errors: action.errors || [],
        networkError: action.networkError,
      })
    }
    default:
      return state
  }
}
