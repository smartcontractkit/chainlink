import reducer from 'connectors/redux/reducers'
import {
  REQUEST_SESSION,
  REQUEST_SIGNOUT,
  RECEIVE_SESSION_SUCCESS,
  RECEIVE_SESSION_ERROR,
  RECEIVE_SIGNOUT_SUCCESS,
  RECEIVE_SIGNOUT_ERROR
} from 'actions'

describe('session reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.session).toEqual({
      fetching: false,
      authenticated: false,
      errors: [],
      networkError: false
    })
  })

  it('REQUEST_SESSION starts fetching and disables the network error', () => {
    const action = {type: REQUEST_SESSION}
    const state = reducer(undefined, action)

    expect(state.session.fetching).toEqual(true)
    expect(state.session.networkError).toEqual(false)
  })

  it('RECEIVE_SESSION_SUCCESS stops fetching and assigns authenticated', () => {
    const previousState = {
      session: {
        fetching: true,
        networkError: true
      }
    }
    const action = {
      type: RECEIVE_SESSION_SUCCESS,
      authenticated: true
    }
    const state = reducer(previousState, action)

    expect(state.session.authenticated).toEqual(true)
    expect(state.session.fetching).toEqual(false)
    expect(state.session.networkError).toEqual(false)
  })

  it('RECEIVE_SESSION_ERROR stops fetching and assigns a network error', () => {
    const previousState = {
      session: {
        authenticated: true,
        fetching: true,
        networkError: false
      }
    }
    const action = {
      type: RECEIVE_SESSION_ERROR,
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

  it('RECEIVE_SIGNOUT_SUCCESS stops fetching and assigns authenticated', () => {
    const previousState = {
      session: {
        authenticated: true,
        fetching: true,
        networkError: true
      }
    }
    const action = {
      type: RECEIVE_SIGNOUT_SUCCESS,
      authenticated: false
    }
    const state = reducer(previousState, action)

    expect(state.session.authenticated).toEqual(false)
    expect(state.session.fetching).toEqual(false)
    expect(state.session.networkError).toEqual(false)
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
