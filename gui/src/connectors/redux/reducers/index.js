import { combineReducers } from 'redux'
import errors from './errors'
import fetching from './fetching'
import accountBalance from './accountBalance'
import bridges from './bridges'
import bridgeSpec from './bridgeSpec'
import jobs from './jobs'
import jobRuns from './jobRuns'
import configuration from './configuration'
import authentication from './authentication'

const reducer = combineReducers({
  errors,
  fetching,
  accountBalance,
  bridges,
  bridgeSpec,
  jobs,
  jobRuns,
  configuration,
  authentication
})

export default reducer
