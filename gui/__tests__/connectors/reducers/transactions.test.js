import reducer from 'connectors/redux/reducers'
import {
  UPSERT_TRANSACTIONS
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
        txattempts: {
          a: { id: 'a' },
          b: { id: 'b' }
        },
        meta: {
          currentPageTransactions: {
            data: [],
            meta: {}
          }
        }
      }
    }
    const state = reducer(undefined, action)

    expect(state.transactions.items).toEqual({
      a: { id: 'a' },
      b: { id: 'b' }
    })
  })
})
