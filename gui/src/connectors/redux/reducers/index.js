import { combineReducers } from 'redux'
import accountBalance from './accountBalance'
import jobs from './jobs'
import jobRuns from './jobRuns'
import configuration from './configuration'

const reducer = combineReducers({
  accountBalance,
  jobs,
  jobRuns,
  configuration
})

export default reducer
