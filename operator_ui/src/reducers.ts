import { combineReducers } from 'redux'
import accountBalances from './reducers/accountBalances'
import authentication from './reducers/authentication'
import bridges, { State as BridgesState } from './reducers/bridges'
import configuration, {
  State as ConfigurationState,
} from './reducers/configuration'
import dashboardIndex, {
  State as DashboardState,
} from './reducers/dashboardIndex'
import fetching from './reducers/fetching'
import jobRuns, { State as JobRunsState } from './reducers/jobRuns'
import jobs from './reducers/jobs'
import notifications, {
  State as NotificationsState,
} from './reducers/notifications'
import redirect from './reducers/redirect'
import transactions from './reducers/transactions'
import transactionsIndex from './reducers/transactionsIndex'

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
