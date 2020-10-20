import { AppState } from 'reducers'
import build from 'redux-object'

export default ({
  transactionsIndex,
  transactions,
}: Pick<AppState, 'transactionsIndex' | 'transactions'>) => {
  return (
    transactionsIndex.currentPage &&
    transactionsIndex.currentPage
      .map((id) => build(transactions, 'items', id))
      .filter((t) => t)
  )
}
