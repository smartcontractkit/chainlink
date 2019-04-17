import bridgesSelector from 'selectors/bridges'

describe('selectors - bridges', () => {
  it('returns the current page of bridges', () => {
    const state = {
      bridges: {
        items: {
          a: { name: 'A' },
          b: { name: 'B' },
          c: { name: 'C' }
        },
        currentPage: ['c', 'a']
      }
    }

    expect(bridgesSelector(state, 'a')).toEqual([{ name: 'C' }, { name: 'A' }])
  })

  it('does not return items that cannot be found', () => {
    const state = {
      bridges: {
        items: {
          a: { name: 'A' },
          b: { name: 'B' },
          c: { name: 'C' }
        },
        currentPage: ['C', 'A']
      }
    }

    expect(bridgesSelector(state, 'a')).toEqual([])
  })
})
