import { parse as parseCookie } from 'cookie'
import { set } from '../../../utils/storage'
import { BadRequestError } from '../../../api/errors'

export interface State {
  errors: any[]
  successes: any[]
  currentUrl?: string
}

const initialState: State = {
  errors: [],
  successes: [],
  currentUrl: undefined,
}

export type Action =
  | { type: 'MATCH_ROUTE'; cookie?: string; match: any }
  | { type: 'RECEIVE_SIGNIN_FAIL'; cookie?: string }
  | { type: 'NOTIFY_SUCCESS'; cookie?: string; component: any; props: any }
  | { type: 'NOTIFY_ERROR'; cookie?: string; component: any; error: any }
// | { type: string; cookie?: string }

const SIGN_IN_FAIL_MSG = 'Your email or password is incorrect. Please try again'

export default function(state: State = initialState, action: Action) {
  const before = beforeCookieState(state, action)
  const after = afterCookieState(before, action)

  return after
}

function beforeCookieState(state: State, action: Action): State {
  switch (action.type) {
    case 'MATCH_ROUTE': {
      if (action.match && state.currentUrl !== action.match.url) {
        return {
          ...state,
          errors: [],
          successes: [],
          currentUrl: action.match.url,
        }
      }

      return state
    }
    case 'RECEIVE_SIGNIN_FAIL': {
      return {
        ...state,
        successes: [],
        errors: [{ props: { msg: SIGN_IN_FAIL_MSG } }],
      }
    }
    case 'NOTIFY_SUCCESS': {
      const success = {
        component: action.component,
        props: action.props,
      }
      if (success.props.data && success.props.data.type === 'specs') {
        set('persistSpec', {})
      } else if (typeof success.props.url === 'string') {
        set('persistBridge', {})
      }

      return {
        ...state,
        successes: [success],
        errors: [],
      }
    }
    case 'NOTIFY_ERROR': {
      let errorNotifications

      if (action.error.errors) {
        errorNotifications = action.error.errors.map((e: any) => ({
          component: action.component,
          props: { msg: e.detail },
        }))
      } else if (action.error.message) {
        errorNotifications = [
          {
            component: action.component,
            props: { msg: action.error.message },
          },
        ]
      } else {
        errorNotifications = [action.error]
      }
      if (action.error instanceof BadRequestError) {
        set('persistBridge', {})
      }

      return {
        ...state,
        successes: [],
        errors: errorNotifications,
      }
    }
    default:
      return state
  }
}

const NOT_CONNECTED = 'not_connected'

const hasExplorerStatus = (errors: any, msg: any) => {
  return errors.find(({ props }: any) => props && props.msg === msg)
}

function afterCookieState(state: State, action: Action): State {
  const cookies = parseCookie(
    action.cookie || (document ? document.cookie : ''),
  )
  let notification

  if (cookies.explorer) {
    try {
      const json = JSON.parse(cookies.explorer)

      if (json.status === NOT_CONNECTED) {
        let msg = `Can't connect to explorer: ${json.url}`
        if (!json.url.match(/^wss?:.+/)) {
          msg = `${msg}. You must use a websocket.`
        }
        if (!hasExplorerStatus(state.errors, msg)) {
          notification = { props: { msg } }
        }
      }
    } catch (e) {
      notification = { props: { msg: 'Invalid explorer status' } }
    }
  }

  return {
    ...state,
    errors: [...state.errors, notification].filter(n => !!n),
  }
}
