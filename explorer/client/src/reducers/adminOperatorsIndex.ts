import { Actions } from './actions'

export interface State {
  items?: string[]
  count?: number
}

const INITIAL_STATE: State = {}

export default (state: State = INITIAL_STATE, action: Actions): State => {
  switch (action.type) {
    case 'FETCH_ADMIN_OPERATORS_SUCCEEDED':
      return {
        items: action.data.meta.currentPageOperators.data.map(o => o.id),
        count: action.data.meta.currentPageOperators.meta.count,
      }
    default:
      return state
  }
}
