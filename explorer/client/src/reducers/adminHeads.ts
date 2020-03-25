import { Actions } from './actions'
import { Reducer } from 'redux'
import { Head } from 'explorer/models'

export interface State {
  items?: Record<string, Head>
}

const INITIAL_STATE: State = {}

export const adminHeads: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_HEADS_SUCCEEDED':
      return { items: { ...action.data.heads } }
    case 'FETCH_ADMIN_HEAD_SUCCEEDED':
      return { items: { ...action.data.heads } }
    default:
      return state
  }
}

export default adminHeads
