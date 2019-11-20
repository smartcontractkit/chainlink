import * as jsonapi from '@chainlink/json-api-client'
import reducer from 'reducers'
import { MATCH_ROUTE, RECEIVE_SIGNIN_FAIL, NOTIFY_SUCCESS } from 'actions'

describe('reducers/notifications', () => {
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
        errors: ['Your email or password is incorrect. Please try again'],
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
        component,
        props,
      }
      const state = reducer(previousState, action)

      expect(state.notifications).toEqual({
        errors: [],
        successes: [{ component, props }],
        currentUrl: undefined,
      })
    })
  })

  describe('NOTIFY_ERROR', () => {
    const component = () => {}
    const jsonApiErrors: jsonapi.ErrorItem[] = [
      { detail: 'Error 1', status: 400 },
      { detail: 'Error 2', status: 400 },
    ]
    const error = { errors: jsonApiErrors }
    const action = {
      type: 'NOTIFY_ERROR',
      component,
      error,
    }

    it('adds a component notification for each JSON-API error', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
        },
      }

      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual([
        { component, props: { msg: 'Error 1' } },
        { component, props: { msg: 'Error 2' } },
      ])
    })

    it('clears successes', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [{}],
        },
      }

      const state = reducer(previousState, action)

      expect(state.notifications.successes).toEqual([])
    })
  })

  describe('NOTIFY_ERROR_MSG', () => {
    it('adds a notification for a single error message', () => {
      const previousState = {
        notifications: {
          errors: [],
          successes: [],
        },
      }

      const action = { type: 'NOTIFY_ERROR_MSG', msg: 'Single Error' }
      const state = reducer(previousState, action)

      expect(state.notifications.errors).toEqual(['Single Error'])
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
      const action = {
        type: 'NOTIFY_ERROR_MSG',
        component,
        error,
      }
      const state = reducer(previousState, action)

      expect(state.notifications.successes).toEqual([])
    })
  })
})
