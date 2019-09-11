const initialState = { items: {} }

export const UPSERT_TRANSACTIONS = 'UPSERT_TRANSACTIONS'
export const UPSERT_TRANSACTION = 'UPSERT_TRANSACTION'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_TRANSACTIONS: {
      return Object.assign({}, state, {
        items: Object.assign({}, state.items, action.data.transactions),
      })
    }
    case UPSERT_TRANSACTION: {
      return Object.assign({}, state, {
        items: Object.assign({}, state.items, action.data.transactions),
      })
    }
    default:
      return state
  }
}
