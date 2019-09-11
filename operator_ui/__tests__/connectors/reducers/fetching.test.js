import { RouterActionType } from 'actions'
import reducer from 'connectors/redux/reducers'

describe('connectors/reducers/fetching', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.fetching).toEqual({ count: 0 })
  })

  it('increments count when the action type starts with "REQUEST_"', () => {
    const action = { type: 'REQUEST_FOO' }
    const state = reducer(undefined, action)

    expect(state.fetching).toEqual({ count: 1 })
  })

  it('decrements count when the action type starts with "RECEIVE_"', () => {
    const action = { type: 'RECEIVE_FOO' }
    const previousState = {
      fetching: { count: 1 },
    }
    const state = reducer(previousState, action)

    expect(state.fetching).toEqual({ count: 0 })
  })

  it('does not negatively decrement count on "RECEIVE_"', () => {
    const action = { type: 'RECEIVE_FOO' }
    const state = reducer(undefined, action)

    expect(state.fetching).toEqual({ count: 0 })
  })

  it('decrements count when the action type starts with "RESPONSE_"', () => {
    const action = { type: 'RESPONSE_FOO' }
    const previousState = {
      fetching: { count: 1 },
    }
    const state = reducer(previousState, action)

    expect(state.fetching).toEqual({ count: 0 })
  })

  it('does not negatively decrement count on "RESPONSE_"', () => {
    const action = { type: 'RESPONSE_FOO' }
    const state = reducer(undefined, action)

    expect(state.fetching).toEqual({ count: 0 })
  })

  it('resets the counter on redirect', () => {
    const action = { type: RouterActionType.REDIRECT }
    const previousState = {
      fetching: { count: 1 },
    }
    const state = reducer(previousState, action)

    expect(state.fetching).toEqual({ count: 0 })
  })
})
