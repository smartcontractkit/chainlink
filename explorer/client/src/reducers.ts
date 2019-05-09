import { combineReducers } from 'redux'
import search, { IState as ISearchState } from './reducers/search'
import jobRuns, { IState as IJobRunsState } from './reducers/jobRuns'
import taskRuns, { IState as ITaskRunsState } from './reducers/taskRuns'
import chainlinkNodes, {
  IState as IChainlinkNodesState
} from './reducers/chainlinkNodes'
import jobRunsIndex, {
  IState as IJobRunsIndexState
} from './reducers/jobRunsIndex'

export interface IState {
  chainlinkNodes: IChainlinkNodesState
  jobRuns: IJobRunsState
  taskRuns: ITaskRunsState
  search: ISearchState
  jobRunsIndex: IJobRunsIndexState
}

const reducer = combineReducers({
  chainlinkNodes,
  jobRuns,
  taskRuns,
  search,
  jobRunsIndex
})

export default reducer
