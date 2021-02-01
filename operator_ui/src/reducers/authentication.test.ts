import { partialAsFull } from 'support/test-helpers/partialAsFull'
import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  ReceiveSigninSuccessAction,
  ReceiveSigninErrorAction,
  ReceiveSignoutErrorAction,
  ReceiveSigninFailAction,
  ReceiveSignoutSuccessAction,
  AuthActionType,
} from '../../src/reducers/actions'
import { getAuthentication } from '../../src/utils/storage'

describe('reducers/authentication', () => {
  const successAction: ReceiveSigninSuccessAction = {
    type: AuthActionType.RECEIVE_SIGNIN_SUCCESS,
    authenticated: true,
  }
  const signInErrorAction: ReceiveSigninErrorAction = {
    type: AuthActionType.RECEIVE_SIGNIN_ERROR,
    errors: ['error 1'],
  }
  const signOutErrorAction: ReceiveSignoutErrorAction = {
    type: AuthActionType.RECEIVE_SIGNOUT_ERROR,
    errors: ['error 2'],
  }

  beforeEach(() => {
    localStorage.clear()
  })

  describe('RECEIVE_SIGNIN_ERROR', () => {
    it('saves allowed false to local storage', () => {
      const action = partialAsFull<ReceiveSigninErrorAction>({
        type: AuthActionType.RECEIVE_SIGNIN_ERROR,
      })

      const state = reducer(INITIAL_STATE, successAction)
      expect(getAuthentication()).toEqual({ allowed: true })

      reducer(state, action)
      expect(getAuthentication()).toEqual({ allowed: false })
    })

    it('assigns errors', () => {
      const state = reducer(INITIAL_STATE, signInErrorAction)

      expect(state.authentication.errors).toEqual(['error 1'])
    })
  })

  describe('RECEIVE_SIGNOUT_ERROR', () => {
    it('saves allowed false to local storage', () => {
      const state = reducer(INITIAL_STATE, successAction)
      expect(getAuthentication()).toEqual({ allowed: true })

      reducer(state, signOutErrorAction)
      expect(getAuthentication()).toEqual({ allowed: false })
    })

    it('assigns errors', () => {
      const state = reducer(INITIAL_STATE, signOutErrorAction)
      expect(state.authentication.errors).toEqual(['error 2'])
    })
  })

  describe('RECEIVE_SIGNIN_SUCCESS', () => {
    it('assigns allowed and saves it in local storage', () => {
      const state = reducer(INITIAL_STATE, successAction)
      expect(state.authentication.allowed).toEqual(true)
      expect(getAuthentication()).toEqual({ allowed: true })
    })
  })

  describe('RECEIVE_SIGNIN_FAIL', () => {
    const failAction: ReceiveSigninFailAction = {
      type: AuthActionType.RECEIVE_SIGNIN_FAIL,
    }

    it('clears authentication errors', () => {
      let state = reducer(INITIAL_STATE, signInErrorAction)
      expect(state.authentication.errors).toEqual(['error 1'])

      state = reducer(state, failAction)
      expect(state.authentication.allowed).toEqual(false)
      expect(state.authentication.errors).toEqual([])
    })

    it('saves allowed false to local storage', () => {
      const state = reducer(INITIAL_STATE, successAction)
      expect(getAuthentication()).toEqual({ allowed: true })

      reducer(state, failAction)
      expect(getAuthentication()).toEqual({ allowed: false })
    })
  })

  describe('RECEIVE_SIGNOUT_SUCCESS', () => {
    it('assigns allowed and saves it to local storage', () => {
      const action: ReceiveSignoutSuccessAction = {
        type: AuthActionType.RECEIVE_SIGNOUT_SUCCESS,
        authenticated: false,
      }
      let state = reducer(INITIAL_STATE, successAction)
      expect(getAuthentication()).toEqual({ allowed: true })

      state = reducer(state, action)
      expect(state.authentication.allowed).toEqual(false)
      expect(getAuthentication()).toEqual({ allowed: false })
    })
  })
})
