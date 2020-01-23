import { Reducer } from 'redux'
import { Actions, RouterActionType } from './actions'

export interface State {
  to?: string
}

const INITIAL_STATE: State = {
  to: undefined,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case RouterActionType.REDIRECT:
      return { ...state, to: action.to }
    case RouterActionType.MATCH_ROUTE:
      return { ...state, to: undefined }
    default:
      return state
  }
}

export default reducer
