import { Actions } from 'state/actions'

export interface State {
  tooltip: null | any
  drawer: null | any
}

export const INITIAL_STATE: State = {
  tooltip: null,
  drawer: null,
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'networkGraph/SET_TOOLTIP':
      return {
        ...state,
        tooltip: action.payload,
      }

    case 'networkGraph/SET_DRAWER':
      return {
        ...state,
        drawer: action.payload,
      }

    default:
      return state
  }
}

export default reducer
