import { combineReducers } from 'redux'
import accountBalance from './accountBalance'
import jobs from './jobs'
import jobRuns from './jobRuns'

const reducer = combineReducers({
  accountBalance,
  jobs,
  jobRuns
})

export default reducer
