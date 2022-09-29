import { partialAsFull } from 'support/test-helpers/partialAsFull'
import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  RedirectAction,
  RouterActionType,
  ResourceActionType,
  RequestCreateAction,
  ReceiveCreateSuccessAction,
} from '../../src/reducers/actions'

describe('connectors/reducers/fetching', () => {
  const incrementAction: RequestCreateAction = {
    type: ResourceActionType.REQUEST_CREATE,
  }
  const receiveDecrementAction: ReceiveCreateSuccessAction = {
    type: ResourceActionType.RECEIVE_CREATE_SUCCESS,
  }

  it('increments count when type starts with REQUEST_ & decrements with RECEIVE_ or RESPONSE_', () => {
    let state = reducer(INITIAL_STATE, incrementAction)

    state = reducer(state, incrementAction)
    expect(state.fetching.count).toEqual(2)

    state = reducer(state, receiveDecrementAction)
    expect(state.fetching.count).toEqual(1)
  })

  it('does not negatively decrement count on RECEIVE_ or RESPONSE_', () => {
    const state = reducer(INITIAL_STATE, receiveDecrementAction)
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
