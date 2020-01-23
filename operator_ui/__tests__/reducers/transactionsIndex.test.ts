import reducer, { INITIAL_STATE } from '../../src/reducers'
import {
  ResourceActionType,
  UpsertTransactionsAction,
} from '../../src/reducers/actions'

describe('reducers/transactionsIndex', () => {
  it('UPSERT_TRANSACTIONS updates the current page & count from meta', () => {
    const action: UpsertTransactionsAction = {
      type: ResourceActionType.UPSERT_TRANSACTIONS,
      data: {
        transactions: {},
        meta: {
          currentPageTransactions: {
            data: [{ id: 'b' }, { id: 'a' }],
            meta: {
              count: 10,
            },
          },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.transactionsIndex.count).toEqual(10)
    expect(state.transactionsIndex.currentPage).toEqual(['b', 'a'])
  })
})
