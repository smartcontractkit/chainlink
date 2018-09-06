import reducer from 'connectors/redux/reducers'
import {
  MATCH_ROUTE,
  RECEIVE_SIGNIN_FAIL,
  RECEIVE_CREATE_SUCCESS,
  RECEIVE_CREATE_ERROR
} from 'actions'

describe('notifications reducer', () => {
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
        errors: [{detail: 'error 1'}],
        successes: [{id: '123'}],
        currentUrl: null
      }
    }

    const sameUrlAction = {type: MATCH_ROUTE, match: {url: null}}
    let state = reducer(previousState, sameUrlAction)

    expect(state.notifications).toEqual({
      errors: [{detail: 'error 1'}],
      successes: [{id: '123'}],
      currentUrl: null
    })

    const changedUrlAction = {type: MATCH_ROUTE, match: {url: '/'}}
    state = reducer(previousState, changedUrlAction)
    expect(state.notifications).toEqual({
      errors: [],
      successes: [],
      currentUrl: '/'
    })
  })

  it('RECEIVE_SIGNIN_FAIL adds a failure', () => {
    const action = {type: RECEIVE_SIGNIN_FAIL}
    const state = reducer(undefined, action)

    expect(state.notifications).toEqual({
      errors: [{detail: 'Your email or password is incorrect. Please try again'}],
      successes: [],
      currentUrl: null
    })
  })

  it('RECEIVE_CREATE_ERROR adds a failure and clears successes', () => {
    const previousState = {
      notifications: {
        errors: [],
        successes: [{id: '123'}],
        currentUrl: null
      }
    }

    const action = {type: RECEIVE_CREATE_ERROR, error: {errors: [{detail: 'Invalid name'}]}}
    const state = reducer(previousState, action)

    expect(state.notifications).toEqual({
      errors: [{detail: 'Invalid name'}],
      successes: [],
      currentUrl: null
    })
  })

  it('RECEIVE_CREATE_SUCCESS adds a success and clears errors', () => {
    const previousState = {
      notifications: {
        errors: [{detail: 'error 1'}],
        successes: [],
        currentUrl: null
      }
    }

    const response = {id: 'SOMEID', name: 'SOMENAME'}
    const action = {type: RECEIVE_CREATE_SUCCESS, response: response}
    const state = reducer(previousState, action)

    expect(state.notifications).toEqual({
      errors: [],
      successes: [response],
      currentUrl: null
    })
  })
})
