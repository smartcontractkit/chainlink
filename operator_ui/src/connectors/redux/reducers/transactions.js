const initialState = { items: {} }

export const UPSERT_TRANSACTIONS = 'UPSERT_TRANSACTIONS'
export const UPSERT_TRANSACTION = 'UPSERT_TRANSACTION'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_TRANSACTIONS: {
      return {
        ...state,

        items: {
          ...state.items,
          ...action.data.transactions
        }
      }
    }
    case UPSERT_TRANSACTION: {
      return {
        ...state,

        items: {
          ...state.items,
          ...action.data.transactions
        }
      }
    }
    default:
      return state
  }
}
