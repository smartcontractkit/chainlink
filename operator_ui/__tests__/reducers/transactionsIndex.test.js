import reducer from 'reducers'
import { UPSERT_TRANSACTIONS } from 'reducers/transactionsIndex'

describe('reducers/transactionsIndex', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.transactionsIndex).toEqual({
      currentPage: null,
      count: 0,
    })
  })

  it('UPSERT_TRANSACTIONS updates the current page & count from meta', () => {
    const action = {
      type: UPSERT_TRANSACTIONS,
      data: {
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
    const state = reducer(undefined, action)

    expect(state.transactionsIndex.count).toEqual(10)
    expect(state.transactionsIndex.currentPage).toEqual(['b', 'a'])
  })
})
