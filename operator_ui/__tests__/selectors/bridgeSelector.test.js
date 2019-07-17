import bridgeSelector from 'selectors/bridge'

describe('selectors - bridge', () => {
  it('returns the bridge with the given id', () => {
    const state = {
      bridges: {
        items: {
          a: { attributes: { name: 'A' } },
          b: { attributes: { name: 'B' } }
        }
      }
    }

    const selected = bridgeSelector(state, 'a')
    expect(selected).toEqual({ id: 'a', name: 'A' })
  })
})
