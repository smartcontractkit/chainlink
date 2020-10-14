import { Reducer } from 'redux'
import pickBy from 'lodash/pickBy'
import { Actions, ResourceActionType } from './actions'
import { JobSpec } from 'operator_ui'

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
    case ResourceActionType.UPSERT_JOBS: {
      const data = action.data
      const currentPage = data.meta.currentPageJobs.data.map((j) => j.id)
      const count = data.meta.currentPageJobs.meta.count
      const items = { ...state.items, ...action.data.specs }

      return {
        ...state,
        currentPage,
        count,
        items,
      }
    }
    case ResourceActionType.UPSERT_RECENTLY_CREATED_JOBS: {
      const data = action.data
      const recentlyCreated = data.meta.recentlyCreatedJobs.data.map(
        (j) => j.id,
      )
      const items = { ...state.items, ...data.specs }

      return {
        ...state,
        recentlyCreated,
        items,
      }
    }
    case ResourceActionType.UPSERT_JOB: {
      const items = { ...state.items, ...action.data.specs }

      return { ...state, items }
    }
    case ResourceActionType.RECEIVE_DELETE_SUCCESS: {
      const items = pickBy(state.items, (i) => i.id !== action.id)

      return { ...state, items }
    }
    case ResourceActionType.DELETE_JOB_SPEC_ERROR: {
      const resource = state.items[action.data.jobSpecID]
      const attributes: JobSpec = resource.attributes
      const newErrors = attributes.errors.filter(
        (error) => error.id != action.data.id,
      )
      const newAttributes = { ...attributes, errors: newErrors }
      const newResource = { ...resource, attributes: newAttributes }
      const items = { ...state.items, [action.data.jobSpecID]: newResource }

      return { ...state, items }
    }
    default:
      return state
  }
}

export default reducer
