import { Actions } from './actions'
import * as authenticationStorage from '../utils/clientStorage'
import { Reducer } from 'redux'

export interface State {
  allowed: boolean
}

const INITIAL_STATE: State = {
  allowed: authenticationStorage.get('adminAllowed') || false,
}

export const adminAuthReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_SIGNIN_SUCCEEDED':
      authenticationStorage.set('adminAllowed', true)
      return { allowed: true }
    case 'FETCH_ADMIN_SIGNIN_ERROR':
    case 'FETCH_ADMIN_SIGNOUT_SUCCEEDED':
      authenticationStorage.set('adminAllowed', false)
      return { allowed: false }
    default:
      return state
  }
}

export default adminAuthReducer
