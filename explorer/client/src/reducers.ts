import { combineReducers } from 'redux'
import adminAuth from './reducers/adminAuth'
import adminOperators from './reducers/adminOperators'
import adminOperatorsIndex from './reducers/adminOperatorsIndex'
import adminOperatorsShow from './reducers/adminOperatorsShow'
import adminHeads from './reducers/adminHeads'
import adminHeadsIndex from './reducers/adminHeadsIndex'
import chainlinkNodes from './reducers/chainlinkNodes'
import config from './reducers/config'
import jobRuns from './reducers/jobRuns'
import jobRunsIndex from './reducers/jobRunsIndex'
import notifications from './reducers/notifications'
import taskRuns from './reducers/taskRuns'

const reducer = combineReducers({
  adminAuth,
  adminOperators,
  adminOperatorsIndex,
  adminOperatorsShow,
  adminHeads,
  adminHeadsIndex,
  chainlinkNodes,
  config,
  jobRuns,
  jobRunsIndex,
  notifications,
  taskRuns,
})

export const INITIAL_STATE = reducer(undefined, {
  type: 'initial_state' as any,
})
export type AppState = typeof INITIAL_STATE

export default reducer
