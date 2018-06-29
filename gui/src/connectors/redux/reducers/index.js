import { combineReducers } from 'redux'
import accountBalance from './accountBalance'
import bridges from './bridges'
import jobs from './jobs'
import jobRuns from './jobRuns'
import configuration from './configuration'

const reducer = combineReducers({
  accountBalance,
  bridges,
  jobs,
  jobRuns,
  configuration
})

export default reducer
