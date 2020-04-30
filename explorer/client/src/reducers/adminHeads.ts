import { Actions } from './actions'
import { Reducer } from 'redux'
import { Head } from 'explorer/models'

export interface State {
  loading: boolean
  error: boolean
  items?: Record<Head['id'], Head>
}

const INITIAL_STATE: State = {
  loading: false,
  error: false,
}

export const adminHeads: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_HEADS_BEGIN':
      return {
        ...state,
        loading: true,
        error: false,
      }

    case 'FETCH_ADMIN_HEADS_SUCCEEDED':
      return {
        ...state,
        items: { ...action.data.heads },
        loading: false,
      }

    case 'FETCH_ADMIN_HEADS_ERROR':
      return {
        ...state,
        loading: false,
        error: true,
      }

    default:
      return state
  }
}

export default adminHeads
