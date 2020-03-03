import transactionSelector from '../../src/selectors/transaction'

describe('selectors - transaction', () => {
  it('returns the transaction item for the given id and null otherwise', () => {
    const state = {
      transactions: {
        items: {
          transactionA: { id: 'transactionA' },
        },
      },
    }

    expect(transactionSelector(state, 'transactionA')).toEqual({
      id: 'transactionA',
    })
    expect(transactionSelector(state, 'transactiona')).toBeNull()
  })
})
