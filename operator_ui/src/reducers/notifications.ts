import * as jsonapi from '@chainlink/json-api-client'
import { Actions, NotifyErrorAction } from 'reducers/actions'

export interface State {
  errors: Notification[]
  successes: Notification[]
  currentUrl?: string
}

const initialState: State = {
  errors: [],
  successes: [],
  currentUrl: undefined,
}

export type TextNotification = string

export interface ComponentNotification {
  component: React.FC<any>
  props: any
}

export type Notification = TextNotification | ComponentNotification

export default function(state: State = initialState, action: Actions) {
  switch (action.type) {
    case 'MATCH_ROUTE': {
      if (action.match && state.currentUrl !== action.match.url) {
        return {
          ...initialState,
          currentUrl: action.match.url,
        }
      }

      return state
    }
    case 'NOTIFY_SUCCESS': {
      const success: ComponentNotification = {
        component: action.component,
        props: action.props,
      }

      return {
        ...state,
        successes: [success],
        errors: [],
      }
    }
    case 'NOTIFY_SUCCESS_MSG': {
      return {
        ...state,
        successes: [action.msg],
        errors: [],
      }
    }
    case 'NOTIFY_ERROR': {
      const notifications = action.error.errors.map(e =>
        buildJsonApiErrorNotification(action, e),
      )

      return {
        ...state,
        successes: [],
        errors: notifications,
      }
    }
    case 'NOTIFY_ERROR_MSG': {
      return {
        ...state,
        successes: [],
        errors: [action.msg],
      }
    }
    case 'RECEIVE_SIGNIN_FAIL': {
      return {
        ...state,
        successes: [],
        errors: ['Your email or password is incorrect. Please try again'],
      }
    }
    default:
      return state
  }
}

function buildJsonApiErrorNotification(
  action: NotifyErrorAction,
  e: jsonapi.ErrorItem,
): ComponentNotification {
  return {
    component: action.component,
    props: { msg: e.detail },
  }
}
