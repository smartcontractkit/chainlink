import { Middleware } from 'redux'
import { parse } from 'cookie'
import { AppState } from 'reducers'
import { Actions } from 'reducers/actions'

/**
 * Create a redux middleware responsible for parsing the explorer status
 * cookie after every MATCH_ROUTE action
 */
export function createExplorerConnectionMiddleware(
  cookie: string = (global.document && global.document.cookie) || '',
): Middleware {
  const explorerConnectionMiddleware: Middleware = store => next => (
    action: Actions,
  ) => {
    // dispatch original action right away
    next(action)

    const state: AppState = store.getState()
    if (
      action.type === 'MATCH_ROUTE' &&
      state.notifications.currentUrl !== '/signin'
    ) {
      const cookies = parse(cookie)

      if (cookies.explorer) {
        try {
          const json = JSON.parse(cookies.explorer)

          if (isErrorStatus(json.status)) {
            const msg = formatMsg(json.url)
            next({ type: 'NOTIFY_ERROR_MSG', msg })
          }
        } catch {
          next({ type: 'NOTIFY_ERROR_MSG', msg: 'Invalid explorer status' })
        }
      }
    }
  }

  return explorerConnectionMiddleware
}

function formatMsg(url: string) {
  const msg = `Can't connect to explorer: ${url}`

  if (url.match(/^wss?:.+/)) {
    return msg
  }

  return `${msg}. You must use a websocket.`
}

function isErrorStatus(status: string) {
  return status === 'error'
}
