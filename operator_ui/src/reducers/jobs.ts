import { Reducer } from 'redux'
import pickBy from 'lodash/pickBy'
import { Actions } from './actions'

export interface State {
  items: Record<string, any>
  currentPage?: string[]
  recentlyCreated?: string[]
  count: number
}

const INITIAL_STATE: State = {
  items: {},
  currentPage: undefined,
  recentlyCreated: undefined,
  count: 0,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case 'UPSERT_JOBS': {
      const data = action.data
      const currentPage = data.meta.currentPageJobs.data.map(j => j.id)
      const count = data.meta.currentPageJobs.meta.count
      const items = { ...state.items, ...action.data.specs }

      return {
        ...state,
        currentPage,
        count,
        items,
      }
    }
    case 'UPSERT_RECENTLY_CREATED_JOBS': {
      const data = action.data
      const recentlyCreated = data.meta.recentlyCreatedJobs.data.map(j => j.id)
      const items = { ...state.items, ...data.specs }

      return {
        ...state,
        recentlyCreated,
        items,
      }
    }
    case 'UPSERT_JOB': {
      const items = { ...state.items, ...action.data.specs }

      return { ...state, items }
    }
    case 'RECEIVE_DELETE_SUCCESS': {
      const items = pickBy(state.items, i => i.id !== action.id)

      return { ...state, items }
    }
    default:
      return state
  }
}

export default reducer
