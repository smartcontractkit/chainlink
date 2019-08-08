import reducer from 'connectors/redux/reducers'
import {
  MATCH_ROUTE,
  RECEIVE_SIGNIN_FAIL,
  NOTIFY_SUCCESS,
  NOTIFY_ERROR
} from 'actions'

describe('connectors/reducers/notifications', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.notifications).toEqual({
      errors: [],
      successes: [],
      currentUrl: null
    })
  })

  it('MATCH_ROUTE clears errors when currentUrl changes', () => {
    const previousState = {
      notifications: {
        errors: [{ detail: 'error 1' }],
        successes: [{ id: '123' }],
        currentUrl: null
      }
    }

    const sameUrlAction = { type: MATCH_ROUTE, match: { url: null } }
    let state = reducer(previousState, sameUrlAction)

    expect(state.notifications).toEqual({
      errors: [{ detail: 'error 1' }],
      successes: [{ id: '123' }],
      currentUrl: null
    })

    const changedUrlAction = { type: MATCH_ROUTE, match: { url: '/' } }
    state = reducer(previousState, changedUrlAction)
    expect(state.notifications).toEqual({
      errors: [],
      successes: [],
      currentUrl: '/'
    })
  })

  it('RECEIVE_SIGNIN_FAIL adds a failure', () => {
    const action = { type: RECEIVE_SIGNIN_FAIL }
    const state = reducer(undefined, action)

    expect(state.notifications).toEqual({
      errors: [
        {
          props: {
            msg: 'Your email or password is incorrect. Please try again'
          }
        }
      ],
      successes: [],
      currentUrl: null
    })
  })

  it('NOTIFY_SUCCESS adds a success component and clears errors', () => {
    const previousState = {
      notifications: {
        errors: [{ detail: 'error 1' }],
        successes: [],
        currentUrl: null
      }
    }

    const component = () => {}
    const props = {}
    const action = { type: NOTIFY_SUCCESS, component, props }
    const state = reducer(previousState, action)

    expect(state.notifications).toEqual({
      errors: [],
      successes: [{ component, props }],
      currentUrl: null
    })
  })

  describe('NOTIFY_ERROR', () => {
    it('adds a notification for each JSON-API errors item detail', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: []
        }
      }

      const component = () => {}
      const error = {
        errors: [{ detail: 'Error 1' }, { detail: 'Error 2' }]
      }
      const action = { type: NOTIFY_ERROR, component, error }
      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual([
        { component, props: { msg: 'Error 1' } },
        { component, props: { msg: 'Error 2' } }
      ])
    })

    it('adds a notification for a single error message', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: []
        }
      }

      const component = () => {}
      const error = { message: 'Single Error' }
      const action = { type: NOTIFY_ERROR, component, error }
      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual([
        { component, props: { msg: 'Single Error' } }
      ])
    })

    it('adds a notification without a component when there are no errors or message attributes', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: []
        }
      }

      const component = () => {}
      const error = {}
      const action = { type: NOTIFY_ERROR, component, error }
      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual([{}])
    })

    it('clears successes', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [{}]
        }
      }

      const component = () => {}
      const error = {}
      const action = { type: NOTIFY_ERROR, component, error }
      const state = reducer(previousState, action)

      expect(state.notifications.successes).toEqual([])
    })
  })
})
