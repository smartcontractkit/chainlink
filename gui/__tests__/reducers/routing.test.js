import reducer from 'connectors/redux/reducers'
import {
  MATCH_ROUTE,
  RECEIVE_SIGNIN_FAIL
} from 'actions'

describe('routing reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.routing).toEqual({
      errors: [],
      currentUrl: null
    })
  })

  it('RECEIVE_SIGNIN_FAIL adds a failure', () => {
    const action = {type: RECEIVE_SIGNIN_FAIL}
    const state = reducer(undefined, action)

    expect(state.routing).toEqual({
      errors: [{detail: 'Your email or password is incorrect. Please try again'}],
      currentUrl: null
    })
  })

  it('MATCH_ROUTE clears errors when currentUrl changes', () => {
    const previousState = {
      routing: {
        errors: [{detail: 'error 1'}],
        currentUrl: null
      }
    }

    const sameUrlAction = {type: MATCH_ROUTE, match: {url: null}}
    let state = reducer(previousState, sameUrlAction)

    expect(state.routing).toEqual({
      errors: [{detail: 'error 1'}],
      currentUrl: null
    })

    const changedUrlAction = {type: MATCH_ROUTE, match: {url: '/'}}
    state = reducer(previousState, changedUrlAction)
    expect(state.routing).toEqual({
      errors: [],
      currentUrl: '/'
    })
  })
})
