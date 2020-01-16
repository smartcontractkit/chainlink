import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  UpsertTransactionAction,
  UpsertTransactionsAction,
  ResourceActionType,
} from '../../src/reducers/actions'

describe('reducers/transactions', () => {
  it('UPSERT_TRANSACTIONS upserts items', () => {
    const action: UpsertTransactionsAction = {
      type: ResourceActionType.UPSERT_TRANSACTIONS,
      data: {
        transactions: {
          a: { id: 'a' },
          b: { id: 'b' },
        },
        meta: {
          currentPageTransactions: {
            data: [],
            meta: { count: 2 },
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
    const action: UpsertTransactionAction = {
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
