import { combineReducers } from 'redux'
import accountBalance from './accountBalance'
import authentication from './authentication'
import bridges from './bridges'
import configuration from './configuration'
import create from './create'
import fetching from './fetching'
import jobRuns from './jobRuns'
import jobs from './jobs'
import notifications from './notifications'
import redirect from './redirect'

const reducer = combineReducers({
  accountBalance,
  authentication,
  bridges,
  configuration,
  create,
  fetching,
  jobRuns,
  jobs,
  notifications,
  redirect
})

export default reducer
