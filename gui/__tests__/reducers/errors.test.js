import reducer from 'connectors/redux/reducers'
import {
  MATCH_ROUTE,
  RECEIVE_SESSION_FAIL
} from 'actions'

describe('errors reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.errors).toEqual({
      messages: [],
      currentUrl: null
    })
  })

  it('RECEIVE_SESSION_FAIL adds a failure message', () => {
    const action = {type: RECEIVE_SESSION_FAIL}
    const state = reducer(undefined, action)

    expect(state.errors).toEqual({
      messages: ['Your email or password are incorrect. Please try again'],
      currentUrl: null
    })
  })

  it('MATCH_ROUTE clears messages when currentUrl changes', () => {
    const previousState = {
      errors: {
        messages: ['error 1'],
        currentUrl: null
      }
    }

    const sameUrlAction = {type: MATCH_ROUTE, match: {url: null}}
    let state = reducer(previousState, sameUrlAction)

    expect(state.errors).toEqual({
      messages: ['error 1'],
      currentUrl: null
    })

    const changedUrlAction = {type: MATCH_ROUTE, match: {url: '/'}}
    state = reducer(previousState, changedUrlAction)
    expect(state.errors).toEqual({
      messages: [],
      currentUrl: '/'
    })
  })
})
