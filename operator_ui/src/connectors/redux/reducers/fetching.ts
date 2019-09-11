import { RouterActionType } from 'actions'
import { Action } from 'redux'

export interface State {
  count: number
}

enum FetchPrefix {
  REQUEST = 'REQUEST_',
  RECEIVE = 'RECEIVE_',
  RESPONSE = 'RESPONSE_',
}

const initialState: State = {
  count: 0,
}

export default (state: State = initialState, action: Action<String>) => {
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
