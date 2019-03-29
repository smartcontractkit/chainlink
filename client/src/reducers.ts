import { combineReducers } from 'redux'
import search, { IState as ISearchState } from './reducers/search'
import jobRuns, { IState as IJobRunsState } from './reducers/jobRuns'
import jobRunsIndex, {
  IState as IJobRunsIndexState
} from './reducers/jobRunsIndex'

export interface IState {
  search: ISearchState
  jobRuns: IJobRunsState
  jobRunsIndex: IJobRunsIndexState
}

const reducer = combineReducers({
  search,
  jobRuns,
  jobRunsIndex
})

export default reducer
