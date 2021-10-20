import { Reducer } from 'redux'
import pickBy from 'lodash/pickBy'
import { Actions, ResourceActionType } from './actions'

export interface State {
  items: Record<string, any>
  currentPage?: string[]
  currentJobRunsCount?: number
}

const INITIAL_STATE: State = {
  items: {},
  currentPage: undefined,
  currentJobRunsCount: undefined,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case ResourceActionType.UPSERT_JOB_RUNS: {
      const data = action.data
      const metaCurrentPage = data.meta.currentPageJobRuns

      return {
        ...state,
        items: { ...state.items, ...data.runs },
        currentPage: metaCurrentPage.data.map((r) => r.id),
        currentJobRunsCount: metaCurrentPage.meta.count,
      }
    }
    case ResourceActionType.UPSERT_RECENT_JOB_RUNS:
    case ResourceActionType.UPSERT_JOB_RUN: {
      return {
        ...state,
        items: { ...state.items, ...action.data.runs },
      }
    }
    case ResourceActionType.RECEIVE_DELETE_SUCCESS: {
      const remainingItems = pickBy(
        state.items,
        ({ attributes }) => attributes.jobId !== action.id,
      )
      return { ...state, items: remainingItems }
    }
    default:
      return state
  }
}

export default reducer
