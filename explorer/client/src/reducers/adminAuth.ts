import { Actions } from './actions'
import * as authenticationStorage from '../utils/clientStorage'

export interface State {
  allowed: boolean
}

const INITIAL_STATE: State = {
  allowed: authenticationStorage.get('adminAllowed') || false,
}

export default (state: State = INITIAL_STATE, action: Actions) => {
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
