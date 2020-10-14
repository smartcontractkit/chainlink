import { Reducer } from 'redux'
import { Actions, ResourceActionType } from './actions'

export interface State {
  items: Record<string, object>
  currentPage: string[]
  count: number
}

const INITIAL_STATE: State = {
  items: {},
  currentPage: [],
  count: 0,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case ResourceActionType.UPSERT_BRIDGES: {
      const { bridges, meta } = action.data
      const items = { ...state.items, ...bridges }
      const currentPage = meta.currentPageBridges.data.map((b) => b.id)
      const count = meta.currentPageBridges.meta.count

      return {
        ...state,
        items,
        currentPage,
        count,
      }
    }
    case ResourceActionType.UPSERT_BRIDGE: {
      const items = { ...state.items, ...action.data.bridges }

      return {
        ...state,
        items,
      }
    }
    default:
      return state
  }
}

export default reducer
