import { Reducer } from 'redux'
import * as storage from 'utils/storage'
import { Actions } from './actions'

export interface State {
  allowed: boolean
  errors: any[]
}

const DEFAULT_STATE = {
  allowed: false,
  errors: [],
}
const INITIAL_AUTH_STATE = storage.getAuthentication()
const INITIAL_STATE = { ...DEFAULT_STATE, ...INITIAL_AUTH_STATE }

const reducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action: Actions,
) => {
  switch (action.type) {
    case 'RECEIVE_SIGNOUT_SUCCESS':
    case 'RECEIVE_SIGNIN_SUCCESS': {
      const allowed = { allowed: action.authenticated }
      storage.setAuthentication(allowed)

      return {
        ...state,
        ...allowed,
        errors: [],
      }
    }
    case 'RECEIVE_SIGNIN_FAIL': {
      const allowed = { allowed: false }
      storage.setAuthentication(allowed)

      return {
        ...state,
        ...allowed,
        errors: [],
      }
    }
    case 'RECEIVE_SIGNIN_ERROR':
    case 'RECEIVE_SIGNOUT_ERROR': {
      const allowed = { allowed: false }
      storage.setAuthentication(allowed)

      return {
        ...state,
        ...allowed,
        errors: action.errors || [],
      }
    }
    default:
      return state
  }
}

export default reducer
