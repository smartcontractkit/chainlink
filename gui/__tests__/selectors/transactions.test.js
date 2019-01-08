import transactionsSelector from 'selectors/transactions'

describe('selectors - transactions', () => {
  const CURRENT_PAGE = ['transactionA', 'transactionB']

  it('returns the transactions in the current page', () => {
    const state = {
      transactionsIndex: { currentPage: CURRENT_PAGE },
      transactions: {
        items: {
          transactionA: { id: 'transactionA' },
          transactionB: { id: 'transactionB' }
        }
      }
    }
    const transactions = transactionsSelector(state)

    expect(transactions).toEqual([
      { id: 'transactionA' },
      { id: 'transactionB' }
    ])
  })

  it('excludes transaction items that are not present', () => {
    const state = {
      transactionsIndex: { currentPage: CURRENT_PAGE },
      transactions: {
        items: {
          transactionA: { id: 'transactionA' }
        }
      }
    }
    const transactions = transactionsSelector(state)

    expect(transactions).toEqual([
      { id: 'transactionA' }
    ])
  })
})
