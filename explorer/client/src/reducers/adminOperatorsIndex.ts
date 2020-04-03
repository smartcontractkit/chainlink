import { Actions } from './actions'
import { Reducer } from 'redux'

export interface State {
  items?: string[]
  count?: number
  loaded: boolean
}

const INITIAL_STATE: State = {
  loaded: false,
}

export const adminOperatorsIndex: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_OPERATORS_SUCCEEDED':
      return {
        ...state,
        items: action.data.meta.currentPageOperators.data.map(o => o.id),
        count: action.data.meta.currentPageOperators.meta.count,
        loaded: true,
      }

    case 'FETCH_ADMIN_OPERATORS_ERROR':
      return {
        ...state,
        loaded: true,
      }

    default:
      return state
  }
}

export default adminOperatorsIndex
