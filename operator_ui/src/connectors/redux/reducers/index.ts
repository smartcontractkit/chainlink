import { combineReducers } from 'redux'
import accountBalances from './accountBalances'
import authentication from './authentication'
import bridges, { IState as IBridgesState } from './bridges'
import configuration, { IState as IConfigurationState } from './configuration'
import fetching from './fetching'
import jobRuns, { IState as IJobRunsState } from './jobRuns'
import jobs from './jobs'
import transactions from './transactions'
import notifications from './notifications'
import redirect from './redirect'
import dashboardIndex, { IState as IDashboardState } from './dashboardIndex'
import transactionsIndex from './transactionsIndex'

export interface IState {
  bridges: IBridgesState
  configuration: IConfigurationState
  dashboardIndex: IDashboardState
  jobRuns: IJobRunsState
}

const reducer = combineReducers({
  accountBalances,
  authentication,
  bridges,
  configuration,
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
