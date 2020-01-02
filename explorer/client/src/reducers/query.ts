import { Actions } from './actions'
import { Reducer } from 'redux'

export type Query = string | undefined
export interface State {
  query?: Query
}

const INITIAL_STATE = { query: undefined }

export const queryReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'QUERY_UPDATED': {
      return { ...state, query: action.data }
    }

    default:
      return state
  }
}

export default queryReducer
