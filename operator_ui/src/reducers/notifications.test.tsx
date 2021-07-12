import * as React from 'react'
import * as jsonapi from 'utils/json-api-client'
import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  AuthActionType,
  NotifyActionType,
  RouterActionType,
  ReceiveSigninFailAction,
  NotifySuccessAction,
  NotifyErrorAction,
  NotifyErrorMsgAction,
  MatchRouteAction,
} from '../../src/reducers/actions'

describe('reducers/notifications', () => {
  const props = {}
  const component: React.FC = () => {
    return <></>
  }

  describe('RECEIVE_SIGNIN_FAIL', () => {
    it('adds an error', () => {
      const action: ReceiveSigninFailAction = {
        type: AuthActionType.RECEIVE_SIGNIN_FAIL,
      }
      const state = reducer(INITIAL_STATE, action)

      expect(state.notifications.errors).toEqual([
        'Your email or password is incorrect. Please try again',
      ])
    })
  })

  describe('NOTIFY_SUCCESS', () => {
    const action: NotifySuccessAction = {
      type: NotifyActionType.NOTIFY_SUCCESS,
      component,
      props,
    }

    it('adds a success component', () => {
      const state = reducer(INITIAL_STATE, action)
      expect(state.notifications.successes).toEqual([{ component, props }])
    })

    it('clears errors', () => {
      const errorAction: ReceiveSigninFailAction = {
        type: AuthActionType.RECEIVE_SIGNIN_FAIL,
      }

      let state = reducer(INITIAL_STATE, errorAction)
      expect(state.notifications.errors.length).toEqual(1)

      state = reducer(state, action)
      expect(state.notifications.errors.length).toEqual(0)
    })
  })

  describe('NOTIFY_ERROR', () => {
    const jsonApiErrors: jsonapi.ErrorItem[] = [
      { detail: 'Error 1', status: 400 },
      { detail: 'Error 2', status: 400 },
    ]
    const error = { errors: jsonApiErrors }
    const action: NotifyErrorAction = {
      type: NotifyActionType.NOTIFY_ERROR,
      component,
      error,
    }

    it('adds a component notification for each JSON-API error', () => {
      const state = reducer(INITIAL_STATE, action)

      expect(state.notifications.errors).toEqual([
        { component, props: { msg: 'Error 1' } },
        { component, props: { msg: 'Error 2' } },
      ])
    })

    it('clears successes', () => {
      const successAction: NotifySuccessAction = {
        type: NotifyActionType.NOTIFY_SUCCESS,
        props,
        component,
      }

      let state = reducer(INITIAL_STATE, successAction)
      expect(state.notifications.successes.length).toEqual(1)

      state = reducer(state, action)
      expect(state.notifications.successes.length).toEqual(0)
    })
  })

  describe('NOTIFY_ERROR_MSG', () => {
    const action: NotifyErrorMsgAction = {
      type: NotifyActionType.NOTIFY_ERROR_MSG,
      msg: 'Single Error',
    }

    it('adds a notification for a single error message', () => {
      const state = reducer(INITIAL_STATE, action)
      expect(state.notifications.errors).toEqual(['Single Error'])
    })

    it('clears successes', () => {
      const successAction: NotifySuccessAction = {
        type: NotifyActionType.NOTIFY_SUCCESS,
        props,
        component,
      }

      let state = reducer(INITIAL_STATE, successAction)
      expect(state.notifications.successes.length).toEqual(1)

      state = reducer(state, action)
      expect(state.notifications.successes.length).toEqual(0)
    })
  })

  describe('MATCH_ROUTE', () => {
    const changeUrlAction: MatchRouteAction = {
      type: RouterActionType.MATCH_ROUTE,
      pathname: '/to',
    }

    it('clears errors when currentUrl changes', () => {
      const errorAction: ReceiveSigninFailAction = {
        type: AuthActionType.RECEIVE_SIGNIN_FAIL,
      }

      let state = reducer(INITIAL_STATE, errorAction)
      expect(state.notifications.errors.length).toEqual(1)
      expect(state.notifications.successes.length).toEqual(0)
      expect(state.notifications.currentUrl).toBeUndefined()

      state = reducer(INITIAL_STATE, changeUrlAction)
      expect(state.notifications.errors.length).toEqual(0)
      expect(state.notifications.successes.length).toEqual(0)
      expect(state.notifications.currentUrl).toEqual('/to')
    })

    it('clears success when currentUrl changes', () => {
      const successAction: NotifySuccessAction = {
        type: NotifyActionType.NOTIFY_SUCCESS,
        props,
        component,
      }

      let state = reducer(INITIAL_STATE, successAction)
      expect(state.notifications.successes.length).toEqual(1)
      expect(state.notifications.currentUrl).toBeUndefined()

      state = reducer(INITIAL_STATE, changeUrlAction)
      expect(state.notifications.successes.length).toEqual(0)
      expect(state.notifications.currentUrl).toEqual('/to')
    })
  })
})
