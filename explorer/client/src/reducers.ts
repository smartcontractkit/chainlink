import { combineReducers } from 'redux'
import config, { State as ConfigState } from './reducers/config'
import search, { State as SearchState } from './reducers/search'
import jobRuns, { State as JobRunsState } from './reducers/jobRuns'
import taskRuns, { State as TaskRunsState } from './reducers/taskRuns'
import chainlinkNodes, {
  State as ChainlinkNodesState,
} from './reducers/chainlinkNodes'
import jobRunsIndex, {
  State as JobRunsIndexState,
} from './reducers/jobRunsIndex'

export interface State {
  config: ConfigState
  chainlinkNodes: ChainlinkNodesState
  jobRuns: JobRunsState
  taskRuns: TaskRunsState
  search: SearchState
  jobRunsIndex: JobRunsIndexState
}

const reducer = combineReducers({
  config,
  chainlinkNodes,
  jobRuns,
  taskRuns,
  search,
  jobRunsIndex,
})

export default reducer
