import { Actions } from './actions'
import { Reducer } from 'redux'

export interface State {
  loaded: boolean
  items?: string[]
  count?: number
}

const INITIAL_STATE: State = { loaded: false }

export const adminHeadsIndex: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_HEADS_SUCCEEDED':
      return {
        items: action.data.meta.currentPageHeads.data.map(o => o.id),
        count: action.data.meta.currentPageHeads.meta.count,
        loaded: true,
      }
    case 'FETCH_ADMIN_HEADS_ERROR':
      return {
        ...state,
        ...{ loaded: true },
      }
    default:
      return state
  }
}

export default adminHeadsIndex
