import bridgeSelector from 'selectors/bridge'

describe('selectors - bridge', () => {
  it('returns the bridge with the given id', () => {
    const state = {
      bridges: {
        items: {
          a: { name: 'A' },
          b: { name: 'B' }
        }
      }
    }

    expect(bridgeSelector(state, 'a')).toEqual({ name: 'A' })
  })
})
