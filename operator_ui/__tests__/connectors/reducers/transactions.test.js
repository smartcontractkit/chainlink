import reducer from 'connectors/redux/reducers'
import {
  UPSERT_TRANSACTIONS,
  UPSERT_TRANSACTION,
} from 'connectors/redux/reducers/transactions'

describe('connectors/reducers/transactions', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.transactions).toEqual({ items: {} })
  })

  it('UPSERT_TRANSACTIONS upserts items', () => {
    const action = {
      type: UPSERT_TRANSACTIONS,
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
    const state = reducer(undefined, action)

    expect(state.transactions.items).toEqual({
      a: { id: 'a' },
      b: { id: 'b' },
    })
  })

  it('UPSERT_TRANSACTION upserts items', () => {
    const action = {
      type: UPSERT_TRANSACTION,
      data: {
        transactions: {
          a: { id: 'a' },
        },
      },
    }
    const state = reducer(undefined, action)

    expect(state.transactions.items).toEqual({ a: { id: 'a' } })
  })
})
