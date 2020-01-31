import * as types from './types'

export const initialState = {
  answers: null,
}

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case types.SET_ANSWERS:
      return {
        ...state,
        answers: action.payload,
      }

    default:
      return state
  }
}

export default reducer
