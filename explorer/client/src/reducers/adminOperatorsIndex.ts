import { Actions } from './actions'
import { Reducer } from 'redux'

export interface State {
  items?: string[]
  count?: number
}

const INITIAL_STATE: State = {}

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
      }

    default:
      return state
  }
}

export default adminOperatorsIndex
