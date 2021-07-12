import { createStore, applyMiddleware, Middleware } from 'redux'
import reducer, { AppState } from '../../src/reducers'
import { MatchRouteAction, RouterActionType } from '../../src/reducers/actions'
import { createExplorerConnectionMiddleware } from '../../src/middleware/explorerConnection'

const WS_ERROR_COOKIE =
  'explorer=%7B%22status%22%3A%22error%22%2C%22url%22%3A%22ws%3A%2F%2Flocalhost%3A8081%22%7D'
const HTTP_ERROR_COOKIE =
  'explorer=%7B%22status%22%3A%22error%22%2C%22url%22%3A%22http%3A%2F%2Flocalhost%3A8081%22%7D'
const PARSE_ERROR_COOKIE = 'explorer=status'

describe('middleware/explorerConnection', () => {
  it('adds a notification after MATCH_ROUTE when the explorer status is "error"', () => {
    const middleware: Middleware[] = [
      createExplorerConnectionMiddleware(WS_ERROR_COOKIE),
    ]
    const store = createStore(reducer, applyMiddleware(...middleware))
    const action: MatchRouteAction = {
      type: RouterActionType.MATCH_ROUTE,
      pathname: '/',
    }

    store.dispatch(action)
    const state: AppState = store.getState()

    expect(state.notifications.errors.length).toEqual(1)
    expect(state.notifications.errors[0]).toEqual(
      "Can't connect to explorer: ws://localhost:8081",
    )
  })

  it('adds a notification when the explorer status cant be parsed', () => {
    const middleware: Middleware[] = [
      createExplorerConnectionMiddleware(PARSE_ERROR_COOKIE),
    ]
    const store = createStore(reducer, applyMiddleware(...middleware))
    const action: MatchRouteAction = {
      type: RouterActionType.MATCH_ROUTE,
      pathname: '/',
    }

    store.dispatch(action)
    const state: AppState = store.getState()

    expect(state.notifications.errors.length).toEqual(1)
    expect(state.notifications.errors[0]).toEqual('Invalid explorer status')
  })

  it('adds an extra help message when the protocol is not a websocket', () => {
    const middleware: Middleware[] = [
      createExplorerConnectionMiddleware(HTTP_ERROR_COOKIE),
    ]
    const store = createStore(reducer, applyMiddleware(...middleware))
    const action: MatchRouteAction = {
      type: RouterActionType.MATCH_ROUTE,
      pathname: '/',
    }

    store.dispatch(action)
    const state: AppState = store.getState()

    expect(state.notifications.errors.length).toEqual(1)
    expect(state.notifications.errors[0]).toEqual(
      "Can't connect to explorer: http://localhost:8081. You must use a websocket.",
    )
  })

  it('doesnt add a notification on sign in', () => {
    const middleware: Middleware[] = [
      createExplorerConnectionMiddleware(WS_ERROR_COOKIE),
    ]
    const store = createStore(reducer, applyMiddleware(...middleware))
    const action: MatchRouteAction = {
      type: RouterActionType.MATCH_ROUTE,
      pathname: '/signin',
    }

    store.dispatch(action)
    const state: AppState = store.getState()

    expect(state.notifications.errors.length).toEqual(0)
  })
})
