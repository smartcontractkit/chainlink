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
  })
})
