import { AppState } from '../../src/reducers'
import fetchCountSelector from '../../src/selectors/fetchCount'

describe('selectors - fetchCount', () => {
  it('returns the value of the counter', () => {
    const state: Pick<AppState, 'fetching'> = { fetching: { count: 1 } }

    expect(fetchCountSelector(state)).toEqual(1)
  })
})
