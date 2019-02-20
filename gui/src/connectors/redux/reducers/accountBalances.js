const initialState = {}

export const UPSERT_ACCOUNT_BALANCE = 'UPSERT_ACCOUNT_BALANCE'

export default (state = initialState, action = {}) => {
  switch (action.type) {
    case UPSERT_ACCOUNT_BALANCE:
      return Object.assign({}, state, action.data.accountBalances)
    default:
      return state
  }
}
