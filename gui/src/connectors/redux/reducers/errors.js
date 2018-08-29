import {
  MATCH_ROUTE,
  RECEIVE_SIGNIN_FAIL,
  RECEIVE_CREATE_ERROR
} from 'actions'

const initialState = {
  errors: [],
  currentUrl: null
}
const SIGN_IN_FAIL_MSG = 'Your email or password is incorrect. Please try again'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case MATCH_ROUTE: {
      if (action.match && state.currentUrl !== action.match.url) {
        return Object.assign(
          {},
          state,
          {errors: [], currentUrl: action.match.url}
        )
      }

      return state
    }
    case RECEIVE_SIGNIN_FAIL:
      return Object.assign(
        {},
        state,
        {errors: [{detail: SIGN_IN_FAIL_MSG}]}
      )
    case RECEIVE_CREATE_ERROR:
      return Object.assign(
        {},
        state,
        {errors: action.error.errors}
      )
    default:
      return state
  }
}
