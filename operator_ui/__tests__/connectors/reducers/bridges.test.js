import reducer from 'connectors/redux/reducers'

describe('connectors/reducers/bridges', () => {
  it('returns the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.bridges).toEqual({
      items: {},
      currentPage: [],
      count: 0
    })
  })

  it('UPSERT_BRIDGES stores the bridge items, current page & count', () => {
    const action = {
      type: 'UPSERT_BRIDGES',
      data: {
        bridges: {
          a: { id: 'a', name: 'A' },
          b: { id: 'b', name: 'B' }
        },
        meta: {
          currentPageBridges: {
            data: [{ id: 'a' }, { id: 'b' }],
            meta: { count: 5 }
          }
        }
      }
    }
    const state = reducer(undefined, action)

    expect(state.bridges.items).toEqual({
      a: { id: 'a', name: 'A' },
      b: { id: 'b', name: 'B' }
    })
    expect(state.bridges.currentPage).toEqual(['a', 'b'])
    expect(state.bridges.count).toEqual(5)
  })

  it('UPSERT_BRIDGE stores the bridge item', () => {
    const action = {
      type: 'UPSERT_BRIDGE',
      data: {
        bridges: {
          a: {
            id: 'a',
            attributes: {
              name: 'A'
            }
          }
        }
      }
    }
    const previousState = {
      bridges: { items: {} }
    }
    const state = reducer(previousState, action)

    expect(state.bridges.items.a).toEqual({
      id: 'a',
      attributes: {
        name: 'A'
      }
    })
  })
})
