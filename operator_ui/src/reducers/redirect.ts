import { Reducer } from 'redux'
import { Actions } from './actions'

export interface State {
  to?: string
}

const INITIAL_STATE: State = {
  to: undefined,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case 'REDIRECT':
      return { ...state, to: action.to }
    case 'MATCH_ROUTE':
      return { ...state, to: undefined }
    default:
      return state
  }
}

export default reducer
