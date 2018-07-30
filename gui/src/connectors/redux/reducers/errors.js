import {
  MATCH_ROUTE,
  RECEIVE_SESSION_FAIL
} from 'actions'

const initialState = {
  messages: [],
  currentUrl: null
}
const SIGN_IN_FAIL_MSG = 'Your email or password are incorrect. Please try again'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case MATCH_ROUTE: {
      if (action.match && state.currentUrl !== action.match.url) {
        return Object.assign(
          {},
          state,
          {messages: [], currentUrl: action.match.url}
        )
      }

      return state
    }
    case RECEIVE_SESSION_FAIL:
      return Object.assign(
        {},
        state,
        {messages: [SIGN_IN_FAIL_MSG]}
      )
    default:
      return state
  }
}
