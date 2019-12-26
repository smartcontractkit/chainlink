import { combineReducers } from 'redux'
import adminAuth from './reducers/adminAuth'
import adminOperators from './reducers/adminOperators'
import adminOperatorsIndex from './reducers/adminOperatorsIndex'
import adminOperatorsShow from './reducers/adminOperatorsShow'
import chainlinkNodes from './reducers/chainlinkNodes'
import config from './reducers/config'
import jobRuns from './reducers/jobRuns'
import jobRunsIndex from './reducers/jobRunsIndex'
import notifications from './reducers/notifications'
import search from './reducers/query'
import taskRuns from './reducers/taskRuns'

const reducer = combineReducers({
  adminAuth,
  adminOperators,
  adminOperatorsIndex,
  adminOperatorsShow,
  chainlinkNodes,
  config,
  jobRuns,
  jobRunsIndex,
  notifications,
  search,
  taskRuns,
})

export const INITIAL_STATE = reducer(undefined, { type: 'initial_state' })
export type AppState = typeof INITIAL_STATE

export default reducer
