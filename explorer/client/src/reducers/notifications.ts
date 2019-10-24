import * as jsonapi from '@chainlink/json-api-client'
import { Actions } from './actions'

export interface State {
  errors: string[]
}

const INITIAL_STATE: State = { errors: [] }

export default (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'FETCH_ADMIN_SIGNIN_ERROR':
      if (isUnauthorized(action.error)) {
        return { errors: ['Invalid username and password.'] }
      }

      return state
    case 'NOTIFY_ERROR':
      return { errors: [action.text] }
    default:
      return state
  }
}

function isUnauthorized(error: Error) {
  return error instanceof jsonapi.AuthenticationError
}
