import build from 'redux-object'

export default ({ transactionsIndex, transactions }) => (
  transactionsIndex.currentPage && transactionsIndex
    .currentPage
    .map(id => build(transactions, 'items', id))
    .filter(t => t)
)
