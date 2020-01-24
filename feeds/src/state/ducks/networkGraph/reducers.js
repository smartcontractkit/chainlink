import * as types from './types'

export const initialState = {
  tooltip: null,
  drawer: null,
}

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case types.SET_TOOLTIP:
      return {
        ...state,
        tooltip: action.payload,
      }

    case types.SET_DRAWER:
      return {
        ...state,
        drawer: action.payload,
      }

    default:
      return state
  }
}

export default reducer
