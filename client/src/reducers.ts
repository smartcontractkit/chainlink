import { combineReducers } from 'redux'
import search, { IState as ISearchState } from './reducers/search'
import jobRuns, { IState as IJobRunsState } from './reducers/jobRuns'

export interface IState {
  search: ISearchState,
  jobRuns: IJobRunsState
}

const reducer = combineReducers({
  search,
  jobRuns
})

export default reducer
