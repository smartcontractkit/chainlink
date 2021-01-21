import { partialAsFull } from 'support/test-helpers/partialAsFull'
import { INITIAL_STATE, AppState } from '../../src/reducers'
import transactionsSelector from '../../src/selectors/transactions'

describe('selectors - transactions', () => {
  type State = Pick<AppState, 'transactions' | 'transactionsIndex'>
  type TransactionsIndexState = typeof INITIAL_STATE.transactionsIndex
  const transactionsIndexState = partialAsFull<TransactionsIndexState>({
    currentPage: ['transactionA', 'transactionB'],
  })

  it('returns the transactions in the current page', () => {
    const transactionsState = {
      items: {
        transactionA: { id: 'transactionA' },
        transactionB: { id: 'transactionB' },
      },
    }
    const state: State = {
      transactionsIndex: transactionsIndexState,
      transactions: transactionsState,
    }
    const transactions = transactionsSelector(state)

    expect(transactions).toEqual([
      { id: 'transactionA' },
      { id: 'transactionB' },
    ])
  })

  it('excludes transaction items that are not present', () => {
    const transactionsState = {
      items: {
        transactionA: { id: 'transactionA' },
      },
    }
    const state: State = {
      transactionsIndex: transactionsIndexState,
      transactions: transactionsState,
    }
    const transactions = transactionsSelector(state)

    expect(transactions).toEqual([{ id: 'transactionA' }])
  })
})
