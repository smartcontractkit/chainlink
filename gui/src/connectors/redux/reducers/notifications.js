import {
  MATCH_ROUTE,
  RECEIVE_SIGNIN_FAIL,
  RECEIVE_CREATE_SUCCESS,
  RECEIVE_CREATE_ERROR,
  NOTIFY_SUCCESS,
  NOTIFY_ERROR
} from 'actions'

const initialState = {
  errors: [],
  successes: [],
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
          {errors: [], successes: [], currentUrl: action.match.url}
        )
      }

      return state
    }
    case RECEIVE_SIGNIN_FAIL: {
      return Object.assign(
        {},
        state,
        {
          successes: [],
          errors: [{detail: SIGN_IN_FAIL_MSG}]
        }
      )
    }
    case RECEIVE_CREATE_ERROR: {
      return Object.assign(
        {},
        state,
        {
          successes: [],
          errors: action.error.errors
        }
      )
    }
    case RECEIVE_CREATE_SUCCESS: {
      return Object.assign(
        {},
        state,
        {
          successes: [action.response],
          errors: []
        }
      )
    }
    case NOTIFY_SUCCESS: {
      return Object.assign(
        {},
        state,
        {
          successes: [{type: 'component', component: action.component, props: action.props}],
          errors: []
        }
      )
    }
    case NOTIFY_ERROR: {
      return Object.assign(
        {},
        state,
        {
          successes: [],
          errors: [{type: 'component', component: action.component, props: action.props}]
        }
      )
    }
    default:
      return state
  }
}
