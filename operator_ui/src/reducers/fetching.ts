import { Reducer } from 'redux'
import { Actions, RouterActionType } from './actions'

export interface State {
  count: number
}

const INITIAL_STATE: State = {
  count: 0,
}

enum FetchPrefix {
  REQUEST = 'REQUEST_',
  RECEIVE = 'RECEIVE_',
  RESPONSE = 'RESPONSE_',
}

const reducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action: Actions,
) => {
  if (!action.type) {
    return state
  }

  if (action.type.startsWith(FetchPrefix.REQUEST)) {
    return { ...state, count: state.count + 1 }
  } else if (action.type.startsWith(FetchPrefix.RECEIVE)) {
    return { ...state, count: Math.max(state.count - 1, 0) }
  } else if (action.type.startsWith(FetchPrefix.RESPONSE)) {
    return { ...state, count: Math.max(state.count - 1, 0) }
  } else if (action.type === RouterActionType.REDIRECT) {
    return { ...state, count: 0 }
  }
  return state
}

export default reducer
