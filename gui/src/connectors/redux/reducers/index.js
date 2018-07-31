import { combineReducers } from 'redux'
import errors from './errors'
import accountBalance from './accountBalance'
import bridges from './bridges'
import bridgeSpec from './bridgeSpec'
import jobs from './jobs'
import jobRuns from './jobRuns'
import configuration from './configuration'
import session from './session'

const reducer = combineReducers({
  errors,
  accountBalance,
  bridges,
  bridgeSpec,
  jobs,
  jobRuns,
  configuration,
  session
})

export default reducer
