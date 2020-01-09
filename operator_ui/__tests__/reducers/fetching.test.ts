import { partialAsFull } from '@chainlink/ts-test-helpers'
import reducer, { INITIAL_STATE } from '../../src/reducers'
import { RedirectAction, RouterActionType } from '../../src/reducers/actions'

describe('connectors/reducers/fetching', () => {
  const incrementAction = { type: 'REQUEST_FOO' }
  const receiveDecrementAction = { type: 'RECEIVE_FOO' }
  const responseDecrementAction = { type: 'RESPONSE_FOO' }

  it('increments count when type starts with REQUEST_ & decrements with RECEIVE_ or RESPONSE_', () => {
    let state = reducer(INITIAL_STATE, incrementAction)

    state = reducer(state, incrementAction)
    expect(state.fetching.count).toEqual(2)

    state = reducer(state, receiveDecrementAction)
    expect(state.fetching.count).toEqual(1)

    state = reducer(state, responseDecrementAction)
    expect(state.fetching.count).toEqual(0)
  })

  it('does not negatively decrement count on RECEIVE_ or RESPONSE_', () => {
    let state = reducer(INITIAL_STATE, receiveDecrementAction)
    expect(state.fetching.count).toEqual(0)

    state = reducer(state, responseDecrementAction)
    expect(state.fetching.count).toEqual(0)
  })

  it('resets the counter on redirect', () => {
    const redirectAction = partialAsFull<RedirectAction>({
      type: RouterActionType.REDIRECT,
    })
    let state = reducer(INITIAL_STATE, incrementAction)

    state = reducer(state, incrementAction)
    expect(state.fetching.count).toEqual(2)

    state = reducer(state, redirectAction)
    expect(state.fetching.count).toEqual(0)
  })
})
