import { REDIRECT } from 'actions'

const initialState = {
  count: 0
}

export default (state = initialState, action = {}) => {
  if (/^REQUEST_/.test(action.type)) {
    return {
      ...state,
      count: state.count + 1
    }
  } else if (/^RECEIVE_/.test(action.type)) {
    return {
      ...state,
      count: Math.max(state.count - 1, 0)
    }
  } else if (/^RESPONSE_/.test(action.type)) {
    return {
      ...state,
      count: Math.max(state.count - 1, 0)
    }
  } else if (action.type === REDIRECT) {
    return {
      ...state,
      count: 0
    }
  }

  return state
}
