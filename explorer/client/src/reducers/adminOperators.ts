import { Actions } from './actions'
import { Reducer } from 'redux'
import { ChainlinkNode } from 'explorer/models'

export interface State {
  loading: boolean
  error: boolean
  items?: Record<ChainlinkNode['id'], ChainlinkNode>
}

const INITIAL_STATE: State = {
  loading: false,
  error: false,
}

export const adminOperators: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_OPERATORS_BEGIN':
      return {
        ...state,
        loading: true,
        error: false,
      }

    case 'FETCH_ADMIN_OPERATORS_SUCCEEDED':
      return {
        ...state,
        items: { ...action.data.chainlinkNodes },
        loading: false,
      }

    case 'FETCH_ADMIN_OPERATORS_ERROR':
      return {
        ...state,
        loading: false,
        error: true,
      }

    default:
      return state
  }
}

export default adminOperators
