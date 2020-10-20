import { AppState } from '../../src/reducers'
import bridgeSelector from '../../src/selectors/bridge'

describe('selectors - bridge', () => {
  it('returns the bridge with the given id', () => {
    const state: Pick<AppState, 'bridges'> = {
      bridges: {
        items: {
          a: { attributes: { name: 'A' } },
          b: { attributes: { name: 'B' } },
        },
        count: 0,
        currentPage: ['0'],
      },
    }

    const selected = bridgeSelector(state, 'a')
    expect(selected).toEqual({ id: 'a', name: 'A' })
  })
})
