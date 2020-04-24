import { Actions } from './actions'
import { Reducer } from 'redux'

export interface State {
  items?: string[]
  count?: number
}

const INITIAL_STATE: State = {}

export const adminHeadsIndex: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_HEADS_SUCCEEDED':
      return {
        items: action.data.meta.currentPageHeads.data.map(o => o.id),
        count: action.data.meta.currentPageHeads.meta.count,
      }

    default:
      return state
  }
}

export default adminHeadsIndex
