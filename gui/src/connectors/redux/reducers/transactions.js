const initialState = { items: {} }

export const UPSERT_TRANSACTIONS = 'UPSERT_TRANSACTIONS'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_TRANSACTIONS: {
      return Object.assign(
        {},
        state,
        { items: Object.assign({}, state.items, action.data.txattempts) }
      )
    }
    default:
      return state
  }
}
