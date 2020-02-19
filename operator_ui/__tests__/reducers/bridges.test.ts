import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  UpsertBridgeAction,
  UpsertBridgesAction,
  ResourceActionType,
} from '../../src/reducers/actions'

describe('reducers/bridges', () => {
  it('UPSERT_BRIDGES stores the bridge items, current page & count', () => {
    const action: UpsertBridgesAction = {
      type: ResourceActionType.UPSERT_BRIDGES,
      data: {
        bridges: {
          a: { id: 'a', name: 'A' },
          b: { id: 'b', name: 'B' },
        },
        meta: {
          currentPageBridges: {
            data: [{ id: 'a' }, { id: 'b' }],
            meta: { count: 5 },
          },
        },
      },
    }
    const state = reducer(undefined, action)

    expect(state.bridges.items).toEqual({
      a: { id: 'a', name: 'A' },
      b: { id: 'b', name: 'B' },
    })
    expect(state.bridges.currentPage).toEqual(['a', 'b'])
    expect(state.bridges.count).toEqual(5)
  })

  it('UPSERT_BRIDGE stores the bridge item', () => {
    const action: UpsertBridgeAction = {
      type: ResourceActionType.UPSERT_BRIDGE,
      data: {
        bridges: {
          a: {
            id: 'a',
            attributes: {
              name: 'A',
            },
          },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.bridges.items.a).toEqual({
      id: 'a',
      attributes: {
        name: 'A',
      },
    })
  })
})
