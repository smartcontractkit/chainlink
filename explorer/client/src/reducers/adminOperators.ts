import { Actions } from './actions'

export interface State {
  items?: ChainlinkNode[]
}

const INITIAL_STATE: State = {}

export default (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'FETCH_ADMIN_OPERATORS_SUCCEEDED':
      return { items: { ...action.data.chainlinkNodes } }
    default:
      return state
  }
}
