import { combineReducers } from 'redux'
import accountBalances from './accountBalances'
import authentication from './authentication'
import bridges from './bridges'
import configuration from './configuration'
import create from './create'
import fetching from './fetching'
import jobRuns from './jobRuns'
import jobs from './jobs'
import transactions from './transactions'
import notifications from './notifications'
import redirect from './redirect'
import dashboardIndex from './dashboardIndex'
import transactionsIndex from './transactionsIndex'

const reducer = combineReducers({
  accountBalances,
  authentication,
  bridges,
  configuration,
  create,
  fetching,
  jobRuns,
  jobs,
  transactions,
  notifications,
  redirect,
  dashboardIndex,
  transactionsIndex
})

export default reducer
