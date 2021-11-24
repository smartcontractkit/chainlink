import { useSelector, TypedUseSelectorHook } from 'react-redux'
import { combineReducers } from 'redux'
import accountBalances from './reducers/accountBalances'
import authentication from './reducers/authentication'
import configuration from './reducers/configuration'
import fetching from './reducers/fetching'
import notifications from './reducers/notifications'
import redirect from './reducers/redirect'
import transactions from './reducers/transactions'
import transactionsIndex from './reducers/transactionsIndex'

const reducer = combineReducers({
  accountBalances,
  authentication,
  configuration,
  fetching,
  notifications,
  redirect,
  transactions,
  transactionsIndex,
})

export const INITIAL_STATE = reducer(undefined, { type: 'INITIAL_STATE' })
export type AppState = typeof INITIAL_STATE
export const useOperatorUiSelector: TypedUseSelectorHook<AppState> = useSelector

export default reducer
