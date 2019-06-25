import {
  MATCH_ROUTE,
  RECEIVE_SIGNIN_FAIL,
  NOTIFY_SUCCESS,
  NOTIFY_ERROR
} from 'actions'
import { set } from 'utils/storage'
import { BadRequestError } from '../../../api/errors'

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
        return Object.assign({}, state, {
          errors: [],
          successes: [],
          currentUrl: action.match.url
        })
      }

      return state
    }
    case RECEIVE_SIGNIN_FAIL: {
      return Object.assign({}, state, {
        successes: [],
        errors: [{ props: { msg: SIGN_IN_FAIL_MSG } }]
      })
    }
    case NOTIFY_SUCCESS: {
      const success = {
        component: action.component,
        props: action.props
      }
      if (success.props.data && success.props.data.type === 'specs')
        set('persistSpec', {})
      else if (typeof success.props.url === 'string') set('persistBridge', {})
      return Object.assign({}, state, {
        successes: [success],
        errors: []
      })
    }
    case NOTIFY_ERROR: {
      const { component, error } = action
      let errorNotifications

      if (error.errors) {
        errorNotifications = error.errors.map(e => ({
          component: component,
          props: { msg: e.detail }
        }))
      } else if (error.message) {
        errorNotifications = [
          {
            component: component,
            props: { msg: error.message }
          }
        ]
      } else {
        errorNotifications = [error]
      }
      if (error instanceof BadRequestError) set('persistBridge', {})

      return Object.assign({}, state, {
        successes: [],
        errors: errorNotifications
      })
    }
    default:
      return state
  }
}
