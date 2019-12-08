import * as types from './types'

export const initialState = {
  options: null,
  contractAddress: null,
  oracles: null,
  currentAnswer: null,
  latestCompletedAnswerId: null,
  pendingAnswerId: null,
  nextAnswerId: null,
  oracleResponse: null,
  requestTime: null,
  minimumResponses: null,
  updateHeight: null,
  answerHistory: null,
  ethGasPrice: null,
}

function clearState(state) {
  const stateCopy = { ...state }
  Object.keys(stateCopy).forEach(k => (stateCopy[k] = null))
  return stateCopy
}

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case types.CLEAR_STATE:
      return {
        ...state,
        ...clearState(state),
      }

    case types.OPTIONS:
      return {
        ...state,
        options: action.payload,
      }

    case types.CONTRACT_ADDRESS:
      return {
        ...state,
        contractAddress: action.payload,
      }

    case types.ORACLES:
      return {
        ...state,
        oracles: action.payload,
      }

    case types.CURRENT_ANSWER:
      return {
        ...state,
        currentAnswer: action.payload,
      }

    case types.LATEST_COMPLETED_ANSWER_ID:
      return {
        ...state,
        latestCompletedAnswerId: action.payload,
      }

    case types.PENDING_ANSWER_ID:
      return {
        ...state,
        pendingAnswerId: action.payload,
      }

    case types.NEXT_ANSWER_ID:
      return {
        ...state,
        nextAnswerId: action.payload,
      }

    case types.ORACLE_RESPONSE:
      return {
        ...state,
        oracleResponse: action.payload,
      }

    case types.REQUEST_TIME:
      return {
        ...state,
        requestTime: action.payload,
      }

    case types.MINUMUM_RESPONSES:
      return {
        ...state,
        minimumResponses: action.payload,
      }

    case types.UPDATE_HEIGHT:
      return {
        ...state,
        updateHeight: action.payload,
      }

    case types.ANSWER_HISTORY:
      return {
        ...state,
        answerHistory: action.payload,
      }

    case types.ETHGAS_PRICE:
      return {
        ...state,
        ethGasPrice: action.payload,
      }

    default:
      return state
  }
}

export default reducer
