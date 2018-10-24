import isFetchingSelector from 'selectors/isFetching'

describe('selectors - isFetching', () => {
  it('is true when count > 0', () => {
    const state = {fetching: {count: 1}}

    expect(isFetchingSelector(state)).toEqual(true)
  })

  it('is false when count = 0', () => {
    const state = {fetching: {count: 0}}

    expect(isFetchingSelector(state)).toEqual(false)
  })
})
