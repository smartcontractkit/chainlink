import { combineReducers } from 'redux'
import search, { IState as ISearchState } from './reducers/search'
import jobRuns, { IState as IJobRunsState } from './reducers/jobRuns'
import taskRuns, { IState as ITaskRunsState } from './reducers/taskRuns'
import jobRunsIndex, {
  IState as IJobRunsIndexState
} from './reducers/jobRunsIndex'

export interface IState {
  search: ISearchState
  taskRuns: ITaskRunsState
  jobRuns: IJobRunsState
  jobRunsIndex: IJobRunsIndexState
}

const reducer = combineReducers({
  search,
  taskRuns,
  jobRuns,
  jobRunsIndex
})

export default reducer
