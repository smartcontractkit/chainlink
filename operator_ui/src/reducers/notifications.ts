import { PropsWithChildren } from 'react'
import { Reducer } from 'redux'
import * as jsonapi from 'utils/json-api-client'
import {
  Actions,
  NotifyErrorAction,
  AuthActionType,
  RouterActionType,
  NotifyActionType,
} from 'reducers/actions'

export interface State {
  errors: Notification[]
  successes: Notification[]
  currentUrl?: string
}

const INITIAL_STATE: State = {
  errors: [],
  successes: [],
  currentUrl: undefined,
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case RouterActionType.MATCH_ROUTE: {
      return {
        ...INITIAL_STATE,
        currentUrl: action.pathname,
      }
    }
    case NotifyActionType.NOTIFY_SUCCESS: {
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
    case NotifyActionType.NOTIFY_SUCCESS_MSG: {
      return {
        ...state,
        successes: [action.msg],
        errors: [],
      }
    }
    case NotifyActionType.NOTIFY_ERROR: {
      const errors = action.error.errors
      const notifications = errors.map((e) =>
        buildJsonApiErrorNotification(action, e),
      )

      return {
        ...state,
        successes: [],
        errors: notifications,
      }
    }
    case NotifyActionType.NOTIFY_ERROR_MSG: {
      return {
        ...state,
        successes: [],
        errors: [action.msg],
      }
    }
    case AuthActionType.RECEIVE_SIGNIN_FAIL: {
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

export type TextNotification = string

export interface ComponentNotification {
  component: React.FC<any>
  props: PropsWithChildren<{ msg?: string }>
}

export type Notification = TextNotification | ComponentNotification

function buildJsonApiErrorNotification(
  action: NotifyErrorAction,
  e: jsonapi.ErrorItem,
): ComponentNotification {
  return {
    component: action.component,
    props: { msg: e.detail },
  }
}

export default reducer
