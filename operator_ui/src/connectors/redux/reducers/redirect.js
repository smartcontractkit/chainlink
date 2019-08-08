import { REDIRECT, MATCH_ROUTE } from 'actions'

const initialState = {
  to: null
}

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case REDIRECT:
      return {
        ...state,
        to: action.to
      }
    case MATCH_ROUTE:
      return {
        ...state,
        to: null
      }
    default:
      return state
  }
}
