import reducer from 'connectors/redux/reducers'
import {
  REQUEST_BRIDGESPEC,
  RECEIVE_BRIDGESPEC_SUCCESS,
  RECEIVE_BRIDGESPEC_ERROR
} from 'actions'

describe('bridgeSpec reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.bridgeSpec).toEqual({
      name: '',
      url: '',
      confirmations: '0',
      networkError: false,
      fetching: false
    })
  })

  it('REQUEST_BRIDGESPEC starts fetching and disables the network error', () => {
    const action = {type: REQUEST_BRIDGESPEC}
    const state = reducer(undefined, action)

    expect(state.bridgeSpec.fetching).toEqual(true)
    expect(state.bridgeSpec.networkError).toEqual(false)
  })

  it('RECEIVE_BRIDGESPEC_SUCCESS stops fetching and assigns properties', () => {
    const previousState = {
      bridgeSpec: {
        fetching: true,
        networkError: true
      }
    }
    const action = {
      type: RECEIVE_BRIDGESPEC_SUCCESS,
      name: 'someRandomName',
      url: 'https://localhost.com:8000/endpoint',
      confirmations: 5,
      incomingToken: 'abc',
      outgoingToken: '123'
    }
    const state = reducer(previousState, action)

    expect(state.bridgeSpec.name).toEqual('someRandomName')
    expect(state.bridgeSpec.url).toEqual('https://localhost.com:8000/endpoint')
    expect(state.bridgeSpec.confirmations).toEqual(5)
    expect(state.bridgeSpec.incomingToken).toEqual('abc')
    expect(state.bridgeSpec.outgoingToken).toEqual('123')
    expect(state.bridgeSpec.fetching).toEqual(false)
    expect(state.bridgeSpec.networkError).toEqual(false)
  })

  it('RECEIVE_BRIDGESPEC_ERROR stops fetching and assigns a network error', () => {
    const previousState = {
      bridgeSpec: {
        fetching: true,
        networkError: false
      }
    }
    const action = {
      type: RECEIVE_BRIDGESPEC_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.bridgeSpec.fetching).toEqual(false)
    expect(state.bridgeSpec.networkError).toEqual(true)
  })
})
