import { Actions } from './actions'
import { Reducer } from 'redux'
import { ChainlinkNode } from 'explorer/models'

export interface State {
  items?: Record<number, ChainlinkNode>
}

const INITIAL_STATE: State = { items: undefined }

export const chainlinkNodesReducer: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_JOB_RUNS_SUCCEEDED':
      return { items: { ...action.data.chainlinkNodes } }
    case 'FETCH_JOB_RUN_SUCCEEDED':
      return { items: { ...action.data.chainlinkNodes } }
    default:
      return state
  }
}

export default chainlinkNodesReducer
