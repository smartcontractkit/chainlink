import { combineReducers } from 'redux'
import accountBalances from './accountBalances'
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
  accountBalances,
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
