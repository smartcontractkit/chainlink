import { Reducer } from 'redux'
import { Actions, ResourceActionType } from './actions'

export interface State {
  currentPage?: string[]
  count: number
}

const INITIAL_STATE: State = {
  currentPage: undefined,
  count: 0,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case ResourceActionType.UPSERT_TRANSACTIONS: {
      const data = action.data
      const metaCurrentPage = data.meta.currentPageTransactions
      const currentPage = metaCurrentPage.data.map((t) => t.id)
      const count = metaCurrentPage.meta.count

      return {
        ...state,
        currentPage,
        count,
      }
    }
    default:
      return state
  }
}

export default reducer
