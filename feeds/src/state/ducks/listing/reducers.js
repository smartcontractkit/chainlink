import * as types from './types'

export const initialState = {
  answers: null,
  healthPrices: {},
}

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case types.SET_ANSWERS:
      return {
        ...state,
        answers: action.payload,
      }

    case types.SET_HEALTH_PRICE: {
      const [item, currentPrice] = action.payload

      return {
        ...state,
        healthPrices: {
          ...state.healthPrices,
          ...{ [item.contractAddress]: currentPrice },
        },
      }
    }

    default:
      return state
  }
}

export default reducer
