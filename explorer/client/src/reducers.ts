import { combineReducers } from 'redux'
import adminAuth, { State as AdminAuthState } from './reducers/adminAuth'
import chainlinkNodes, {
  State as ChainlinkNodesState,
} from './reducers/chainlinkNodes'
import config, { State as ConfigState } from './reducers/config'
import jobRuns, { State as JobRunsState } from './reducers/jobRuns'
import jobRunsIndex, {
  State as JobRunsIndexState,
} from './reducers/jobRunsIndex'
import notifications, {
  State as NotificationsState,
} from './reducers/notifications'
import search, { State as SearchState } from './reducers/search'
import taskRuns, { State as TaskRunsState } from './reducers/taskRuns'

export interface State {
  adminAuth: AdminAuthState
  chainlinkNodes: ChainlinkNodesState
  config: ConfigState
  jobRuns: JobRunsState
  jobRunsIndex: JobRunsIndexState
  notifications: NotificationsState
  search: SearchState
  taskRuns: TaskRunsState
}

const reducer = combineReducers({
  adminAuth,
  chainlinkNodes,
  config,
  jobRuns,
  jobRunsIndex,
  notifications,
  search,
  taskRuns,
})

export default reducer
