import { combineReducers } from 'redux'
import notifications from './notifications'
import fetching from './fetching'
import accountBalance from './accountBalance'
import bridges from './bridges'
import bridgeSpec from './bridgeSpec'
import jobs from './jobs'
import jobRuns from './jobRuns'
import create from './create'
import configuration from './configuration'
import authentication from './authentication'

const reducer = combineReducers({
  notifications,
  fetching,
  accountBalance,
  bridges,
  bridgeSpec,
  jobs,
  jobRuns,
  create,
  configuration,
  authentication
})

export default reducer
