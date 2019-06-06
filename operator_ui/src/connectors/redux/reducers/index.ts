import { combineReducers } from 'redux'
import accountBalances from './accountBalances'
import authentication from './authentication'
import bridges, { State as BridgesState } from './bridges'
import configuration, { State as ConfigurationState } from './configuration'
import dashboardIndex, { State as DashboardState } from './dashboardIndex'
import fetching from './fetching'
import jobRuns, { State as JobRunsState } from './jobRuns'
import jobs from './jobs'
import notifications, { State as NotificationsState } from './notifications'
import redirect from './redirect'
import transactions from './transactions'
import transactionsIndex from './transactionsIndex'

export interface AppState {
  bridges: BridgesState
  configuration: ConfigurationState
  dashboardIndex: DashboardState
  jobRuns: JobRunsState
  notifications: NotificationsState
}

const reducer = combineReducers({
  accountBalances,
  authentication,
  bridges,
  configuration,
  dashboardIndex,
  fetching,
  jobRuns,
  jobs,
  notifications,
  redirect,
  transactions,
  transactionsIndex,
})

export default reducer
