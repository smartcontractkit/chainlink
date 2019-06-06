import reducer from '../../../src/connectors/redux/reducers'
import {
  MATCH_ROUTE,
  RECEIVE_SIGNIN_FAIL,
  NOTIFY_SUCCESS,
  NOTIFY_ERROR,
} from '../../../src/actions'

describe('connectors/reducers/notifications', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.notifications).toEqual({
      errors: [],
      successes: [],
      currentUrl: undefined,
    })
  })

  describe('explorer status', () => {
    it('adds a notification when the explorer status is "not_connected"', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
          currentUrl: undefined,
        },
      }

      const cookie =
        'explorer=%7B%22status%22%3A%22not_connected%22%2C%22url%22%3A%22ws%3A%2F%2Flocalhost%3A8081%22%7D'
      const action = { type: 'ANY', cookie: cookie }
      let state = reducer(previousState, action)

      expect(state.notifications.errors.length).toEqual(1)
      expect(state.notifications.errors[0].props.msg).toEqual(
        "Can't connect to explorer: ws://localhost:8081",
      )
    })

    it('adds a notification when the explorer status cant be parsed', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
          currentUrl: undefined,
        },
      }

      const cookie = 'explorer=status'
      const action = { type: 'ANY', cookie: cookie }
      let state = reducer(previousState, action)

      expect(state.notifications.errors.length).toEqual(1)
      expect(state.notifications.errors[0].props.msg).toEqual(
        'Invalid explorer status',
      )
    })

    it('adds an extra help message when the protocol is not a websocket', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
          currentUrl: undefined,
        },
      }

      const cookie =
        'explorer=%7B%22status%22%3A%22not_connected%22%2C%22url%22%3A%22http%3A%2F%2Flocalhost%3A8081%22%7D'
      const action = { type: 'ANY', cookie: cookie }
      let state = reducer(previousState, action)

      expect(state.notifications.errors.length).toEqual(1)
      expect(state.notifications.errors[0].props.msg).toEqual(
        "Can't connect to explorer: http://localhost:8081. You must use a websocket.",
      )
    })
  })

  describe('MATCH_ROUTE', () => {
    it('clears errors when currentUrl changes', () => {
      const previousState = {
        notifications: {
          errors: [{ detail: 'error 1' }],
          successes: [{ id: '123' }],
          currentUrl: undefined,
        },
      }

      const sameUrlAction = { type: MATCH_ROUTE, match: { url: undefined } }
      let state = reducer(previousState, sameUrlAction)

      expect(state.notifications).toEqual({
        errors: [{ detail: 'error 1' }],
        successes: [{ id: '123' }],
        currentUrl: undefined,
      })

      const changedUrlAction = { type: MATCH_ROUTE, match: { url: '/' } }
      state = reducer(previousState, changedUrlAction)
      expect(state.notifications).toEqual({
        errors: [],
        successes: [],
        currentUrl: '/',
      })
    })
  })

  describe('RECEIVE_SIGNIN_FAIL', () => {
    it('adds a failure', () => {
      const action = { type: RECEIVE_SIGNIN_FAIL }
      const state = reducer(undefined, action)

      expect(state.notifications).toEqual({
        errors: [
          {
            props: {
              msg: 'Your email or password is incorrect. Please try again',
            },
          },
        ],
        successes: [],
        currentUrl: undefined,
      })
    })
  })

  describe('NOTIFY_SUCCESS', () => {
    it('adds a success component and clears errors', () => {
      const previousState = {
        notifications: {
          errors: [{ detail: 'error 1' }],
          successes: [],
          currentUrl: undefined,
        },
      }

      const component = () => {}
      const props = {}
      const action = {
        type: NOTIFY_SUCCESS,
        component: component,
        props: props,
      }
      const state = reducer(previousState, action)

      expect(state.notifications).toEqual({
        errors: [],
        successes: [{ component: component, props: props }],
        currentUrl: undefined,
      })
    })
  })

  describe('NOTIFY_ERROR', () => {
    it('adds a notification for each JSON-API errors item detail', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
        },
      }

      const component = () => {}
      const error = {
        errors: [{ detail: 'Error 1' }, { detail: 'Error 2' }],
      }
      const action = { type: NOTIFY_ERROR, component: component, error: error }
      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual([
        { component: component, props: { msg: 'Error 1' } },
        { component: component, props: { msg: 'Error 2' } },
      ])
    })

    it('adds a notification for a single error message', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
        },
      }

      const component = () => {}
      const error = { message: 'Single Error' }
      const action = { type: NOTIFY_ERROR, component: component, error: error }
      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual([
        { component: component, props: { msg: 'Single Error' } },
      ])
    })

    it('adds a notification without a component when there are no errors or message attributes', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
        },
      }

      const component = () => {}
      const error = {}
      const action = { type: NOTIFY_ERROR, component: component, error: error }
      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual([{}])
    })

    it('clears successes', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [{}],
        },
      }

      const component = () => {}
      const error = {}
      const action = { type: NOTIFY_ERROR, component: component, error: error }
      const state = reducer(previousState, action)

      expect(state.notifications.successes).toEqual([])
    })
  })
})
