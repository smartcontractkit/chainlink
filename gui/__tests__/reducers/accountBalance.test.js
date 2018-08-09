import reducer from 'connectors/redux/reducers'
import {
  REQUEST_ACCOUNT_BALANCE,
  RECEIVE_ACCOUNT_BALANCE_SUCCESS,
  RECEIVE_ACCOUNT_BALANCE_ERROR
} from 'actions'

describe('accountBalance reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.accountBalance).toEqual({
      eth: null,
      link: null,
      networkError: false
    })
  })

  it('REQUEST_ACCOUNT_BALANCE disables the network error', () => {
    const action = {type: REQUEST_ACCOUNT_BALANCE}
    const state = reducer(undefined, action)

    expect(state.accountBalance.networkError).toEqual(false)
  })

  it('RECEIVE_ACCOUNT_BALANCE_SUCCESS assigns the eth & link balance', () => {
    const previousState = {
      accountBalance: {
        networkError: true
      }
    }
    const action = {
      type: RECEIVE_ACCOUNT_BALANCE_SUCCESS,
      eth: '100',
      link: '200'
    }
    const state = reducer(previousState, action)

    expect(state.accountBalance.eth).toEqual('100')
    expect(state.accountBalance.link).toEqual('200')
    expect(state.accountBalance.networkError).toEqual(false)
  })

  it('RECEIVE_ACCOUNT_BALANCE_ERROR assigns a network error', () => {
    const previousState = {
      accountBalance: {
        networkError: false
      }
    }
    const action = {
      type: RECEIVE_ACCOUNT_BALANCE_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.accountBalance.networkError).toEqual(true)
  })
})
