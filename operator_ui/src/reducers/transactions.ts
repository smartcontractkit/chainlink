import { Reducer } from 'redux'
import { Actions, ResourceActionType } from './actions'

export interface State {
  items: Record<string, any>
}

const INITIAL_STATE: State = {
  items: {},
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case ResourceActionType.UPSERT_TRANSACTIONS: {
      const items = { ...state.items, ...action.data.transactions }
      return { ...state, items }
    }
    case ResourceActionType.UPSERT_TRANSACTION: {
      const items = { ...state.items, ...action.data.transactions }
      return { ...state, items }
    }
    default:
      return state
  }
}

export default reducer
