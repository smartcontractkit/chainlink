import { InitialStateAction } from 'state/actions'
import reducer, { INITIAL_STATE } from './reducers'

describe('state/ducks/feeds/reducers', () => {
  describe('initial state', () => {
    it('stores a map from feeds.json by contract address', () => {
      const action: InitialStateAction = { type: 'INITIAL_STATE' }
      const state = reducer(INITIAL_STATE, action)
      const keys = Object.keys(state.items)

      expect(keys.length).toBeGreaterThan(0)

      const firstKey = keys[0]
      const firstFeed = state.items[firstKey]
      expect(firstFeed.contractAddress).toEqual(firstKey)
    })

    it('maintains an index with the order', () => {
      const action: InitialStateAction = { type: 'INITIAL_STATE' }
      const state = reducer(INITIAL_STATE, action)

      expect(state.order.length).toBeGreaterThan(0)
      expect(state.order[0]).toEqual(
        '0xF79D6aFBb6dA890132F9D7c355e3015f15F3406F',
      )
    })

    it('maintains an index by pair and network', () => {
      const action: InitialStateAction = { type: 'INITIAL_STATE' }
      const state = reducer(INITIAL_STATE, action)

      expect(state.pairPaths.length).toBeGreaterThan(0)
      expect(state.pairPaths[0]).toEqual(['eth-usd', 1, '0xF79D6aFBb6dA890132F9D7c355e3015f15F3406F'])
    })
  })
})
