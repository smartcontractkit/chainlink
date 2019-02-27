import reducer from 'connectors/redux/reducers'
import {
  REQUEST_BRIDGES,
  RECEIVE_BRIDGES_SUCCESS,
  RECEIVE_BRIDGES_ERROR,
  REQUEST_BRIDGE,
  RECEIVE_BRIDGE_SUCCESS,
  RECEIVE_BRIDGE_ERROR
} from 'actions'

describe('connectors/reducers/bridges', () => {
  it('returns the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.bridges).toEqual({
      items: {},
      currentPage: [],
      count: 0,
      networkError: false
    })
  })

  it('REQUEST_BRIDGES disables the network error', () => {
    const action = { type: REQUEST_BRIDGES }
    const previousState = {
      bridges: { networkError: true }
    }
    const state = reducer(previousState, action)

    expect(state.bridges.networkError).toEqual(false)
  })

  it('RECEIVE_BRIDGES_SUCCESS stores the bridge items and the current page', () => {
    const action = {
      type: RECEIVE_BRIDGES_SUCCESS,
      items: [{ id: 'a', name: 'A' }, { id: 'b', name: 'B' }]
    }
    const state = reducer(undefined, action)

    expect(state.bridges.items).toEqual({
      a: { id: 'a', name: 'A' },
      b: { id: 'b', name: 'B' }
    })
    expect(state.bridges.currentPage).toEqual(['a', 'b'])
    expect(state.bridges.networkError).toEqual(false)
  })

  it('RECEIVE_BRIDGES_ERROR updates the network error', () => {
    const previousState = {
      bridges: { networkError: false }
    }
    const action = {
      type: RECEIVE_BRIDGES_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.bridges.networkError).toEqual(true)
  })

  it('REQUEST_BRIDGE disables the network error', () => {
    const previousState = {
      bridges: { networkError: true }
    }
    const action = { type: REQUEST_BRIDGE }
    const state = reducer(previousState, action)

    expect(state.bridges.networkError).toEqual(false)
  })

  it('RECEIVE_BRIDGE_SUCCESS adds to the list of items', () => {
    const action = {
      type: RECEIVE_BRIDGE_SUCCESS,
      item: {
        id: 'a',
        name: 'A'
      }
    }
    const previousState = {
      bridges: { items: [] }
    }
    const state = reducer(previousState, action)

    expect(state.bridges.items.a).toEqual({
      id: 'a',
      name: 'A'
    })
  })

  it('RECEIVE_BRIDGE_ERROR assigns a network error', () => {
    const previousState = {
      bridges: { networkError: false }
    }
    const action = {
      type: RECEIVE_BRIDGE_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.bridges.networkError).toEqual(true)
  })
})
