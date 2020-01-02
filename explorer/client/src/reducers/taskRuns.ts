import { Actions } from './actions'
import { Reducer } from 'redux'
import { TaskRun } from 'explorer/models'

export interface State {
  items?: TaskRun[]
}

const INITIAL_STATE: State = { items: undefined }

const taskRunsReducer: Reducer<State, Actions> = (
  state: State = INITIAL_STATE,
  action: Actions,
) => {
  switch (action.type) {
    case 'FETCH_JOB_RUN_SUCCEEDED':
      return { items: action.data.taskRuns }
    default:
      return state
  }
}

export default taskRunsReducer
