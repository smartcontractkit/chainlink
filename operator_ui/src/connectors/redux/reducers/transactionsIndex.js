const initialState = {
  currentPage: null,
  count: 0
}

export const UPSERT_TRANSACTIONS = 'UPSERT_TRANSACTIONS'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_TRANSACTIONS: {
      const { data } = action
      return Object.assign(
        {},
        state,
        { currentPage: data.meta.currentPageTransactions.data.map(t => t.id) },
        { count: data.meta.currentPageTransactions.meta.count }
      )
    }
    default:
      return state
  }
}
