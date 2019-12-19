import { Actions } from './actions'
import { Reducer } from 'redux'
import { ChainlinkNode } from 'explorer/models'

export interface State {
  items?: Record<string, ChainlinkNode>
}

const INITIAL_STATE: State = {}

export const adminOperators: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_OPERATORS_SUCCEEDED':
      return { items: { ...action.data.chainlinkNodes } }
    default:
      return state
  }
}

export default adminOperators
