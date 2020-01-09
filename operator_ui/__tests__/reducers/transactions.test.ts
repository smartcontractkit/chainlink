import reducer, { INITIAL_STATE } from '../../src/reducers'
import { ResourceActionType } from '../../src/reducers/actions'

describe('reducers/transactions', () => {
  it('UPSERT_TRANSACTIONS upserts items', () => {
    const action = {
      type: ResourceActionType.UPSERT_TRANSACTIONS,
      data: {
        transactions: {
          a: { id: 'a' },
          b: { id: 'b' },
        },
        meta: {
          currentPageTransactions: {
            data: [],
            meta: {},
          },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.transactions.items).toEqual({
      a: { id: 'a' },
      b: { id: 'b' },
    })
  })

  it('UPSERT_TRANSACTION upserts items', () => {
    const action = {
      type: ResourceActionType.UPSERT_TRANSACTION,
      data: {
        transactions: {
          a: { id: 'a' },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.transactions.items).toEqual({ a: { id: 'a' } })
  })
})
