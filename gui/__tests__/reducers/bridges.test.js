import reducer from 'connectors/redux/reducers'
import {
  REQUEST_BRIDGES,
  RECEIVE_BRIDGES_SUCCESS,
  RECEIVE_BRIDGES_ERROR
} from 'actions'

describe('bridges reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.bridges).toEqual({
      items: [],
      currentPage: [],
      count: 0,
      networkError: false
    })
  })

  it('REQUEST_BRIDGES disables the network error', () => {
    const action = {type: REQUEST_BRIDGES}
    const state = reducer(undefined, action)

    expect(state.bridges.networkError).toEqual(false)
  })

  it('RECEIVE_BRIDGES_SUCCESS stores the bridge items and the current page', () => {
    const action = {
      type: RECEIVE_BRIDGES_SUCCESS,
      items: [{name: 'a'}, {name: 'b'}]
    }
    const state = reducer(undefined, action)

    expect(state.bridges.items).toEqual([{name: 'a'}, {name: 'b'}])
    expect(state.bridges.currentPage).toEqual(['a', 'b'])
    expect(state.bridges.networkError).toEqual(false)
  })

  it('RECEIVE_BRIDGES_ERROR updates the network error', () => {
    const previousState = {
      bridges: {networkError: false}
    }
    const action = {
      type: RECEIVE_BRIDGES_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.bridges.networkError).toEqual(true)
  })
})
