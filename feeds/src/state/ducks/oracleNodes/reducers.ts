import 'core-js/stable/object/from-entries'
import { OracleNode, getOracleNodes } from '../../../config'
import { Actions } from 'state/actions'

export interface State {
  items: Record<OracleNode['address'], OracleNode>
  order: Array<OracleNode['address']>
}

export const INITIAL_STATE: State = {
  items: Object.fromEntries(getOracleNodes().map(o => [o.address, o])),
  order: getOracleNodes().map(o => o.address),
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    default:
      return state
  }
}

export default reducer
