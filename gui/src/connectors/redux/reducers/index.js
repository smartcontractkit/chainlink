import { combineReducers } from 'redux'
import accountBalance from './accountBalance'
import bridges from './bridges'
import bridgeSpec from './bridgeSpec'
import jobs from './jobs'
import jobRuns from './jobRuns'
import configuration from './configuration'

const reducer = combineReducers({
  accountBalance,
  bridges,
  bridgeSpec,
  jobs,
  jobRuns,
  configuration
})

export default reducer
