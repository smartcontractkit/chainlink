import reducer from 'connectors/redux/reducers'
import { get as getSessionStorage } from 'utils/sessionStorage'
import {
  REQUEST_SIGNIN,
  RECEIVE_SIGNIN_SUCCESS,
  RECEIVE_SIGNIN_FAIL,
  RECEIVE_SIGNIN_ERROR,
  REQUEST_SIGNOUT,
  RECEIVE_SIGNOUT_SUCCESS,
  RECEIVE_SIGNOUT_ERROR
} from 'actions'

describe('session reducer', () => {
  beforeEach(() => {
    global.localStorage.clear()
  })

  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.session).toEqual({
      fetching: false,
      authenticated: false,
      errors: [],
      networkError: false
    })
  })

  it('REQUEST_SIGNIN starts fetching and disables the network error', () => {
    const action = {type: REQUEST_SIGNIN}
    const state = reducer(undefined, action)

    expect(state.session.fetching).toEqual(true)
    expect(state.session.networkError).toEqual(false)
  })

  describe('RECEIVE_SIGNIN_SUCCESS', () => {
    it('stops fetching and assigns authenticated', () => {
      const previousState = {
        session: {
          fetching: true,
          networkError: true
        }
      }
      const action = {type: RECEIVE_SIGNIN_SUCCESS, authenticated: true}
      const state = reducer(previousState, action)

      expect(state.session.authenticated).toEqual(true)
      expect(state.session.fetching).toEqual(false)
      expect(state.session.networkError).toEqual(false)
    })

    it('saves authenticated true to local storage', () => {
      const action = {type: RECEIVE_SIGNIN_SUCCESS, authenticated: true}
      reducer(undefined, action)

      expect(getSessionStorage()).toEqual({authenticated: true})
    })
  })

  describe('RECEIVE_SIGNIN_FAIL', () => {
    it('stops fetching and clears session errors', () => {
      const previousState = {
        session: {
          authenticated: true,
          fetching: true,
          errors: ['error 1']
        }
      }
      const action = {type: RECEIVE_SIGNIN_FAIL}
      const state = reducer(previousState, action)

      expect(state.session.authenticated).toEqual(false)
      expect(state.session.fetching).toEqual(false)
      expect(state.session.errors).toEqual([])
    })

    it('saves authenticated false to local storage', () => {
      const action = {type: RECEIVE_SIGNIN_FAIL}
      reducer(undefined, action)

      expect(getSessionStorage()).toEqual({authenticated: false})
    })
  })

  it('RECEIVE_SIGNIN_ERROR stops fetching and assigns a network error', () => {
    const previousState = {
      session: {
        authenticated: true,
        fetching: true,
        networkError: false
      }
    }
    const action = {
      type: RECEIVE_SIGNIN_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.session.fetching).toEqual(false)
    expect(state.session.networkError).toEqual(true)
    expect(state.session.authenticated).toEqual(false)
  })

  it('REQUEST_SIGNOUT starts fetching and disables the network error', () => {
    const action = {type: REQUEST_SIGNOUT}
    const state = reducer(undefined, action)

    expect(state.session.fetching).toEqual(true)
    expect(state.session.networkError).toEqual(false)
  })

  describe('RECEIVE_SIGNOUT_SUCCESS', () => {
    it('stops fetching and assigns authenticated', () => {
      const previousState = {
        session: {
          authenticated: true,
          fetching: true,
          networkError: true
        }
      }
      const action = {type: RECEIVE_SIGNOUT_SUCCESS, authenticated: false}
      const state = reducer(previousState, action)

      expect(state.session.authenticated).toEqual(false)
      expect(state.session.fetching).toEqual(false)
      expect(state.session.networkError).toEqual(false)
    })

    it('saves authenticated false to local storage', () => {
      const action = {type: RECEIVE_SIGNOUT_SUCCESS, authenticated: false}
      reducer(undefined, action)

      expect(getSessionStorage()).toEqual({authenticated: false})
    })
  })

  it('RECEIVE_SIGNOUT_ERROR stops fetching and assigns a network error', () => {
    const previousState = {
      session: {
        authenticated: true,
        fetching: true,
        networkError: false
      }
    }
    const action = {
      type: RECEIVE_SIGNOUT_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.session.fetching).toEqual(false)
    expect(state.session.networkError).toEqual(true)
    expect(state.session.authenticated).toEqual(false)
  })
})
