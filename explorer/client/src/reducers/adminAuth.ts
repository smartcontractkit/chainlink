import { Actions } from './actions'
import {
  getAdminAllowed,
  setAdminAllowed,
} from '../utils/authenticationStorage'
import { Reducer } from 'redux'

export interface State {
  allowed: boolean
}

const INITIAL_STATE: State = {
  allowed: getAdminAllowed(),
}

export const adminAuthReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_SIGNIN_SUCCEEDED':
      setAdminAllowed(true)
      return { allowed: true }
    case 'FETCH_ADMIN_SIGNIN_ERROR':
    case 'FETCH_ADMIN_SIGNOUT_SUCCEEDED':
      setAdminAllowed(false)
      return { allowed: false }
    default:
      return state
  }
}

export default adminAuthReducer
