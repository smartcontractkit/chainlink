import * as jsonapi from '@chainlink/json-api-client'
import { Actions } from './actions'
import { Reducer } from 'redux'

export interface State {
  errors: string[]
}

const INITIAL_STATE: State = { errors: [] }

const notificationsReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
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

export default notificationsReducer
