import { AnyAction } from 'redux'

const initialState = {}

export const UPSERT_ACCOUNT_BALANCE = 'UPSERT_ACCOUNT_BALANCE'

type Action = { type: 'UPSERT_ACCOUNT_BALANCE' }

export default (state = initialState, action: AnyAction = {}) => {
  switch (action.type) {
    case UPSERT_ACCOUNT_BALANCE:
      return Object.assign({}, state, action.data.accountBalances)
    default:
      return state
  }
}
