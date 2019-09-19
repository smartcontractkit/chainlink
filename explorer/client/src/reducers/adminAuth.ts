import * as authenticationStorage from '../utils/clientStorage'

export interface State {
  allowed: boolean
}

export type Action =
  | { type: 'ADMIN_SIGNIN_SUCCEEDED' }
  | { type: 'ADMIN_SIGNIN_FAILED' }
  | { type: 'ADMIN_SIGNIN_ERROR' }
  | { type: 'ADMIN_SIGNOUT_SUCCEEDED' }
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }

const INITIAL_STATE: State = { allowed: false }

export default (state: State = INITIAL_STATE, action: Action) => {
  switch (action.type) {
    case '@@redux/INIT':
    case '@@INIT':
      return { allowed: authenticationStorage.get('adminAllowed') }
    case 'ADMIN_SIGNIN_SUCCEEDED':
      authenticationStorage.set('adminAllowed', true)
      return { allowed: true }
    case 'ADMIN_SIGNOUT_SUCCEEDED':
    case 'ADMIN_SIGNIN_FAILED':
    case 'ADMIN_SIGNIN_ERROR':
      authenticationStorage.set('adminAllowed', false)
      return { allowed: false }
    default:
      return state
  }
}
