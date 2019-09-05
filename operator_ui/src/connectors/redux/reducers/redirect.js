import { RouterActionType } from 'actions'

const initialState = {
  to: null
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case RouterActionType.REDIRECT:
      return Object.assign({}, state, { to: action.to })
    case RouterActionType.MATCH_ROUTE:
      return Object.assign({}, state, { to: null })
    default:
      return state
  }
}
