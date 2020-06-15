import { Reducer } from 'redux'
import { SET_TOOLTIP, SET_DRAWER, NetworkGraphActionTypes } from './types'

export interface State {
  tooltip: null | any
  drawer: null | any
}

export const INITIAL_STATE: State = {
  tooltip: null,
  drawer: null,
}

const reducer: Reducer<State, NetworkGraphActionTypes> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case SET_TOOLTIP:
      return {
        ...state,
        tooltip: action.payload,
      }

    case SET_DRAWER:
      return {
        ...state,
        drawer: action.payload,
      }

    default:
      return state
  }
}

export default reducer
