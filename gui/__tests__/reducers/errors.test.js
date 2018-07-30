import reducer from 'connectors/redux/reducers'
import {
  RECEIVE_SESSION_FAIL
} from 'actions'

describe('errors reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.errors).toEqual([])
  })

  it('RECEIVE_SESSION_FAIL adds a failure message', () => {
    const action = {type: RECEIVE_SESSION_FAIL}
    const state = reducer(undefined, action)

    expect(state.errors).toEqual([
      'Your email or password are incorrect. Please try again'
    ])
  })
})
