import { combineReducers } from 'redux'
import jobs from './jobs'
import accountBalance from './accountBalance'

const reducer = combineReducers({
  accountBalance,
  jobs
})

export default reducer
