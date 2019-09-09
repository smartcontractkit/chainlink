import fetchCountSelector from 'selectors/fetchCount'

describe('selectors - fetchCount', () => {
  it('returns the value of the counter', () => {
    const state = { fetching: { count: 1 } }

    expect(fetchCountSelector(state)).toEqual(1)
  })
})
