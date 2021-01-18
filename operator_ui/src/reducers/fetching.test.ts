import { partialAsFull } from 'support/test-helpers/partialAsFull'
import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  RedirectAction,
  RouterActionType,
  ResourceActionType,
  RequestCreateAction,
  ReceiveCreateSuccessAction,
  ResponseAccountBalanceAction,
} from '../../src/reducers/actions'

describe('connectors/reducers/fetching', () => {
  const incrementAction: RequestCreateAction = {
    type: ResourceActionType.REQUEST_CREATE,
  }
  const receiveDecrementAction: ReceiveCreateSuccessAction = {
    type: ResourceActionType.RECEIVE_CREATE_SUCCESS,
  }
  const responseDecrementAction: ResponseAccountBalanceAction = {
    type: ResourceActionType.RESPONSE_ACCOUNT_BALANCE,
  }

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
