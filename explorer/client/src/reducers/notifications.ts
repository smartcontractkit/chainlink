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
    case 'FETCH_ADMIN_SIGNIN_ERROR': {
      const errors = action.errors.map(e => {
        if (e.status === 500) {
          return 'Error processing your request. Please ensure your connection is active and try again'
        } else if (e.status === 401) {
          return 'Invalid username and password'
        }
        return e.detail
      })
      return { errors }
    }
    default:
      return state
  }
}

export default notificationsReducer
